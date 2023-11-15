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
	config := config.Get()
	err := config.LoadConfig()
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CONFIG", "ERROR", err)
		return
	}
	log.Debugw("CARDAGO", "PACKAGE", "CONFIG", "RUNTIME", config)

	KESPeriodInfo, err := cardano.GetKESPeriodInfo(config.Cardano.NodeCertPath)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
		return
	}

	KESExpiryDate, err := time.Parse("2006-01-02T15:04:05Z", KESPeriodInfo.KesKesKeyExpiry)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
	}
	log.Infow("CARDAGO", "PACKAGE", "CARDANO", "KES EXPIRY", KESExpiryDate.Format("2006-01-02 15:04:05"))

	KESExpiryMessage := fmt.Sprintf("<@%s> KES Expiry Date: %s", config.Discord.UserID, KESExpiryDate.Format("2006-01-02 15:04:05"))

	discord.ExecuteWebhook(config.Discord, KESExpiryMessage)
}
