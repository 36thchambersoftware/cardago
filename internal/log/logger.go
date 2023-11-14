package log

import (
	"log/slog"

	"go.uber.org/zap"
)

func InitializeLogger() *zap.SugaredLogger {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"/home/cardano/scripts/cardago/cardago.log"}
	logger, err := config.Build()
	if err != nil {
		slog.Error("could not initialize config", "ERROR", err)
	}

	sugar := logger.Sugar()
	sugar.Info("Logger initialized")

	return sugar
}
