package main

import (
	"errors"
	"fmt"
	"os"

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
	// Check if the Cardano Shelley genesis file path exists.
	_, err = os.Stat(cfg.Cardano.ShelleyGenesisFilePath)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CONFIG", "Cardano.ShelleyGenesisFilePath", err)
		return
	}

	// Check if the Cardano VRFS key file path exists.
	_, err = os.Stat(cfg.Cardano.VRFSKeyFilePath)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CONFIG", "Cardano.VRFSKeyFilePath", err)
		return
	}
	log.Debugw("CARDAGO", "PACKAGE", "CONFIG", "RUNTIME", cfg)

	// Check if the leader log directory exists.
	_, err = os.Stat(cfg.Cardano.Leader.Directory)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CONFIG", "Cardano.Leader.Directory", err)
		log.Errorw("CARDAGO", "PACKAGE", "CONFIG", "Cardano.Leader.Directory", cfg.Leader)

		return
	}

	// Get the scheduled Cardano blocks.
	scheduledBlocks, err := cardano.GetScheduledBlocks(cfg.Cardano)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
		if errors.Is(err, cardano.ErrorEpochAlreadyChecked) {
			data, readErr := cardano.ReadScheduledBlocks(cfg.Cardano.Leader)
			if readErr != nil {
				log.Errorw("CARDAGO", "PACKAGE", "CARDANO", "ERROR", readErr)
				return
			}

			discord.ExecuteWebhook(cfg.Discord, string(data))
		}
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
