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
	config.LoadConfig()
	slog.Info("CARDAGO", "PACKAGE", "CONFIG", "RUNTIME", config)

	KESPeriodInfo := cardano.GetKESPeriodInfo("mainnet", config.NodeCertPath)

	KESExpiryDate, err := time.Parse("2006-01-02T15:04:05Z", KESPeriodInfo.KesKesKeyExpiry)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
	}
	slog.Info("CARDAGO", "PACKAGE", "CARDANO", "EXPIRY", KESExpiryDate)

	KESExpiryMessage := fmt.Sprintf("<@%s> KES Expiry Date: %s", config.Discord.UserID, KESExpiryDate.Format("2006-01-02 15:04:05"))

	discord.NotifyChannel(config.Discord, KESExpiryMessage)
}
