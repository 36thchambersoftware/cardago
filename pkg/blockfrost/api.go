package blockfrost

import (
	"context"
	"strconv"

	"cardano/cardago/internal/log"
	"cardano/cardago/pkg/preeb"

	bfg "github.com/blockfrost/blockfrost-go"
)

type Config struct {
	ProjectID string `yaml:"projectid"`
	APIURL    string `yaml:"apiurl"`
	Timeout   string `yaml:"timeout"`
}

func GetDelegatorsByPoolID(ctx context.Context, bfc bfg.APIClient, poolID string) ([]*preeb.Delegator, error) {
	result, err := bfc.PoolDelegators(ctx, poolID, bfg.APIQueryParams{})
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "BLOCKFROST", "ERROR", err)
		return nil, err
	}

	delegators := []*preeb.Delegator{}

	for _, v := range result {
		delegator := preeb.Delegator{}
		delegator.StakeAddress = v.Address
		v, err := strconv.ParseUint(v.LiveStake, 10, 64)
		if err != nil {
			return nil, err
		}
		delegator.LiveStake = v

		// addresses, err := bfc.AccountAssociatedAddresses(ctx, v.Address, bfg.APIQueryParams{})
		// if err != nil {
		// 	return nil, err
		// }

		// delegator.Address = addresses[0].Address
		delegators = append(delegators, &delegator)
	}

	return delegators, nil
}

type (
	Address string
	Epoch   int
)

func StakeAddressHistoryByEpoch(ctx context.Context, bfc bfg.APIClient, delegator *preeb.Delegator, poolID string, currentEpoch int) error {
	delegationHistory, err := bfc.AccountDelegationHistory(ctx, string(delegator.StakeAddress), bfg.APIQueryParams{Order: "desc"})
	if err != nil {
		log.Errorw("BLOCKFROST", "GET", "AccountDelegation History", "ERROR", err)
		return err
	}

	// take the first one - we got these addresses from known delegators, so 0 always exists.
	if delegationHistory[0].PoolID == poolID {
		amount, err := strconv.ParseUint(delegationHistory[0].Amount, 10, 64)
		if err != nil {
			return err
		}

		delegator.ActiveEpoch = delegationHistory[0].ActiveEpoch
		delegator.CurrentEpoch = int32(currentEpoch)
		delegator.InitialAmount = amount

		if (delegator.CurrentEpoch - delegator.ActiveEpoch) > 10 {
			preeb.GetTiers(delegator)
		}
	}

	return nil
}
