// ZAP logger will be configured and built here
package zaplogger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var cfg zap.Config

func buildLogger() *zap.Logger {
	cfg = zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "message",
			LevelKey:      "level",
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			TimeKey:       "time",
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			CallerKey:     "caller",
			EncodeCaller:  zapcore.ShortCallerEncoder,
			NameKey:       "name",
			StacktraceKey: "stackTrace",
		},
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger

}
