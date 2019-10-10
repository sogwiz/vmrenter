// ZAP logger will be configured and built here
package zaplogger

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
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

func ConfigureLogger(logLevel string) error {

	for lvl, zapLevel := range logLevels {
		if logLevel == lvl {
			cfg = zap.Config{
				Encoding:    "json",
				Level:       zap.NewAtomicLevelAt(zapLevel),
				OutputPaths: []string{"stdout", "./logs.log"},
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

			fmt.Printf("Log level %v\n", logLevel)
			logger, err := cfg.Build()
			zap.ReplaceGlobals(logger)
			if err != nil {
				return err
			}
			zap.S().Info("Successfully configured logger")
			return nil
		}
	}

	var possibleLogLevels []string

	for k := range logLevels {
		possibleLogLevels = append(possibleLogLevels, k)
	}

	possibleLogLevelsStr := strings.Join(possibleLogLevels, ", ")

	errorMsg := fmt.Sprintf("invalid log level - %v, possible log levels - %v", logLevel, possibleLogLevelsStr)

	return errors.New(errorMsg)
}
