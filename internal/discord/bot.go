package discord

import (
	"fmt"
	"log/slog"
	"time"

	"cardano/cardago/internal/cardano"

	"github.com/bwmarrin/discordgo"
)

type DiscordConfig struct {
	AuthenticationToken string `yaml:"authenticationToken"`
	ServerID            string `yaml:"serverID"`
	UserID              string `yaml:"userID"`
	ChannelID           string `yaml:"channelID"`
	VoiceChannelID      string `yaml:"voiceChannelID"`
}

func NotifyKESExpiryDate(config DiscordConfig, KESExpiryDate time.Time) {
	discord, err := discordgo.New(config.AuthenticationToken)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "DISCORD", "ERROR", err)
	}
	KESExpiryMessage := fmt.Sprintf("<@%s> KES Expiry Date: %s", config.UserID, KESExpiryDate.Format("2006-01-02 15:04:05"))
	message := discordgo.MessageSend{
		Content: KESExpiryMessage,
	}
	response, err := discord.ChannelMessageSendComplex(config.ChannelID, &message)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "DISCORD", "ChannelMessageSend", err)
	}
	slog.Info("CARDAGO", "PACKAGE", "DISCORD", "MESSAGE", response)

	startDate := KESExpiryDate
	endDate := KESExpiryDate.AddDate(0, 0, 1)
	event := discordgo.GuildScheduledEventParams{
		ChannelID:          config.VoiceChannelID,
		Name:               "kes-expiry-date",
		Description:        KESExpiryMessage,
		ScheduledStartTime: &startDate,
		ScheduledEndTime:   &endDate,
		PrivacyLevel:       discordgo.GuildScheduledEventPrivacyLevelGuildOnly,
		Status:             discordgo.GuildScheduledEventStatusScheduled,
		EntityType:         discordgo.GuildScheduledEventEntityTypeVoice,
	}

	st, err := discord.GuildScheduledEventCreate(config.ServerID, &event)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "DISCORD", "GuildScheduledEventParams", err)
	}
	slog.Info("CARDAGO", "PACKAGE", "DISCORD", "MESSAGE", st)
}

func NotifyScheduledBlocks(config DiscordConfig, scheduledBlocks []cardano.ScheduledBlock) {
	discord, err := discordgo.New(config.AuthenticationToken)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "DISCORD", "ERROR", err)
	}

	scheduledBlocksMessage := fmt.Sprintf("<@%s> SCHEDULED BLOCKS: %s", config.UserID, scheduledBlocks)
	if len(scheduledBlocks) == 0 {
		scheduledBlocksMessage = fmt.Sprintf("<@%s> NO SCHEDULED BLOCKS", config.UserID)
	}

	message := discordgo.MessageSend{
		Content: scheduledBlocksMessage,
	}
	response, err := discord.ChannelMessageSendComplex(config.ChannelID, &message)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "DISCORD", "ChannelMessageSend", err)
	}
	slog.Info("CARDAGO", "PACKAGE", "DISCORD", "MESSAGE", response)
}
