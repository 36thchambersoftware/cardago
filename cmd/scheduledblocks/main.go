package main

import (
	"fmt"

	"cardano/cardago/internal/cardano"
	"cardano/cardago/internal/config"
	"cardano/cardago/internal/discord"
	"cardano/cardago/internal/log"
)

/**
* Retrieves the scheduled Cardano blocks and sends a Discord message to the specified user with the list of scheduled blocks.
 */
func main() {
	// Get the configuration from the configuration file.
	cfg := config.Get()
	err := cfg.LoadConfig()
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CONFIG", "ERROR", err)
		return
	}
	log.Debugw("CARDAGO", "PACKAGE", "CONFIG", "RUNTIME", cfg)

	// Get the scheduled Cardano blocks.
	scheduledBlocks, err := cardano.GetScheduledBlocks(cfg.Cardano)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
		return
	}

	scheduledBlocksText := fmt.Sprintln(scheduledBlocks)
	// Create a Discord message with the list of scheduled blocks.
	scheduledBlocksMessage := fmt.Sprintf("<@%s> SCHEDULED BLOCKS: %s", cfg.Discord.UserID, scheduledBlocksText)
	if len(scheduledBlocks) == 0 {
		scheduledBlocksMessage = fmt.Sprintf("<@%s> NO SCHEDULED BLOCKS", cfg.Discord.UserID)
	}

	// Execute the Discord webhook to send the message.
	discord.ExecuteWebhook(cfg.Discord, scheduledBlocksMessage)
}
