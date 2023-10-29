package main

import (
	"log/slog"

	"cardano/cardago/internal/cardano"
	"cardano/cardago/internal/config"
	"cardano/cardago/internal/discord"
)

func main() {
	config := config.Get()
	config.LoadConfig()
	slog.Info("CARDAGO", "PACKAGE", "CONFIG", "RUNTIME", config)

	scheduledBlocks := cardano.GetScheduledBlocks(config.LeaderLogDirectory, config.LeaderLogPrefix, config.LeaderLogExtension)

	discord.NotifyScheduledBlocks(config.Discord, scheduledBlocks)
}
