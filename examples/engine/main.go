package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/core/support/logger"
)

var log = logger.GetLogger("main-engine")
var cpuProfile = flag.String("cpuprofile", "", "Writes CPU profiling for the current process to the specified file")
var memProfile = flag.String("memprofile", "", "Writes memory profiling for the current process to the specified file")

var (
	configProvider engine.AppConfigProvider
)

func main() {

	flag.Parse()
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			fmt.Println(fmt.Sprintf("Failed to create CPU profiling file due to error - %s", err.Error()))
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	e, err := engine.NewFromConfigProvider(configProvider)
	if err != nil {
		log.Errorf("Failed to create engine instance due to error: %s", err.Error())
		os.Exit(1)
	}

	code := engine.RunEngine(e)

	if *memProfile != "" {
		f, err := os.Create(*memProfile)
		if err != nil {
			fmt.Println(fmt.Sprintf("Failed to create memory profiling file due to error - %s", err.Error()))
			os.Exit(1)
		}

		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			fmt.Println(fmt.Sprintf("Failed to write memory profiling data to file due to error - %s", err.Error()))
			os.Exit(1)
		}
		f.Close()
	}

	os.Exit(code)
}
