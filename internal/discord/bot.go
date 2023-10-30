package discord

import (
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
