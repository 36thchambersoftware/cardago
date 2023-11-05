package discord

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

func ExecuteWebhook(config Config, content string) {
	// Check config url
	webhookURL, err := url.Parse(config.WebhookURL)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "DISCORD", "ERROR", err)
	}
	slog.Info("CARDAGO", "PACKAGE", "DISCORD", "INFO", webhookURL.String())

	// build post
	payload := WebhookPayload{
		Content: content,
	}

	message, err := json.Marshal(payload)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "DISCORD", "ERROR", err)
	}
	slog.Info("CARDAGO", "PACKAGE", "DISCORD", "INFO", string(message))

	client := &http.Client{}
	data := strings.NewReader(string(message))
	slog.Info("CARDAGO", "PACKAGE", "DISCORD", "INFO", data)

	// Create a new request
	request, err := http.NewRequest(http.MethodPost, webhookURL.String(), data)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "DISCORD", "ERROR", err)
	}
	slog.Info("CARDAGO", "PACKAGE", "DISCORD", "INFO", request)

	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "DISCORD", "ERROR", err)
	}

	if response.StatusCode != http.StatusNoContent {
		slog.Info("CARDAGO", "PACKAGE", "DISCORD", "INFO", "You are not waiting for a response. Add ?wait=true to webhook url")
	}

	if response.StatusCode != http.StatusOK {
		slog.Error("CARDAGO", "PACKAGE", "DISCORD", "ERROR", response.StatusCode)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "DISCORD", "ERROR", err)
	}
	slog.Info("CARDAGO", "PACKAGE", "DISCORD", "INFO", body)
}
