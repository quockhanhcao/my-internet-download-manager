package utils

import (
	"github.com/quockhanhcao/my-internet-download-manager/internal/configs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getLoggerLevel(level string) zap.AtomicLevel {
	switch level {
	case "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
    case "panic":
        return zap.NewAtomicLevelAt(zap.PanicLevel)
	default:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}
}

func InitializeLogger(configs configs.LogConfig) (*zap.Logger, func(), error) {
    config := zap.NewProductionConfig()
    config.Level = getLoggerLevel(configs.Level)
    config.OutputPaths = configs.OutputPaths
    config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    logger, err := config.Build()
    if err != nil {
        return nil, nil, err
    }
    cleanup := func() {
        logger.Sync()
    }
    return logger, cleanup, nil
}
