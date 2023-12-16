package command

import (
	"os/exec"

	"cardano/cardago/internal/log"
)

func Run(cmd string, args []string) ([]byte, error) {
	output, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "COMMAND", "ERROR", err, "OUTPUT", string(output))
	}
	log.Debugw("CARDAGO", "PACKAGE", "COMMAND", "OUTPUT", string(output))

	return output, err
}
