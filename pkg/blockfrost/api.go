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

func PoolDelegators(ctx context.Context, bfc bfg.APIClient, poolID string) ([]bfg.PoolDelegator, error) {
	count := 100
	results := []bfg.PoolDelegator{}
	page := 1

	for count == 100 {
		result, err := bfc.PoolDelegators(ctx, poolID, bfg.APIQueryParams{
			Page: page,
		})
		if err != nil {
			// log.Errorw("CARDAGO", "PACKAGE", "BLOCKFROST", "ERROR", err)
			return nil, err
		}
		// log.Infow("BLOCKFROST", "GET", "DELEGATORS", "RESULT", result)
		results = append(results, result...)
		count = len(result)
		page++
		// log.Infow("BLOCKFROST", "GET", "DELEGATORS", "COUNT", len(result))
	}
	// log.Infow("BLOCKFROST", "GET", "DELEGATORS", "DATA", results)
	return results, nil
}

func GetDelegatorsByPoolID(ctx context.Context, bfc bfg.APIClient, poolID string) (map[string]*preeb.Delegator, error) {
	result, err := PoolDelegators(ctx, bfc, poolID)
	if err != nil {
		return nil, err
	}

	// log.Infow("BLOCKFROST", "GET", "DELEGATORS", "COUNT", len(result))

	data := make(map[string]*preeb.Delegator)

	for _, v := range result {
		delegator := preeb.Delegator{}
		delegator.StakeAddress = v.Address
		liveStake, err := strconv.ParseUint(v.LiveStake, 10, 64)
		if err != nil {
			return nil, err
		}
		delegator.LiveStake = liveStake
		delegator.LiveStakeADA = float64(delegator.LiveStake) / preeb.LovelaceToADA

		addresses, err := bfc.AccountAssociatedAddresses(ctx, v.Address, bfg.APIQueryParams{})
		if err != nil {
			return nil, err
		}

		delegator.Address = addresses[0].Address
		data[v.Address] = &delegator
	}

	// log.Infow("DATA", "DELEGATORS", data)

	return data, nil
}

func StakeAddressHistoryByEpoch(ctx context.Context, bfc bfg.APIClient, delegator *preeb.Delegator, poolID string, currentEpoch int) error {
	delegationHistory, err := bfc.AccountDelegationHistory(ctx, string(delegator.StakeAddress), bfg.APIQueryParams{Order: "desc"})
	if err != nil {
		log.Errorw("BLOCKFROST", "GET", "AccountDelegation History", "ERROR", err)
		return err
	}

	// log.Infow("CARDAGO", "PACKAGE", "BLOCKFROST", "DELEGATOR", delegator, "HISTORY", delegationHistory)

	// take the first one - we got these addresses from known delegators, so 0 always exists.
	if delegationHistory[0].PoolID == poolID {
		amount, err := strconv.ParseUint(delegationHistory[0].Amount, 10, 64)
		if err != nil {
			return err
		}

		delegator.ActiveEpoch = delegationHistory[0].ActiveEpoch
		delegator.CurrentEpoch = int32(currentEpoch)
		delegator.EpochTotal = delegator.CurrentEpoch - delegator.ActiveEpoch
		delegator.InitialAmount = amount
		delegator.InitialAmountADA = float64(delegator.InitialAmount) / preeb.LovelaceToADA

		preeb.GetTiers(delegator)
	}

	return nil
}

/**
https://cardano-mainnet.blockfrost.io/api/v0/assets/2d73c7107e14b8cae9efa6d9794f1f4db5f0a5d23ad3147cb4c3b7345052454542353030412d3130452d30303031
Asset:
{
    "asset": "2d73c7107e14b8cae9efa6d9794f1f4db5f0a5d23ad3147cb4c3b7345052454542353030412d3130452d30303031",
    "policy_id": "2d73c7107e14b8cae9efa6d9794f1f4db5f0a5d23ad3147cb4c3b734",
    "asset_name": "5052454542353030412d3130452d30303031",
    "fingerprint": "asset12k454zzldch6n38aph8q0kpcurqcczh06yw2u9",
    "quantity": "1",
    "initial_mint_tx_hash": "8db2a8df36b20f048ed5fe5b29d70b819b48972560468f8ad79dab80d29954ce",
    "mint_or_burn_count": 1,
    "onchain_metadata": {
        "name": "500A-10E-0001",
        "files": [
            {
                "src": "ipfs://Qmat3Fnqe6M1pqM3KM5EAAHZN6d4tKBb1R7XGY7kwREkPW",
                "name": "500A-10E-0001",
                "mediaType": "image/png"
            }
        ],
        "image": "ipfs://Qmat3Fnqe6M1pqM3KM5EAAHZN6d4tKBb1R7XGY7kwREkPW",
        "Discord": "thepriebe",
        "Twitter": "PREEBPool",
        "Website": "https://preeb.cloud",
        "mediaType": "image/png",
        "description": "500A-10E"
    },
    "onchain_metadata_standard": "CIP25v1",
    "onchain_metadata_extra": null,
    "metadata": null
}
*/

/**

1. Get all assets from all policies
https://cardano-mainnet.blockfrost.io/api/v0/assets/policy/2d73c7107e14b8cae9efa6d9794f1f4db5f0a5d23ad3147cb4c3b734
policy: []assets
[{"asset":"2d73c7107e14b8cae9efa6d9794f1f4db5f0a5d23ad3147cb4c3b7345052454542353030412d3130452d30303031","quantity":"1"},{"asset":"2d73c7107e14b8cae9efa6d9794f1f4db5f0a5d23ad3147cb4c3b7345052454542353030412d3130452d30303032","quantity":"1"}]
policy1: asset1, asset2, assetn
policy2: asset1, asset2, assetn
policyn: assetn

2. Get asset's addresses
https://cardano-mainnet.blockfrost.io/api/v0/assets/2d73c7107e14b8cae9efa6d9794f1f4db5f0a5d23ad3147cb4c3b7345052454542353030412d3130452d30303031/addresses
[{"address":"addr1","quantity":"1"},{"address":"addr2","quantity":"1"}]

3.

*/

// func getAddressesByAssets(ctx context.Context, bfc bfg.APIClient, assets []string) {
// 	assetsByAddress := []string{}
// 	for i, a := range assets {
// 		assetAddresses := bfc.AssetAddresses(ctx, a)
// 	}

// }
