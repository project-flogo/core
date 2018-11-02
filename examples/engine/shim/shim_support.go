package embedded

import (
	"os"

	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/core/support/log"
)

var (
	cfgJson       string
	cfgCompressed bool
)

func init() {
	log.SetLogLevel(log.RootLogger(), log.ErrorLevel)

	cfg, err := engine.LoadAppConfig(cfgJson, cfgCompressed)
	if err != nil {
		log.RootLogger().Errorf("Failed to create engine: %s", err.Error())
		os.Exit(1)
	}

	_, err = engine.New(cfg)
	if err != nil {
		log.RootLogger().Errorf("Failed to create engine: %s", err.Error())
		os.Exit(1)
	}
}
