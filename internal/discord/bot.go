package discord

import (
	"fmt"
	"log/slog"
	"time"

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

func NotifyChannel(config DiscordConfig, content string) {
	discord, err := discordgo.New(config.AuthenticationToken)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "DISCORD", "ERROR", err)
	}

	message := discordgo.MessageSend{
		Content: content,
	}
	response, err := discord.ChannelMessageSendComplex(config.ChannelID, &message)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "DISCORD", "ChannelMessageSend", err)
	}
	slog.Info("CARDAGO", "PACKAGE", "DISCORD", "MESSAGE", response)
}

func ScheduleEvent(config DiscordConfig, name string, start time.Time, end time.Time, content string) {
	discord, err := discordgo.New(config.AuthenticationToken)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "DISCORD", "ERROR", err)
	}

	startDate := start
	endDate := end
	event := discordgo.GuildScheduledEventParams{
		ChannelID:          config.VoiceChannelID,
		Name:               name,
		Description:        content,
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
