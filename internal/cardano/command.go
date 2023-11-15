package cardano

import (
	"os/exec"

	"cardano/cardago/internal/log"
)

func Run(args []string) ([]byte, error) {
	output, err := exec.Command("/usr/local/bin/cardano-cli", args...).CombinedOutput()
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err, "OUTPUT", string(output))
	}
	log.Debugw("CARDAGO", "PACKAGE", "CARDANO", "OUTPUT", string(output))

	return output, err
}
