package engine

import (
	"os"
	"strconv"

	"github.com/project-flogo/core/engine/runner"
	"github.com/project-flogo/core/support/logger"
)

const (
	EnvKeyLogDateFormat      = "FLOGO_LOG_DTFORMAT"
	DefaultLogDateFormat     = "2006-01-02 15:04:05.000"
	EnvKeyLogLevel           = "FLOGO_LOG_LEVEL"
	DefaultLogLevel          = "INFO"
	EnvKeyAppConfigLocation  = "FLOGO_CONFIG_PATH"
	DefaultAppConfigLocation = "flogo.json"
	EnvKeyStopEngineOnError  = "FLOGO_ENGINE_STOP_ON_ERROR"
	DefaultStopEngineOnError = true
	EnvKeyRunnerType         = "FLOGO_RUNNER_TYPE"
	DefaultRunnerType        = ValueRunnerTypePooled
	EnvKeyRunnerWorkers      = "FLOGO_RUNNER_WORKERS"
	DefaultRunnerWorkers     = 5
	EnvKeyRunnerQueueSize    = "FLOGO_RUNNER_QUEUE"
	DefaultRunnerQueueSize   = 50
	EnvAppPropertyOverride   = "FLOGO_APP_PROP_OVERRIDE"
	EnvAppPropertyProvider   = "FLOGO_APP_PROP_PROVIDER"

	ValueRunnerTypePooled = "POOLED"
	ValueRunnerTypeDirect = "DIRECT"
)

var defaultLogLevel = DefaultLogLevel

//GetFlogoConfigPath returns the flogo config path
func GetFlogoConfigPath() string {
	flogoConfigPathEnv := os.Getenv(EnvKeyAppConfigLocation)
	if len(flogoConfigPathEnv) > 0 {
		return flogoConfigPathEnv
	}
	return DefaultAppConfigLocation
}

func SetDefaultLogLevel(logLevel string) {
	defaultLogLevel = logLevel
}

//GetLogLevel returns the log level
func GetLogLevel() string {
	logLevelEnv := os.Getenv(EnvKeyLogLevel)
	if len(logLevelEnv) > 0 {
		return logLevelEnv
	}
	return defaultLogLevel
}

func GetLogDateTimeFormat() string {
	logLevelEnv := os.Getenv(EnvKeyLogDateFormat)
	if len(logLevelEnv) > 0 {
		return logLevelEnv
	}
	return DefaultLogDateFormat
}

func StopEngineOnError() bool {
	stopEngineOnError := os.Getenv(EnvKeyStopEngineOnError)
	if len(stopEngineOnError) == 0 {
		return DefaultStopEngineOnError
	}
	b, _ := strconv.ParseBool(stopEngineOnError)
	return b
}

func GetAppPropertyOverride() string {
	key := os.Getenv(EnvAppPropertyOverride)
	if len(key) > 0 {
		return key
	}
	return ""
}

func GetAppPropertyProvider() string {
	key := os.Getenv(EnvAppPropertyProvider)
	if len(key) > 0 {
		return key
	}
	return ""
}

//GetRunnerType returns the runner type
func GetRunnerType() string {
	runnerTypeEnv := os.Getenv(EnvKeyRunnerType)
	if len(runnerTypeEnv) > 0 {
		return runnerTypeEnv
	}
	return DefaultRunnerType
}

//GetRunnerWorkers returns the number of workers to use
func GetRunnerWorkers() int {
	numWorkers := DefaultRunnerWorkers
	workersEnv := os.Getenv(EnvKeyRunnerWorkers)
	if len(workersEnv) > 0 {
		i, err := strconv.Atoi(workersEnv)
		if err == nil {
			numWorkers = i
		}
	}
	return numWorkers
}

//GetRunnerQueueSize returns the runner queue size
func GetRunnerQueueSize() int {
	queueSize := DefaultRunnerQueueSize
	queueSizeEnv := os.Getenv(EnvKeyRunnerQueueSize)
	if len(queueSizeEnv) > 0 {
		i, err := strconv.Atoi(queueSizeEnv)
		if err == nil {
			queueSize = i
		}
	}
	return queueSize
}

//NewPooledRunnerConfig creates a new Pooled config, looks for environment variables to override default values
func NewPooledRunnerConfig() *runner.PooledConfig {
	return &runner.PooledConfig{NumWorkers: GetRunnerWorkers(), WorkQueueSize: GetRunnerQueueSize()}
}

type Config struct {
	LogLevel          string
	StopEngineOnError bool
	RunnerType        string
}

func ConfigViaEnv(e *engineImpl) {

	config := &Config{}
	config.LogLevel = GetLogLevel()
	config.RunnerType = GetRunnerType()
	config.StopEngineOnError = StopEngineOnError()

	e.config = config
}

func DirectRunner(e *engineImpl) {
	logger.Debugf("Using 'DIRECT' Action Runner")
	e.actionRunner = runner.NewDirect()
}
