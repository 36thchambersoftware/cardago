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

	scheduledBlocks := cardano.GetScheduledBlocks(config.LeaderLogDirectory, config.LeaderLogPrefix, config.LeaderLogExtension)

	scheduledBlocksMessage := fmt.Sprintf("<@%s> SCHEDULED BLOCKS: %s", config.Discord.UserID, scheduledBlocks)
	if len(scheduledBlocks) == 0 {
		scheduledBlocksMessage = fmt.Sprintf("<@%s> NO SCHEDULED BLOCKS", config.Discord.UserID)
	}

	discord.NotifyChannel(config.Discord, scheduledBlocksMessage)
}
