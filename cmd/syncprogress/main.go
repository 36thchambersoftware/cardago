package main

import (
	"fmt"

	"cardano/cardago/internal/cardano"
	"cardano/cardago/internal/config"
	"cardano/cardago/internal/discord"
	"cardano/cardago/internal/log"
)

func main() {
	config := config.Get()
	err := config.LoadConfig()
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CONFIG", "ERROR", err)
		return
	}
	log.Debugw("CARDAGO", "PACKAGE", "CONFIG", "RUNTIME", config)

	tip, err := cardano.QueryTip()
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
		return
	}

	log.Infow("CARDAGO", "PACKAGE", "CARDANO", "SYNCPROGRESS", tip.SyncProgress)

	message := fmt.Sprintf("<@%s> SYNC PROGRESS: %s", config.Discord.UserID, tip.SyncProgress)

	discord.ExecuteWebhook(config.Discord, message)
}
