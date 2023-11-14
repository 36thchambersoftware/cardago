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
	logger := log.InitializeLogger()

	// Check config url
	webhookURL, err := url.Parse(config.WebhookURL)
	if err != nil {
		logger.Errorw("CARDAGO", "PACKAGE", "DISCORD", "ERROR", err)
	}
	logger.Infow("CARDAGO", "PACKAGE", "DISCORD", "webhook", webhookURL.String())

	// build post
	payload := WebhookPayload{
		Content: content,
	}

	message, err := json.Marshal(payload)
	if err != nil {
		logger.Errorw("CARDAGO", "PACKAGE", "DISCORD", "marshal", err)
	}
	logger.Infow("CARDAGO", "PACKAGE", "DISCORD", "message", string(message))

	client := &http.Client{}
	data := strings.NewReader(string(message))

	// Create a new request
	request, err := http.NewRequest(http.MethodPost, webhookURL.String(), data)
	if err != nil {
		logger.Errorw("CARDAGO", "PACKAGE", "DISCORD", "request", err)
	}

	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		logger.Errorw("CARDAGO", "PACKAGE", "DISCORD", "response", err)
	}

	if response.StatusCode == http.StatusNoContent {
		logger.Infow("CARDAGO", "PACKAGE", "DISCORD", "INFO", "You are not waiting for a response. Add ?wait=true to webhook url")
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNoContent {
		logger.Errorw("CARDAGO", "PACKAGE", "DISCORD", "ERROR", response.StatusCode)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Errorw("CARDAGO", "PACKAGE", "DISCORD", "body", body, "error", err)
	}
	logger.Infow("CARDAGO", "PACKAGE", "DISCORD", "success", body)
}
