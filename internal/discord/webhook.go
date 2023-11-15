package discord

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"cardano/cardago/internal/log"
)

func ExecuteWebhook(config Config, content string) {
	// Check config url
	webhookURL, err := url.Parse(config.WebhookURL)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "DISCORD", "ERROR", err)
	}
	log.Debugw("CARDAGO", "PACKAGE", "DISCORD", "webhook", webhookURL.String())

	// build post
	payload := WebhookPayload{
		Content: content,
	}

	message, err := json.Marshal(payload)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "DISCORD", "marshal", err)
	}
	log.Debugw("CARDAGO", "PACKAGE", "DISCORD", "message", string(message))

	client := &http.Client{}
	data := strings.NewReader(string(message))

	// Create a new request
	request, err := http.NewRequest(http.MethodPost, webhookURL.String(), data)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "DISCORD", "request", err)
	}

	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "DISCORD", "response", err)
	}

	if response.StatusCode == http.StatusNoContent {
		log.Warnw("CARDAGO", "PACKAGE", "DISCORD", "INFO", "You are not waiting for a response. Add ?wait=true to webhook url")
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNoContent {
		log.Errorw("CARDAGO", "PACKAGE", "DISCORD", "ERROR", response.StatusCode)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "DISCORD", "body", body, "error", err)
	}
	log.Debugw("CARDAGO", "PACKAGE", "DISCORD", "body", string(body))
	log.Infow("CARDAGO", "PACKAGE", "DISCORD", "success", "message sent")
}
