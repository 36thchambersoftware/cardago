package cardano

import (
	"errors"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

type ScheduledBlock struct {
	Datetime   time.Time
	SlotNumber int64
}

func GetScheduledBlocks(config Config, logs Logs) ([]ScheduledBlock, error) {
	slog.Info("CARDAGO", "PACKAGE", "CARDANO", "ACTION", "GetScheduledBlocks")
	scheduledBlocks := []ScheduledBlock{}
	tip, err := QueryTip()
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
		return []ScheduledBlock{}, err
	}

	if tip.SyncProgress != "100.00" {
		return []ScheduledBlock{}, errors.New("sync progress is less than 100 - wait for full sync")
	}

	epochLength := config.GetShelleyGenesis().EpochLength
	slog.Info("CARDAGO", "PACKAGE", "CARDANO", "EPOCHLENGTH", epochLength)
	if epochLength < 1 {
		return []ScheduledBlock{}, errors.New("epoch length incorrect - check shelley genesis file")
	}

	epochProgress := float32(tip.SlotInEpoch) / float32(epochLength)
	slog.Info("CARDAGO", "PACKAGE", "CARDANO", "Epoch progress", epochProgress)
	if epochProgress < 0.75 {
		return []ScheduledBlock{}, errors.New("too early to check")
	}

	nextEpoch := tip.Epoch + 1
	path := logs.GetLeaderPath(nextEpoch)
	slog.Info("CARDAGO", "PACKAGE", "CARDANO", "PATH", path)
	existingLogFile, err := os.Stat(path)
	doesNotExist := os.IsNotExist(err)
	slog.Info("CARDAGO", "PACKAGE", "CARDANO", "doesNotExist", doesNotExist)
	if !doesNotExist {
		slog.Info("CARDAGO", "PACKAGE", "CARDANO", "name", existingLogFile.Name(), "modified", existingLogFile.ModTime())
		slog.Info("CARDAGO", "PACKAGE", "CARDANO", "epoch already checked", nextEpoch)
		return scheduledBlocks, errors.New("epoch already checked")
	}

	args := []string{
		"query", "leadership-schedule",
		"--mainnet",
		"--genesis", config.ShelleyGenesisFilePath,
		"--stake-pool-id", config.StakePoolID,
		"--vrf-signing-key-file", config.VRFSKeyFilePath,
		"--next",
	}
	slog.Info("CARDAGO", "PACKAGE", "CARDANO", "ARGS", args)

	output, err := Run(args)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err, "OUTPUT", output)
		return scheduledBlocks, err
	}

	err = logScheduledBlocks(logs, nextEpoch, output)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CARDANO", "logScheduledBlocks", err)
		return scheduledBlocks, err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) == 0 || len(output) == 0 {
		slog.Info("CARDAGO", "PACKAGE", "CARDANO", "BLOCKS", "No scheduled blocks")
		return scheduledBlocks, err
	}

	// the first two lines of the output are ignored
	//      SlotNo                          UTC Time
	// -------------------------------------------------------------
	blockLines := lines[2:]

	if len(blockLines) < 1 {
		slog.Info("CARDAGO", "PACKAGE", "CARDANO", "lines", lines, "blockLines", blockLines)
		return scheduledBlocks, err
	}

	for _, line := range blockLines {
		scheduledBlock := ScheduledBlock{}
		pieces := strings.Fields(line)
		slog.Info("CARDAGO", "PACKAGE", "CARDANO", "pieces", pieces)
		if len(pieces) == 0 {
			slog.Info("CARDAGO", "PACKAGE", "CARDANO", "no scheduled blocks")
			return scheduledBlocks, err
		}

		if len(pieces) != 2 {
			slog.Info("CARDAGO", "PACKAGE", "CARDANO", "pieces mismatch - check output")
			return scheduledBlocks, err
		}

		slotNumber, err := strconv.ParseInt(pieces[0], 10, 64) // e.g. 97470387
		if err != nil {
			slog.Error("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
		}

		slotTime, err := time.Parse("2006-01-02T15:04:05Z", pieces[1])
		if err != nil {
			slog.Error("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
		}

		scheduledBlock.SlotNumber = slotNumber
		scheduledBlock.Datetime = slotTime

		scheduledBlocks = append(scheduledBlocks, scheduledBlock)
	}

	return scheduledBlocks, err
}

func logScheduledBlocks(logs Logs, nextEpoch int, content []byte) error {
	path := logs.GetLeaderPath(nextEpoch)
	slog.Info("CARDAGO", "PACKAGE", "CARDANO", "logScheduledBlocks", path)

	f, err := os.Create(path)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CARDANO", "cannot create leader log file", err)
		return err
	}

	defer f.Close()
	_, err = f.WriteString(string(content))
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CARDANO", "cannot write to leader log file", err)
		return err
	}

	return err
}
