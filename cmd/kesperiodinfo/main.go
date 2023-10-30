package main

import (
	"fmt"
	"log/slog"
	"time"

	"cardano/cardago/internal/cardano"
	"cardano/cardago/internal/config"
	"cardano/cardago/internal/discord"
)

func main() {
	config := config.Get()
	err := config.LoadConfig()
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CONFIG", "ERROR", err)
		return
	}
	slog.Info("CARDAGO", "PACKAGE", "CONFIG", "RUNTIME", config)

	KESPeriodInfo, err := cardano.GetKESPeriodInfo("mainnet", config.NodeCertPath)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
		return
	}

	KESExpiryDate, err := time.Parse("2006-01-02T15:04:05Z", KESPeriodInfo.KesKesKeyExpiry)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
	}
	slog.Info("CARDAGO", "PACKAGE", "CARDANO", "EXPIRY", KESExpiryDate)

	KESExpiryMessage := fmt.Sprintf("<@%s> KES Expiry Date: %s", config.Discord.UserID, KESExpiryDate.Format("2006-01-02 15:04:05"))

	discord.NotifyChannel(config.Discord, KESExpiryMessage)
}
