package cardano

import (
	"encoding/json"
	"log/slog"
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
	tip := Tip{}

	args := []string{
		"query",
		"tip",
		"--mainnet",
	}
	slog.Info("CARDAGO", "PACKAGE", "CARDANO", "ARGS", args)

	output, err := Run(args)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CARDANO", "TYPE", "query tip", "ERROR", err, "OUTPUT", output)
		return tip, err
	}

	err = json.Unmarshal(output, &tip)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
	}
	slog.Info("CARDAGO", "PACKAGE", "CARDANO", "TIP", tip)

	return tip, err
}
