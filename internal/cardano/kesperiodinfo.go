package cardano

import (
	"encoding/json"
	"log/slog"
	"strings"
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

func GetKESPeriodInfo(opCertPath string) (KESPeriodInfo, error) {
	kes := KESPeriodInfo{}
	slog.Info("CARDAGO", "PACKAGE", "CARDANO", "PATH", opCertPath)

	args := []string{
		"query",
		"kes-period-info",
		"--mainnet",
		"--op-cert-file",
		opCertPath,
	}
	slog.Info("CARDAGO", "PACKAGE", "CARDANO", "ARGS", args)

	output, err := Run(args)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
	}

	// convert raw output to KES model
	individualLines := strings.Split(string(output), "\n")
	dataArray := individualLines[2:]
	data := strings.Join(dataArray, "")

	// Decode the json data
	err = json.Unmarshal([]byte(data), &kes)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
		return kes, err
	}
	slog.Info("CARDAGO", "PACKAGE", "CARDANO", "KES", kes)

	return kes, err
}
