package cardano

import (
	"encoding/json"
	"log/slog"
	"os/exec"
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

func QueryTip() Tip {
	tip := Tip{}

	output, err := exec.Command("cardano-cli", "query", "tip", "--mainnet").CombinedOutput()
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CARDANO", "TYPE", "query tip", "ERROR", err, "OUTPUT", output)
	}

	err = json.Unmarshal(output, &tip)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
	}
	slog.Info("CARDAGO", "PACKAGE", "CARDANO", "KES", tip)

	return tip
}
