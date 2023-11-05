package cardano

import (
	"log/slog"
	"os/exec"
)

func Run(args []string) ([]byte, error) {
	output, err := exec.Command("cardano-cli", args...).CombinedOutput()
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err, "OUTPUT", output)
	}
	slog.Info("CARDAGO", "PACKAGE", "CARDANO", "OUTPUT", output)

	return output, err
}
