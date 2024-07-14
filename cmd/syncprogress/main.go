package main

import (
	"fmt"

	"cardano/cardago/internal/cardano"
	"cardano/cardago/internal/config"
	"cardano/cardago/internal/discord"
	"cardano/cardago/internal/log"

	"github.com/bwmarrin/discordgo"
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

	message := fmt.Sprintf("SYNC PROGRESS: %s", tip.SyncProgress)

	preeb, err := discordgo.New("")
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
		return
	}
	log.Infow("CARDAGO", "PACKAGE", "CARDANO", "DISCORD", preeb.State)

	data := discordgo.ChannelEdit{
		Name: message,
	}

	log.Infow("CARDAGO", "PACKAGE", "CARDANO", "data", data)

	channel, err := preeb.ChannelEdit("1262933138896846900", &data, nil)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
		return
	}

	log.Infow("CARDAGO", "PACKAGE", "CARDANO", "channel", channel)

	if channel.Name != message {
		log.Errorw("CARDAGO", "PACKAGE", "CARDANO", "ERROR", "Could not set name")
	}

	discord.ExecuteWebhook(config.Discord, message)
}
