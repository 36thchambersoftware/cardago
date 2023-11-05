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

	path := logs.GetLeaderPath(tip.Epoch)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// The file does not exist
		slog.Info("CARDAGO", "PACKAGE", "CARDANO", "Epoch", tip.Epoch)
	} else {
		// The file exists
		slog.Info("CARDAGO", "PACKAGE", "CARDANO", "epoch already checked", tip.Epoch)
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

	err = logScheduledBlocks(logs, tip.Epoch, output)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CARDANO", "logScheduledBlocks", err)
		return scheduledBlocks, err
	}

	err = logScheduledBlocks(logs, tip.Epoch, output)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CARDANO", "logScheduledBlocks", err)
		return scheduledBlocks, err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) == 0 || len(output) == 0 {
		slog.Info("CARDAGO", "PACKAGE", "CARDANO", "BLOCKS", "No scheduled blocks")
		return scheduledBlocks, err
	}

	for _, line := range lines {
		scheduledBlock := ScheduledBlock{}
		pieces := strings.Fields(line)
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

func logScheduledBlocks(logs Logs, epoch int, content []byte) error {
	nextEpoch := epoch + 1
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
