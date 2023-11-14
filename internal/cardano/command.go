package cardano

import (
	"os/exec"

	"cardano/cardago/internal/log"
)

func Run(args []string) ([]byte, error) {
	logger := log.InitializeLogger()

	output, err := exec.Command("/usr/local/bin/cardano-cli", args...).CombinedOutput()
	if err != nil {
		logger.Errorw("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err, "OUTPUT", output)
	}
	logger.Infow("CARDAGO", "PACKAGE", "CARDANO", "OUTPUT", output)

	return output, err
}
