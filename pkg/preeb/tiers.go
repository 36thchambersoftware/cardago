package preeb

// Lovelace in 1 ADA
const LovelaceToADA = 1000000

// Sequential epochs delegated to pool
const (
	EpochTier10  = 10
	EpochTier20  = 20
	EpochTier50  = 50
	EpochTier100 = 100
)

type EpochTier int

// Amount of ADA delegated to pool
const (
	LovelaceTier500     = 500000000
	LovelaceTier5000    = 5000000000
	LovelaceTier50000   = 50000000000
	LovelaceTier500000  = 500000000000
	LovelaceTier1000000 = 1000000000000
	LovelaceTier5000000 = 5000000000000
)

type Delegator struct {
	LiveStake    uint64
	Address      string
	StakeAddress string
	EpochLoyalty int
	StakeLoyalty uint64
	ActiveEpoch  int32
	CurrentEpoch int32
	// Rewards for given epoch in Lovelaces
	InitialAmount uint64
	// Bech32 ID of pool being delegated to
	PoolID string
}

func GetTiers(delegator *Delegator) {
	getEpochTier(delegator)
	getADATier(delegator)
}

func getEpochTier(delegator *Delegator) {
	epochsDelegated := delegator.CurrentEpoch - delegator.ActiveEpoch
	if epochsDelegated < 10 {
		return
	}

	if epochsDelegated > EpochTier10 {
		delegator.EpochLoyalty = EpochTier10
	}

	if epochsDelegated > EpochTier20 {
		delegator.EpochLoyalty = EpochTier20
	}

	if epochsDelegated > EpochTier50 {
		delegator.EpochLoyalty = EpochTier50
	}

	if epochsDelegated > EpochTier100 {
		delegator.EpochLoyalty = EpochTier100
	}
}

func getADATier(delegator *Delegator) {
	if delegator.LiveStake < LovelaceTier500 {
		return
	}

	if delegator.LiveStake > LovelaceTier500 {
		delegator.StakeLoyalty = LovelaceTier500
	}

	if delegator.LiveStake > LovelaceTier5000 {
		delegator.StakeLoyalty = LovelaceTier5000
	}

	if delegator.LiveStake > LovelaceTier50000 {
		delegator.StakeLoyalty = LovelaceTier50000
	}

	if delegator.LiveStake > LovelaceTier500000 {
		delegator.StakeLoyalty = LovelaceTier500000
	}

	if delegator.LiveStake > LovelaceTier1000000 {
		delegator.StakeLoyalty = LovelaceTier1000000
	}

	if delegator.LiveStake > LovelaceTier5000000 {
		delegator.StakeLoyalty = LovelaceTier5000000
	}
}
