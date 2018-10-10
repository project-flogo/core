package engine

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/project-flogo/core/app"
	"github.com/project-flogo/core/support/logger"
	"github.com/project-flogo/core/support/managed"
)

var managedServices []managed.Managed
var lock = &sync.Mutex{}

// Interface for the engine behaviour
type Engine interface {
	//App get the application associated with this engine
	App() *app.App

	// Start starts the engine
	Start() error

	// Stop stop the engine
	Stop() error
}

func LifeCycle(managedEntity managed.Managed) {
	defer lock.Unlock()
	lock.Lock()
	managedServices = append(managedServices, managedEntity)
}

// NewFromConfigProvider creates a new Engine using the specified AppConfigProvider
func NewFromConfigProvider(provider AppConfigProvider) (Engine, error) {

	if provider == nil {
		provider = DefaultAppConfigProvider()
	}

	appConfig, err := provider.GetAppConfig()
	if err != nil {
		return nil, err
	}

	return New(appConfig)
}

func RunEngine(e Engine) int {

	err := e.Start()
	if err != nil {
		fmt.Printf("Failed to start engine due to error: %s\n", err.Error())
		os.Exit(1)
	}

	exitChan := setupSignalHandling()

	code := <-exitChan

	e.Stop()

	return code
}

func setupSignalHandling() chan int {

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	exitChan := make(chan int, 1)
	select {
	case s := <-signalChan:
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			exitChan <- 0
		default:
			logger.Debug("Unknown signal.")
			exitChan <- 1
		}
	}
	return exitChan
}
