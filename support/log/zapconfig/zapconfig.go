package zapconfig

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Format int

const (
	EnvKeyLogFormat         = "FLOGO_LOG_FORMAT"
	DefaultLogFormat        = FormatConsole
	FormatConsole    Format = iota
	FormatJSON
)

type defaultCfgImpl struct {
	logConfig      zap.Config
	logLevel       *zap.AtomicLevel
	traceLogConfig zap.Config
	traceLogLevel  *zap.AtomicLevel
}

// DefaultConfig returns default configuration values
type DefaultConfig interface {
	LogCfg() zap.Config
	LogLvl() *zap.AtomicLevel
	TraceLogCfg() zap.Config
	TraceLogLvl() *zap.AtomicLevel
}

var defaultCfg DefaultConfig

func init() {
	defaultCfg = createDefaultConfiguration()
}

// DefaultCfg returns default configuration
func DefaultCfg() DefaultConfig {
	return defaultCfg
}

func (d *defaultCfgImpl) LogCfg() zap.Config {
	return d.logConfig
}

func (d *defaultCfgImpl) LogLvl() *zap.AtomicLevel {
	return d.logLevel
}

func (d *defaultCfgImpl) TraceLogCfg() zap.Config {
	return d.traceLogConfig
}

func (d *defaultCfgImpl) TraceLogLvl() *zap.AtomicLevel {
	return d.traceLogLevel
}

func createDefaultConfiguration() DefaultConfig {

	logFormat := DefaultLogFormat
	envLogFormat := strings.ToUpper(os.Getenv(EnvKeyLogFormat))
	if envLogFormat == "JSON" {
		logFormat = FormatJSON
	}

	cfg := zap.NewProductionConfig()
	cfg.DisableCaller = true

	eCfg := cfg.EncoderConfig
	eCfg.TimeKey = "timestamp"
	eCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	//eCfg.EncodeTime = zapcore.EpochNanosTimeEncoder

	if logFormat == FormatConsole {
		eCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		cfg.Encoding = "console"
		eCfg.EncodeName = nameEncoder
	}

	cfg.EncoderConfig = eCfg

	lvl := cfg.Level

	// trace log configuration
	tcfg := cfg

	if strings.Compare(tcfg.Encoding, "console") == 0 {
		tcfg.EncoderConfig.EncodeLevel = traceLevelEncoder
	}

	tlvl := tcfg.Level
	tlvl.SetLevel(zapcore.DebugLevel)

	defaultCfg := &defaultCfgImpl{
		logConfig:      cfg,
		logLevel:       &lvl,
		traceLogConfig: tcfg,
		traceLogLevel:  &tlvl,
	}

	return defaultCfg
}

func nameEncoder(loggerName string, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + loggerName + "] -")
}

func traceLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[TRACE]")
}
