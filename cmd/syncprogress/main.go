package main

import (
	"fmt"
	"log/slog"

	"cardano/cardago/internal/cardano"
	"cardano/cardago/internal/config"
	"cardano/cardago/internal/discord"
)

func main() {
	config := config.Get()
	err := config.LoadConfig()
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CONFIG", "ERROR", err)
		return
	}
	slog.Info("CARDAGO", "PACKAGE", "CONFIG", "RUNTIME", config)

	tip, err := cardano.QueryTip()
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
		return
	}

	slog.Info("CARDAGO", "PACKAGE", "CARDANO", "SYNCPROGRESS", tip.SyncProgress)

	message := fmt.Sprintf("<@%s> SYNC PROGRESS: %s", config.Discord.UserID, tip.SyncProgress)

	discord.NotifyChannel(config.Discord, message)
}
