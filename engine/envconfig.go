package engine

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/project-flogo/core/app/propertyresolver"
	"github.com/project-flogo/core/data/property"
	"github.com/project-flogo/core/engine/runner"
	"github.com/project-flogo/core/support/log"
)

const (
	EnvKeyAppConfigLocation     = "FLOGO_CONFIG_PATH"
	DefaultAppConfigLocation    = "flogo.json"
	EnvKeyEngineConfigLocation  = "FLOGO_ENG_CONFIG_PATH"
	DefaultEngineConfigLocation = "engine.json"

	EnvKeyStopEngineOnError  = "FLOGO_ENGINE_STOP_ON_ERROR"
	DefaultStopEngineOnError = true
	EnvKeyRunnerType         = "FLOGO_RUNNER_TYPE"
	DefaultRunnerType        = ValueRunnerTypePooled
	EnvKeyRunnerWorkers      = "FLOGO_RUNNER_WORKERS"
	DefaultRunnerWorkers     = 5

	//Deprecated
	EnvKeyRunnerQueueSizeLegacy = "FLOGO_RUNNER_QUEUE"

	EnvKeyRunnerQueueSize  = "FLOGO_RUNNER_QUEUE_SIZE"
	DefaultRunnerQueueSize = 50

	EnvAppPropertyResolvers   = "FLOGO_APP_PROP_RESOLVERS"
	EnvEnableSchemaSupport    = "FLOGO_SCHEMA_SUPPORT"
	EnvEnableSchemaValidation = "FLOGO_SCHEMA_VALIDATION"

	ValueRunnerTypePooled = "POOLED"
	ValueRunnerTypeDirect = "DIRECT"
)

func IsSchemaSupportEnabled() bool {
	schemaValidationEnv := os.Getenv(EnvEnableSchemaSupport)
	if strings.EqualFold(schemaValidationEnv, "true") {
		return true
	}

	return false
}

func IsSchemaValidationEnabled() bool {
	schemaValidationEnv := os.Getenv(EnvEnableSchemaValidation)
	if !strings.EqualFold(schemaValidationEnv, "true") {
		return false
	}

	return true
}

//GetFlogoAppConfigPath returns the flogo config path
func GetFlogoAppConfigPath() string {

	flogoConfigPathEnv := os.Getenv(EnvKeyAppConfigLocation)
	if len(flogoConfigPathEnv) > 0 {
		return flogoConfigPathEnv
	}

	if _, err := os.Stat(DefaultAppConfigLocation); err != nil {
		upDirConfig := filepath.Join("..", DefaultAppConfigLocation)
		if _, err := os.Stat(upDirConfig); err == nil {
			return upDirConfig
		}
	}

	return DefaultAppConfigLocation
}

//GetFlogoEngineConfigPath returns the flogo engine config path
func GetFlogoEngineConfigPath() string {

	flogoConfigPathEnv := os.Getenv(EnvKeyEngineConfigLocation)
	if len(flogoConfigPathEnv) > 0 {
		return flogoConfigPathEnv
	}

	if _, err := os.Stat(DefaultEngineConfigLocation); err != nil {
		upDirConfig := filepath.Join("..", DefaultEngineConfigLocation)
		if _, err := os.Stat(upDirConfig); err == nil {
			return upDirConfig
		}
	}

	return DefaultEngineConfigLocation
}

//func SetDefaultLogLevel(logLevel string) {
//	defaultLogLevel = logLevel
//}
//
////GetLogLevel returns the log level
//func GetLogLevel() string {
//	logLevelEnv := os.Getenv(EnvKeyLogLevel)
//	if len(logLevelEnv) > 0 {
//		return logLevelEnv
//	}
//	return defaultLogLevel
//}
//
//func GetLogDateTimeFormat() string {
//	logLevelEnv := os.Getenv(EnvKeyLogDateFormat)
//	if len(logLevelEnv) > 0 {
//		return logLevelEnv
//	}
//	return DefaultLogDateFormat
//}

func StopEngineOnError() bool {
	stopEngineOnError := os.Getenv(EnvKeyStopEngineOnError)
	if len(stopEngineOnError) == 0 {
		return DefaultStopEngineOnError
	}
	b, _ := strconv.ParseBool(stopEngineOnError)
	return b
}

func displayAppPropertyValueResolversHelp(logger log.Logger, resolvers []string) {
	logger.Warn("Multiple property resolvers where defined without setting a priority order!")
	logger.Infof("Set environment variable '%s' with a comma-separated list of resolvers to use (definition order is decreasing order of priority)", EnvAppPropertyResolvers)
	logger.Infof("List of available resolvers: %v", resolvers)
	logger.Warn("No property resolver will be used")
}

func GetAppPropertyValueResolvers(logger log.Logger) string {
	key := os.Getenv(EnvAppPropertyResolvers)
	if len(key) > 0 {
		if key == "disabled" {
			return ""
		}
		return key
	}

	// EnvAppPropertyResolvers is not set, let's guess some convenient default behaviours
	switch len(property.RegisteredResolvers) {
	case 0: // no resolver, do nothing
		return ""
	case 1: // only one resolver has been registered, use it
		for resolver := range property.RegisteredResolvers {
			return resolver
		}
	case 2, 3:
		var resolvers, builtinResolvers []string

		for resolver := range property.RegisteredResolvers {
			if resolver != propertyresolver.ResolverNameEnv && resolver != propertyresolver.ResolverNameJson {
				resolvers = append(resolvers, resolver)
			} else {
				builtinResolvers = append(builtinResolvers, resolver)
			}
		}

		if len(resolvers) > 1 { // multiple (excluding builtin) resolvers defined, do nothing and hint to enforce an order
			resolvers = append(resolvers, builtinResolvers...)
			displayAppPropertyValueResolversHelp(logger, resolvers)
			return ""
		}

		if len(builtinResolvers) == 2 { // force priority between the two builtin resolvers
			builtinResolvers = []string{propertyresolver.ResolverNameEnv, propertyresolver.ResolverNameJson}
		}

		resolvers = append(resolvers, builtinResolvers...)

		return strings.Join(resolvers[:], ",")
	default: // multiple (excluding builtin) resolvers defined, do nothing and hint to enforce an order
		var resolvers []string

		for resolver := range property.RegisteredResolvers {
			resolvers = append(resolvers, resolver)
		}

		displayAppPropertyValueResolversHelp(logger, resolvers)

		return ""
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
	} else {
		//For backward compatible.
		legacyQueueSize := os.Getenv(EnvKeyRunnerQueueSizeLegacy)
		if len(legacyQueueSize) > 0 {
			i, err := strconv.Atoi(queueSizeEnv)
			if err == nil {
				queueSize = i
			}
		}
	}
	return queueSize
}

//NewPooledRunnerConfig creates a new Pooled config, looks for environment variables to override default values
func NewPooledRunnerConfig() *runner.PooledConfig {
	return &runner.PooledConfig{NumWorkers: GetRunnerWorkers(), WorkQueueSize: GetRunnerQueueSize()}
}

func ConfigViaEnv(e *engineImpl) {

	config := &Config{}
	//config.LogLevel = GetLogLevel()
	config.RunnerType = GetRunnerType()
	config.StopEngineOnError = StopEngineOnError()

	e.config = config
}

func DirectRunner(e *engineImpl) error {
	e.logger.Debugf("Using 'DIRECT' Action Runner")
	e.actionRunner = runner.NewDirect()
	return nil
}
