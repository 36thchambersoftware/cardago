package cardano

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"cardano/cardago/internal/log"
)

type ScheduledBlock struct {
	Datetime   time.Time
	SlotNumber int64
}

func GetScheduledBlocks(cfg Config, logs Logs) ([]ScheduledBlock, error) {
	logger := log.InitializeLogger()
	logger.Infow("CARDAGO", "PACKAGE", "CARDANO", "ACTION", "GetScheduledBlocks")
	scheduledBlocks := []ScheduledBlock{}
	tip, err := QueryTip()
	if err != nil {
		logger.Errorw("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
		return []ScheduledBlock{}, err
	}

	if tip.SyncProgress != "100.00" {
		return []ScheduledBlock{}, errors.New("sync progress is less than 100 - wait for full sync")
	}

	epochLength := cfg.GetShelleyGenesis().EpochLength
	logger.Infow("CARDAGO", "PACKAGE", "CARDANO", "EPOCHLENGTH", epochLength)
	if epochLength < 1 {
		return []ScheduledBlock{}, errors.New("epoch length incorrect - check shelley genesis file")
	}

	epochProgress := float32(tip.SlotInEpoch) / float32(epochLength)
	logger.Infow("CARDAGO", "PACKAGE", "CARDANO", "Epoch progress", epochProgress)
	if epochProgress < 0.75 {
		return []ScheduledBlock{}, errors.New("too early to check")
	}

	nextEpoch := tip.Epoch + 1
	path := logs.GetLeaderPath(nextEpoch)
	logger.Infow("CARDAGO", "PACKAGE", "CARDANO", "PATH", path)
	existingLogFile, err := os.Stat(path)
	doesNotExist := os.IsNotExist(err)
	logger.Infow("CARDAGO", "PACKAGE", "CARDANO", "doesNotExist", doesNotExist)
	if !doesNotExist {
		logger.Infow("CARDAGO", "PACKAGE", "CARDANO", "name", existingLogFile.Name(), "modified", existingLogFile.ModTime())
		logger.Infow("CARDAGO", "PACKAGE", "CARDANO", "epoch already checked", nextEpoch)
		return scheduledBlocks, errors.New("epoch already checked")
	}

	args := []string{
		"query", "leadership-schedule",
		"--mainnet",
		"--genesis", cfg.ShelleyGenesisFilePath,
		"--stake-pool-id", cfg.StakePoolID,
		"--vrf-signing-key-file", cfg.VRFSKeyFilePath,
		"--next",
	}
	logger.Infow("CARDAGO", "PACKAGE", "CARDANO", "ARGS", args)

	output, err := Run(args)
	if err != nil {
		logger.Errorw("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err, "OUTPUT", output)
		return scheduledBlocks, err
	}

	err = logScheduledBlocks(logs, nextEpoch, output)
	if err != nil {
		logger.Errorw("CARDAGO", "PACKAGE", "CARDANO", "logScheduledBlocks", err)
		return scheduledBlocks, err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) == 0 || len(output) == 0 {
		logger.Infow("CARDAGO", "PACKAGE", "CARDANO", "BLOCKS", "No scheduled blocks")
		return scheduledBlocks, err
	}

	// the first two lines of the output are ignored
	//      SlotNo                          UTC Time
	// -------------------------------------------------------------
	blockLines := lines[2:]

	if len(blockLines) < 1 {
		logger.Infow("CARDAGO", "PACKAGE", "CARDANO", "lines", lines, "blockLines", blockLines)
		return scheduledBlocks, err
	}

	for _, line := range blockLines {
		scheduledBlock := ScheduledBlock{}
		pieces := strings.Fields(line)
		logger.Infow("CARDAGO", "PACKAGE", "CARDANO", "pieces", pieces)
		if len(pieces) == 0 {
			logger.Infow("CARDAGO", "PACKAGE", "CARDANO", "scheduled blocks", "ZERO")
			return scheduledBlocks, err
		}

		if len(pieces) != 2 {
			logger.Infow("CARDAGO", "PACKAGE", "CARDANO", "pieces mismatch", "check output")
			return scheduledBlocks, err
		}

		slotNumber, errIntParse := strconv.ParseInt(pieces[0], 10, 64) // e.g. 97470387
		if errIntParse != nil {
			logger.Errorw("CARDAGO", "PACKAGE", "CARDANO", "ERROR", errIntParse)
		}

		slotTime, errTimeParse := time.Parse("2006-01-02T15:04:05Z", pieces[1])
		if errTimeParse != nil {
			logger.Errorw("CARDAGO", "PACKAGE", "CARDANO", "ERROR", errTimeParse)
		}

		scheduledBlock.SlotNumber = slotNumber
		scheduledBlock.Datetime = slotTime

		scheduledBlocks = append(scheduledBlocks, scheduledBlock)
	}

	return scheduledBlocks, err
}

/**
 * Logs the scheduled Cardano blocks to the leader log file.
 *
 * @param logs The logs directory.
 * @param nextEpoch The next epoch.
 * @param content The scheduled Cardano blocks.
 * @return error Any error that occurred while logging the scheduled Cardano blocks.
 */
func logScheduledBlocks(logs Logs, nextEpoch int, content []byte) error {
	logger := log.InitializeLogger()
	// Get the leader log file path.
	path := logs.GetLeaderPath(nextEpoch)
	logger.Infow("CARDAGO", "PACKAGE", "CARDANO", "logScheduledBlocks", path)

	// Create the leader log file.
	f, err := os.Create(path)
	if err != nil {
		logger.Errorw("CARDAGO", "PACKAGE", "CARDANO", "cannot create leader log file", err)
		return err
	}

	defer f.Close()

	// Write the scheduled Cardano blocks to the leader log file.
	_, err = f.WriteString(string(content))
	if err != nil {
		logger.Errorw("CARDAGO", "PACKAGE", "CARDANO", "cannot write to leader log file", err)
		return err
	}

	// Return nil to indicate that the scheduled Cardano blocks were logged successfully.
	return err
}
