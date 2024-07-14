package cardano

import (
	"encoding/json"
	"strings"

	"cardano/cardago/internal/log"
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
	log.Infow("CARDAGO", "PACKAGE", "CARDANO", "ACTION", "GetKESPeriodInfo")
	kes := KESPeriodInfo{}
	log.Debugw("CARDAGO", "PACKAGE", "CARDANO", "PATH", opCertPath)

	args := []string{
		"conway", "query",
		"kes-period-info",
		"--mainnet",
		"--op-cert-file",
		opCertPath,
	}
	log.Debugw("CARDAGO", "PACKAGE", "CARDANO", "ARGS", args)

	output, err := Run(args)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
	}

	// convert raw output to KES model
	individualLines := strings.Split(string(output), "\n")
	dataArray := individualLines[2:]
	data := strings.Join(dataArray, "")

	// Decode the json data
	err = json.Unmarshal([]byte(data), &kes)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CARDANO", "ERROR", err)
		return kes, err
	}
	log.Infow("CARDAGO", "PACKAGE", "CARDANO", "KES", kes)

	return kes, err
}
