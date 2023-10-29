package cardano

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type ScheduledBlock struct {
	SlotNumber int64
	Datetime   time.Time
}

func GetScheduledBlocks(leaderLogDirectory string, leaderLogPrefix string, leaderLogExtension string) []ScheduledBlock {
	slog.Info("CARDAGO", "PACKAGE", "CARDANO", "ACTION", "GetScheduledBlocks")
	scheduledBlocks := []ScheduledBlock{}
	epoch := QueryTip().Epoch
	sEpoch := strconv.Itoa(epoch)
	path := fmt.Sprintf("%s/%s%s.%s", leaderLogDirectory, leaderLogPrefix, sEpoch, leaderLogExtension)
	command := fmt.Sprintf("tail -n +3 %s", path)
	slog.Info("CARDAGO", "PACKAGE", "CARDANO", "COMMAND", command)

	output, err := exec.Command("bash", "-c", command).CombinedOutput()
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err, "OUTPUT", output)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) == 0 || len(output) == 0 {
		slog.Info("CARDAGO", "PACKAGE", "CARDANO", "BLOCKS", "No scheduled blocks")
		return scheduledBlocks
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

	return scheduledBlocks
}
