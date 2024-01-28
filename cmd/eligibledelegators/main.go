package main

import (
	"context"
	"time"

	"cardano/cardago/internal/config"
	"cardano/cardago/internal/log"
	"cardano/cardago/pkg/blockfrost"

	bfg "github.com/blockfrost/blockfrost-go"
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

	bfc := bfg.NewAPIClient(bfg.APIClientOptions{ProjectID: cfg.Blockfrost.ProjectID})

	delegatorAddresses, err := blockfrost.GetDelegatorsByPoolID(ctx, bfc, cfg.Cardano.StakePoolID)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "BLOCKFROST", "ERROR", err)
		return
	}

	ep, err := bfc.EpochLatest(ctx)
	if err != nil {
		log.Errorw("BLOCKFROST", "GET", "EpochLatest", "ERROR", err)
	}

	for _, a := range delegatorAddresses {
		blockfrost.StakeAddressHistoryByEpoch(ctx, bfc, a, cfg.Cardano.Bech32PoolId, ep.Epoch)
	}

	log.Infow("delegator", "addresses", delegatorAddresses)
}
