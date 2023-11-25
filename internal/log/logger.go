package log

import (
	"log/slog"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Log struct {
	FilePath string `yaml:"filePath"`
}

var logger *zap.SugaredLogger

func init() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"cardago.log"}
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapper, err := config.Build()
	if err != nil {
		slog.Error("could not initialize config", "ERROR", err)
	}

	logger = zapper.Sugar()
	logger.Info("Logger initialized")
}

func Infow(msg string, keysAndValues ...interface{}) {
	logger.Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	logger.Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	logger.Errorw(msg, keysAndValues...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	logger.Debugw(msg, keysAndValues...)
}
