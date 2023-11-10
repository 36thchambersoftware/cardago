package main

import (
	"fmt"
	"log/slog"

	"cardano/cardago/internal/cardano"
	"cardano/cardago/internal/config"
	"cardano/cardago/internal/discord"
)

/**
* Retrieves the scheduled Cardano blocks and sends a Discord message to the specified user with the list of scheduled blocks.
 */
func main() {
	// Get the configuration from the configuration file.
	config := config.Get()
	err := config.LoadConfig()
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CONFIG", "ERROR", err)
		return
	}
	slog.Info("CARDAGO", "PACKAGE", "CONFIG", "RUNTIME", config)

	// Get the scheduled Cardano blocks.
	scheduledBlocks, err := cardano.GetScheduledBlocks(config.Cardano, config.Logs)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
		return
	}

	scheduledBlocksText := fmt.Sprintln(scheduledBlocks)
	// Create a Discord message with the list of scheduled blocks.
	scheduledBlocksMessage := fmt.Sprintf("<@%s> SCHEDULED BLOCKS: %s", config.Discord.UserID, scheduledBlocksText)
	if len(scheduledBlocks) == 0 {
		scheduledBlocksMessage = fmt.Sprintf("<@%s> NO SCHEDULED BLOCKS", config.Discord.UserID)
	}

	// Execute the Discord webhook to send the message.
	discord.ExecuteWebhook(config.Discord, scheduledBlocksMessage)
}
