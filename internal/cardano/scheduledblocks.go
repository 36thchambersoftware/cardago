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

var (
	ErrorEpochAlreadyChecked = errors.New("epoch already checked")
	ErrorOutOfSync           = errors.New("sync progress is less than 100 - wait for full sync")
	ErrorEpochLength         = errors.New("epoch length incorrect - check shelley genesis file")
	ErrorTooEarly            = errors.New("too early to check")
)

func GetScheduledBlocks(cfg Config) ([]ScheduledBlock, error) {
	log.Infow("CARDAGO", "PACKAGE", "CARDANO", "ACTION", "GetScheduledBlocks")
	scheduledBlocks := []ScheduledBlock{}
	tip, err := QueryTip()
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
		return []ScheduledBlock{}, err
	}

	if tip.SyncProgress != "100.00" {
		return []ScheduledBlock{}, ErrorOutOfSync
	}

	epochLength := cfg.GetShelleyGenesis().EpochLength
	log.Infow("CARDAGO", "PACKAGE", "CARDANO", "EPOCHLENGTH", epochLength)
	if epochLength < 1 {
		return []ScheduledBlock{}, ErrorEpochLength
	}

	epochProgress := float32(tip.SlotInEpoch) / float32(epochLength)
	log.Infow("CARDAGO", "PACKAGE", "CARDANO", "Epoch progress", epochProgress)
	if epochProgress < 0.75 {
		return []ScheduledBlock{}, ErrorTooEarly
	}

	nextEpoch := tip.Epoch + 1
	path := cfg.Leader.GetLeaderPath(nextEpoch)
	log.Infow("CARDAGO", "PACKAGE", "CARDANO", "PATH", path)
	existingLogFile, err := os.Stat(path)
	doesNotExist := os.IsNotExist(err)
	log.Infow("CARDAGO", "PACKAGE", "CARDANO", "doesNotExist", doesNotExist)
	if !doesNotExist {
		log.Infow("CARDAGO", "PACKAGE", "CARDANO", "name", existingLogFile.Name(), "modified", existingLogFile.ModTime())
		log.Infow("CARDAGO", "PACKAGE", "CARDANO", "epoch already checked", nextEpoch)
		return scheduledBlocks, ErrorEpochAlreadyChecked
	}


	args := []string{
		"conway", "query", "leadership-schedule",
		"--mainnet",
		"--genesis", cfg.ShelleyGenesisFilePath,
		"--stake-pool-id", cfg.StakePoolID,
		"--vrf-signing-key-file", cfg.VRFSKeyFilePath,
		"--next",
	}
	log.Debugw("CARDAGO", "PACKAGE", "CARDANO", "ARGS", args)

	output, err := Run(args)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err, "OUTPUT", string(output))
		return scheduledBlocks, err
	}

	err = logScheduledBlocks(cfg.Leader, nextEpoch, output)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CARDANO", "logScheduledBlocks", err)
		return scheduledBlocks, err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) == 0 || len(output) == 0 {
		log.Infow("CARDAGO", "PACKAGE", "CARDANO", "BLOCKS", "No scheduled blocks")
		return scheduledBlocks, err
	}

	// the first two lines of the output are ignored
	//      SlotNo                          UTC Time
	// -------------------------------------------------------------
	blockLines := lines[2:]

	if len(blockLines) < 1 {
		log.Infow("CARDAGO", "PACKAGE", "CARDANO", "lines", lines, "blockLines", blockLines)
		return scheduledBlocks, err
	}

	for _, line := range blockLines {
		scheduledBlock := ScheduledBlock{}
		pieces := strings.Fields(line)
		log.Infow("CARDAGO", "PACKAGE", "CARDANO", "pieces", pieces)
		if len(pieces) == 0 {
			log.Infow("CARDAGO", "PACKAGE", "CARDANO", "scheduled blocks", "ZERO")
			return scheduledBlocks, err
		}

		if len(pieces) != 4 {
			log.Infow("CARDAGO", "PACKAGE", "CARDANO", "pieces mismatch", "check output")
			return scheduledBlocks, err
		}

		slotNumber, errIntParse := strconv.ParseInt(pieces[0], 10, 64) // e.g. 97470387
		if errIntParse != nil {
			log.Errorw("CARDAGO", "PACKAGE", "CARDANO", "ERROR", errIntParse)
			return scheduledBlocks, err
		}

		slotTime, errTimeParse := time.Parse("2006-01-02 15:04:05 MST", strings.Join(pieces[1:], " "))
		if errTimeParse != nil {
			log.Errorw("CARDAGO", "PACKAGE", "CARDANO", "ERROR", errTimeParse)
			return scheduledBlocks, err
		}

		local, _ := time.LoadLocation("America/Chicago")

		scheduledBlock.SlotNumber = slotNumber
		scheduledBlock.Datetime = slotTime.In(local)

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
func logScheduledBlocks(logs Leader, nextEpoch int, content []byte) error {
	// Get the leader log file path.
	path := logs.GetLeaderPath(nextEpoch)
	log.Infow("CARDAGO", "PACKAGE", "CARDANO", "logScheduledBlocks", path)

	// Create the leader log file.
	f, err := os.Create(path)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CARDANO", "cannot create leader log file", err)
		return err
	}

	defer f.Close()

	// Write the scheduled Cardano blocks to the leader log file.
	_, err = f.WriteString(string(content))
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CARDANO", "cannot write to leader log file", err)
		return err
	}

	// Return nil to indicate that the scheduled Cardano blocks were logged successfully.
	return err
}

func ReadScheduledBlocks(logs Leader) ([]byte, error) {
	tip, err := QueryTip()
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
		return []byte{}, err
	}
	return os.ReadFile(logs.GetLeaderPath(tip.Epoch + 1))
}
