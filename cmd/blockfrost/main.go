package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cardano/cardago/internal/config"
	"cardano/cardago/internal/discord"
	"cardano/cardago/internal/log"
	"cardano/cardago/pkg/blockfrost"
)

func main() {
	cfg := config.Get()
	err := cfg.LoadConfig()
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CONFIG", "ERROR", err)
		return
	}
	log.Debugw("CARDAGO", "PACKAGE", "CONFIG", "RUNTIME", cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	delegatorAddresses, err := blockfrost.GetDelegatorsByAddress(ctx, cfg.Blockfrost, cfg.Cardano.StakePoolID)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "BLOCKFROST", "ERROR", err)
		return
	}

	log.Infow("delegator", "addresses", delegatorAddresses)

	// Group addresses by tens to not hit the max 2000 character limit in discord, and hopefully avoid rate limiting.
	count := len(delegatorAddresses)
	discord.ExecuteWebhook(cfg.Discord, fmt.Sprintf("There are %v unique delegators", count))

	groups := []string{}
	remaining := count
	for i := 0; i < count; i += 10 {
		if count > i+10 {
			groups = append(groups, strings.Join(delegatorAddresses[i:i+10], ", "))
			remaining -= 10
		} else {
			groups = append(groups, strings.Join(delegatorAddresses[i:i+remaining], ", "))
		}
	}

	log.Infow("groups", "addresses", groups)
	for _, group := range groups {
		discord.ExecuteWebhook(cfg.Discord, group)
	}
}
