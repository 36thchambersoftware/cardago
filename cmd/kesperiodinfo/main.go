package main

import (
	"fmt"
	"time"

	"cardano/cardago/internal/cardano"
	"cardano/cardago/internal/config"
	"cardano/cardago/internal/discord"
	"cardano/cardago/internal/log"
)

func main() {
	logger := log.InitializeLogger()
	config := config.Get()
	err := config.LoadConfig()
	if err != nil {
		logger.Errorw("CARDAGO", "PACKAGE", "CONFIG", "ERROR", err)
		return
	}
	logger.Infow("CARDAGO", "PACKAGE", "CONFIG", "RUNTIME", config)

	KESPeriodInfo, err := cardano.GetKESPeriodInfo(config.Cardano.NodeCertPath)
	if err != nil {
		logger.Errorw("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
		return
	}

	KESExpiryDate, err := time.Parse("2006-01-02T15:04:05Z", KESPeriodInfo.KesKesKeyExpiry)
	if err != nil {
		logger.Errorw("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
	}
	logger.Infow("CARDAGO", "PACKAGE", "CARDANO", "EXPIRY", KESExpiryDate)

	KESExpiryMessage := fmt.Sprintf("<@%s> KES Expiry Date: %s", config.Discord.UserID, KESExpiryDate.Format("2006-01-02 15:04:05"))

	discord.ExecuteWebhook(config.Discord, KESExpiryMessage)
}
