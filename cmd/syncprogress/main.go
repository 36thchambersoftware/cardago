package main

import (
	"fmt"

	"cardano/cardago/internal/cardano"
	"cardano/cardago/internal/config"
	"cardano/cardago/internal/discord"
	"cardano/cardago/internal/log"
)

func main() {
	logger := log.InitializeLogger()
	config := config.Get()
	err := config.LoadConfig()
	if err != nil {
		logger.Errorw("CARDAGO", "PACKAGE", "CONFIG", "ERROR", err)
		return
	}
	logger.Infow("CARDAGO", "PACKAGE", "CONFIG", "RUNTIME", config)

	tip, err := cardano.QueryTip()
	if err != nil {
		logger.Errorw("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
		return
	}

	logger.Infow("CARDAGO", "PACKAGE", "CARDANO", "SYNCPROGRESS", tip.SyncProgress)

	message := fmt.Sprintf("<@%s> SYNC PROGRESS: %s", config.Discord.UserID, tip.SyncProgress)

	discord.ExecuteWebhook(config.Discord, message)
}
