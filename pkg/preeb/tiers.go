package preeb

import (
	"cardano/cardago/internal/log"

	_ "github.com/lib/pq"
)

// Lovelace in 1 ADA
const LovelaceToADA = 1000000

// Sequential epochs delegated to pool
const (
	EpochTier1 = 10
	EpochTier2 = 20
	EpochTier3 = 50
	EpochTier4 = 100
)

type EpochTier int

// Amount of ADA delegated to pool
const (
	LovelaceTier1 = 500000000
	LovelaceTier2 = 5000000000
	LovelaceTier3 = 50000000000
	LovelaceTier4 = 500000000000
	LovelaceTier5 = 1000000000000
)

type Delegator struct {
	Address          string
	StakeAddress     string
	LiveStake        uint64
	LiveStakeADA     float64
	EpochTotal       int32
	EpochLoyalty     int
	StakeLoyalty     uint64
	StakeLoyaltyADA  float64
	InitialAmount    uint64
	InitialAmountADA float64
	ActiveEpoch      int32
	CurrentEpoch     int32
}

func GetTiers(delegator *Delegator) {
	getEpochTier(delegator)
	getADATier(delegator)
}

func getEpochTier(delegator *Delegator) {
	if delegator.EpochTotal < EpochTier1 {
		return
	}

	if delegator.EpochTotal >= EpochTier1 {
		delegator.EpochLoyalty = EpochTier1
	}

	if delegator.EpochTotal >= EpochTier2 {
		delegator.EpochLoyalty = EpochTier2
	}

	if delegator.EpochTotal >= EpochTier3 {
		delegator.EpochLoyalty = EpochTier3
	}

	if delegator.EpochTotal >= EpochTier4 {
		delegator.EpochLoyalty = EpochTier4
	}
}

func getADATier(delegator *Delegator) {
	if delegator.LiveStake < LovelaceTier1 {
		return
	}

	if delegator.LiveStake >= LovelaceTier1 {
		delegator.StakeLoyalty = LovelaceTier1
	}

	if delegator.LiveStake >= LovelaceTier2 {
		delegator.StakeLoyalty = LovelaceTier2
	}

	if delegator.LiveStake >= LovelaceTier3 {
		delegator.StakeLoyalty = LovelaceTier3
	}

	if delegator.LiveStake >= LovelaceTier4 {
		delegator.StakeLoyalty = LovelaceTier4
	}

	if delegator.LiveStake >= LovelaceTier5 {
		delegator.StakeLoyalty = LovelaceTier5
	}

	delegator.StakeLoyaltyADA = float64(delegator.StakeLoyalty) / LovelaceToADA
}

func GetQualifiers(delegators map[string]*Delegator) map[string]*Delegator {
	qualifiedDelegators := make(map[string]*Delegator)
	for skey, d := range delegators {
		if d.EpochLoyalty > 0 && d.StakeLoyalty > 0 {
			qualifiedDelegators[skey] = d
		}
	}

	return qualifiedDelegators
}

type Stats struct {
	TotalEpochTier1    int `json:"totalEpochTier1"`
	TotalEpochTier2    int `json:"totalEpochTier2"`
	TotalEpochTier3    int `json:"totalEpochTier3"`
	TotalEpochTier4    int `json:"totalEpochTier4"`
	TotalLovelaceTier1 int `json:"totalLovelaceTier1"`
	TotalLovelaceTier2 int `json:"totalLovelaceTier2"`
	TotalLovelaceTier3 int `json:"totalLovelaceTier3"`
	TotalLovelaceTier4 int `json:"totalLovelaceTier4"`
	TotalLovelaceTier5 int `json:"totalLovelaceTier5"`
	TotalDelegators    int `json:"totalDelegators"`
}

func GetDelegatorStatistics(delegators map[string]Delegator) Stats {
	stats := Stats{
		TotalEpochTier1:    0,
		TotalEpochTier2:    0,
		TotalEpochTier3:    0,
		TotalEpochTier4:    0,
		TotalLovelaceTier1: 0,
		TotalLovelaceTier2: 0,
		TotalLovelaceTier3: 0,
		TotalLovelaceTier4: 0,
		TotalLovelaceTier5: 0,
		TotalDelegators:    len(delegators),
	}

	for _, delegator := range delegators {
		if delegator.LiveStake >= LovelaceTier1 && delegator.LiveStake < LovelaceTier2 {
			stats.TotalLovelaceTier1++
		}

		if delegator.LiveStake >= LovelaceTier2 && delegator.LiveStake < LovelaceTier3 {
			stats.TotalLovelaceTier2++
		}

		if delegator.LiveStake >= LovelaceTier3 && delegator.LiveStake < LovelaceTier4 {
			stats.TotalLovelaceTier3++
		}

		if delegator.LiveStake >= LovelaceTier4 && delegator.LiveStake < LovelaceTier5 {
			stats.TotalLovelaceTier4++
		}

		if delegator.LiveStake >= LovelaceTier5 {
			stats.TotalLovelaceTier5++
		}

		if delegator.EpochTotal >= EpochTier1 && delegator.EpochTotal < EpochTier2 {
			stats.TotalEpochTier1++
		}

		if delegator.EpochTotal >= EpochTier2 && delegator.EpochTotal < EpochTier3 {
			stats.TotalEpochTier2++
		}

		if delegator.EpochTotal >= EpochTier3 && delegator.EpochTotal < EpochTier4 {
			stats.TotalEpochTier3++
		}

		if delegator.EpochTotal >= EpochTier4 {
			stats.TotalEpochTier4++
		}
	}
	log.Infow("PREEB", "STATS", stats)
	return stats
}
