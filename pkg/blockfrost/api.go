package blockfrost

import (
	"context"

	"cardano/cardago/internal/log"

	bfg "github.com/blockfrost/blockfrost-go"
)

type Config struct {
	ProjectID string `yaml:"projectid"`
	APIURL    string `yaml:"apiurl"`
	Timeout   string `yaml:"timeout"`
}

func GetDelegatorsByAddress(ctx context.Context, cfg Config, poolID string) ([]string, error) {
	bfc := bfg.NewAPIClient(bfg.APIClientOptions{ProjectID: cfg.ProjectID})
	delegators, err := bfc.PoolDelegators(ctx, poolID, bfg.APIQueryParams{})
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "BLOCKFROST", "ERROR", err)
		return nil, err
	}

	address := []string{}

	for _, delegator := range delegators {
		addresses, err := bfc.AccountAssociatedAddresses(ctx, delegator.Address, bfg.APIQueryParams{})
		if err != nil {
			return nil, err
		}

		for _, addr := range addresses {
			address = append(address, addr.Address)
			break
		}
	}

	return address, nil
}
