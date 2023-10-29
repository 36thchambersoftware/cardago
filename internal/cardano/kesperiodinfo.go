package cardano

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
)

type KESPeriodInfo struct {
	KesKesKeyExpiry                          string `json:"qKesKesKeyExpiry"`
	KesCurrentKesPeriod                      int64  `json:"qKesCurrentKesPeriod"`
	KesEndKesInterval                        int64  `json:"qKesEndKesInterval"`
	KesMaxKESEvolutions                      int64  `json:"qKesMaxKESEvolutions"`
	KesNodeStateOperationalCertificateNumber int64  `json:"qKesNodeStateOperationalCertificateNumber"`
	KesOnDiskOperationalCertificateNumber    int64  `json:"qKesOnDiskOperationalCertificateNumber"`
	KesRemainingSlotsInKesPeriod             int64  `json:"qKesRemainingSlotsInKesPeriod"`
	KesSlotsPerKesPeriod                     int64  `json:"qKesSlotsPerKesPeriod"`
	KesStartKesInterval                      int64  `json:"qKesStartKesInterval"`
}

func GetKESPeriodInfo(network string, nodeCertPath string) KESPeriodInfo {
	kes := KESPeriodInfo{}
	// Create a new command
	slog.Info("CARDAGO", "PACKAGE", "CARDANO", "NETWORK", network, "PATH", nodeCertPath)
	command := fmt.Sprintf("cardano-cli query kes-period-info --%s --op-cert-file %s | tail -n +3", network, nodeCertPath)
	slog.Info("CARDAGO", "PACKAGE", "CARDANO", "COMMAND", command)
	output, err := exec.Command("bash", "-c", command).CombinedOutput()
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err, "OUTPUT", output)
	}
	slog.Info("CARDAGO", "PACKAGE", "CARDANO", "OUTPUT", output)

	// Decode the json data
	err = json.Unmarshal(output, &kes)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
	}
	slog.Info("CARDAGO", "PACKAGE", "CARDANO", "KES", kes)

	return kes

	// For testing, these numbers might be like this..
	// kestemp := KESPeriodInfo{
	// 	KesCurrentKesPeriod:                      816,
	// 	KesEndKesInterval:                        878,
	// 	KesKesKeyExpiry:                          "2024-01-15T21:44:51Z",
	// 	KesMaxKESEvolutions:                      62,
	// 	KesNodeStateOperationalCertificateNumber: 4,
	// 	KesOnDiskOperationalCertificateNumber:    5,
	// 	KesRemainingSlotsInKesPeriod:             7977800,
	// 	KesSlotsPerKesPeriod:                     129600,
	// 	KesStartKesInterval:                      816,
	// }

	// return kestemp
}
