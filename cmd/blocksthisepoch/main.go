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

	scheduledBlocks, err := cardano.GetScheduledBlocks(config.Cardano, config.Logs)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
		return
	}

	scheduledBlocksText := fmt.Sprintln(scheduledBlocks)

	scheduledBlocksMessage := fmt.Sprintf("<@%s> SCHEDULED BLOCKS: %s", config.Discord.UserID, scheduledBlocksText)
	if len(scheduledBlocks) == 0 {
		scheduledBlocksMessage = fmt.Sprintf("<@%s> NO SCHEDULED BLOCKS", config.Discord.UserID)
	}

	discord.ExecuteWebhook(config.Discord, scheduledBlocksMessage)
}
