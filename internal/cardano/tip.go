package cardano

import (
	"encoding/json"

	"cardano/cardago/internal/log"
)

type Tip struct {
	Era             string `json:"era,omitempty"`
	Hash            string `json:"hash,omitempty"`
	SyncProgress    string `json:"syncProgress,omitempty"`
	Block           int    `json:"block,omitempty"`
	Epoch           int    `json:"epoch,omitempty"`
	Slot            int    `json:"slot,omitempty"`
	SlotInEpoch     int    `json:"slotInEpoch,omitempty"`
	SlotsToEpochEnd int    `json:"slotsToEpochEnd,omitempty"`
}

func QueryTip() (Tip, error) {
	log.Infow("CARDAGO", "PACKAGE", "CARDANO", "ACTION", "QueryTip")
	tip := Tip{}

	args := []string{
		"conway", "query",
		"tip",
		"--mainnet",
	}
	log.Debugw("CARDAGO", "PACKAGE", "CARDANO", "ARGS", args)

	output, err := Run(args)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CARDANO", "TYPE", "query tip", "ERROR", err, "OUTPUT", string(output))
		return tip, err
	}

	err = json.Unmarshal(output, &tip)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
	}
	log.Infow("CARDAGO", "PACKAGE", "CARDANO", "TIP", tip)

	return tip, err
}
