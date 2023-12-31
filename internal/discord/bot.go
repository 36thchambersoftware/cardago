package discord

import (
	"time"

	"cardano/cardago/internal/log"

	"github.com/bwmarrin/discordgo"
)

func NotifyChannel(config Config, content string) {
	discord, err := discordgo.New(config.AuthenticationToken)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "DISCORD", "ERROR", err)
	}

	message := discordgo.MessageSend{
		Content: content,
	}
	response, err := discord.ChannelMessageSendComplex(config.ChannelID, &message)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "DISCORD", "ChannelMessageSend", err)
	}
	log.Infow("CARDAGO", "PACKAGE", "DISCORD", "MESSAGE", response)
}

func ScheduleEvent(config Config, name string, start time.Time, end time.Time, content string) {
	discord, err := discordgo.New(config.AuthenticationToken)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "DISCORD", "ERROR", err)
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
		log.Errorw("CARDAGO", "PACKAGE", "DISCORD", "GuildScheduledEventParams", err)
	}
	log.Infow("CARDAGO", "PACKAGE", "DISCORD", "MESSAGE", st)
}
