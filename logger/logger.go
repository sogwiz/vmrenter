// ZAP logger will be configured and built here
package zaplogger

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var cfg zap.Config

var logLevels = map[string]zapcore.Level{
	"Debug":  zapcore.DebugLevel,
	"Info":   zapcore.InfoLevel,
	"Warn":   zapcore.WarnLevel,
	"Error":  zapcore.ErrorLevel,
	"DPanic": zapcore.DPanicLevel,
	"Panic":  zapcore.PanicLevel,
	"Fatal":  zapcore.FatalLevel,
}

func BuildLogger(logLevel string) (*zap.Logger, error) {

	for lvl, zapLevel := range logLevels {
		if logLevel == lvl {
			cfg = zap.Config{
				Encoding:    "json",
				Level:       zap.NewAtomicLevelAt(zapLevel),
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

			fmt.Printf("Log level %v", logLevel)
			logger, err := cfg.Build()
			return logger, err
		}
	}

	return nil, errors.New("invalid log level")
}
