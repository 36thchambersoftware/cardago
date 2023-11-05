package discord

type Config struct {
	AuthenticationToken string `yaml:"authenticationToken"`
	ServerID            string `yaml:"serverID"`
	UserID              string `yaml:"userID"`
	ChannelID           string `yaml:"channelID"`
	VoiceChannelID      string `yaml:"voiceChannelID"`
	WebhookURL          string `yaml:"webhookURL"`
}

// Requires one of content, file, embeds
type WebhookPayload struct {
	Content string `json:"content"` // REQUIRED the message contents (up to 2000 characters)
}
