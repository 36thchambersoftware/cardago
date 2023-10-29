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
	config.LoadConfig()
	slog.Info("CARDAGO", "PACKAGE", "CONFIG", "RUNTIME", config)

	syncProgress := cardano.QueryTip().SyncProgress
	slog.Info("CARDAGO", "PACKAGE", "CARDANO", "SYNCPROGRESS", syncProgress)

	message := fmt.Sprintf("<@%s> SYNC PROGRESS: %s", config.Discord.UserID, syncProgress)

	discord.NotifyChannel(config.Discord, message)
}
