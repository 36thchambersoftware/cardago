package discord

type Config struct {
	AuthenticationToken string `yaml:"authenticationToken,omitempty"`
	ServerID            string `yaml:"serverID,omitempty"`
	UserID              string `yaml:"userID"`
	ChannelID           string `yaml:"channelID,omitempty"`
	VoiceChannelID      string `yaml:"voiceChannelID,omitempty"`
	WebhookURL          string `yaml:"webhookURL"`
}

// Requires one of content, file, embeds
type WebhookPayload struct {
	Content string `json:"content"` // REQUIRED the message contents (up to 2000 characters)
}
