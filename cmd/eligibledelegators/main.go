package main

import (
	"context"
	"time"

	"cardano/cardago/internal/config"
	"cardano/cardago/internal/log"
	"cardano/cardago/pkg/blockfrost"
	"cardano/cardago/pkg/preeb"

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

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	bfc := bfg.NewAPIClient(bfg.APIClientOptions{ProjectID: cfg.Blockfrost.ProjectID})

	delegatorAddresses, err := blockfrost.GetDelegatorsByPoolID(ctx, bfc, cfg.Cardano.StakePoolID)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "BLOCKFROST", "ERROR", err)
		return
	}

	// log.Infow("BLOCKFROST", "GET", "DELEGATORS", "DATA", delegatorAddresses)

	ep, err := bfc.EpochLatest(ctx)
	if err != nil {
		log.Errorw("BLOCKFROST", "GET", "EpochLatest", "ERROR", err)
	}

	log.Infow("CARDAGO", "PACKAGE", "BLOCKFROST", "DELEGATORS", "BEGIN PROCESSING")
	for _, a := range delegatorAddresses {
		blockfrost.StakeAddressHistoryByEpoch(ctx, bfc, a, cfg.Cardano.Bech32PoolId, ep.Epoch)
	}
	log.Infow("CARDAGO", "PACKAGE", "BLOCKFROST", "DELEGATORS", "END PROCESSING")
	// log.Infow("BLOCKFROST", "GET", "DELEGATORS", "DATA", delegatorAddresses)

	qualifiedDelegators := preeb.GetQualifiers(delegatorAddresses)

	log.Infow("delegator", "addresses", qualifiedDelegators)
}
