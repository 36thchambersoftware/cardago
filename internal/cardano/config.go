package cardano

import (
	"encoding/json"
	"fmt"
	"os"

	"cardano/cardago/internal/log"
)

type ShelleyGenesis struct {
	EpochLength int `json:"epochLength"`
}

type Config struct {
	NodeCertPath           string `yaml:"nodeCertPath"`
	ShelleyGenesisFilePath string `yaml:"shelleyGenesisFilePath"`
	StakePoolID            string `yaml:"stakepoolid"`
	VRFSKeyFilePath        string `yaml:"vrfskeyfilepath"`
}

type Logs struct {
	Leader Leader `yaml:"leader"`
}

type Leader struct {
	Directory string `yaml:"directory"`
	Prefix    string `yaml:"prefix"`
	Extension string `yaml:"extension"`
}

func (cfg *Config) GetShelleyGenesis() *ShelleyGenesis {
	logger := log.InitializeLogger()
	ShelleyGenesis := &ShelleyGenesis{}
	data, err := os.ReadFile(cfg.ShelleyGenesisFilePath)
	if err != nil {
		logger.Errorw("CARDAGO", "PACKAGE", "CARDANO", "GetEpochLength", err)
	}
	err = json.Unmarshal(data, ShelleyGenesis)
	if err != nil {
		logger.Errorw("CARDAGO", "PACKAGE", "CARDANO", "GetEpochLength", err)
	}

	return ShelleyGenesis
}

func (logs *Logs) GetLeaderPath(epoch int) string {
	logger := log.InitializeLogger()
	logger.Infow("GetLeaderPath", "directory", logs.Leader.Directory, "prefix", logs.Leader.Prefix, "epoch", epoch, "extension", logs.Leader.Extension)
	return fmt.Sprintf("%s/%s%d%s", logs.Leader.Directory, logs.Leader.Prefix, epoch, logs.Leader.Extension)
}
