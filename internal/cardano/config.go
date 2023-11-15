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
	Leader                 Leader `yaml:"leader"`
}

type Leader struct {
	Directory string `yaml:"directory"`
	Prefix    string `yaml:"prefix"`
	Extension string `yaml:"extension"`
}

func (cfg *Config) GetShelleyGenesis() *ShelleyGenesis {
	ShelleyGenesis := &ShelleyGenesis{}
	data, err := os.ReadFile(cfg.ShelleyGenesisFilePath)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CARDANO", "GetEpochLength", err)
	}
	err = json.Unmarshal(data, ShelleyGenesis)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CARDANO", "GetEpochLength", err)
	}

	return ShelleyGenesis
}

func (leader Leader) GetLeaderPath(epoch int) string {
	log.Infow("GetLeaderPath", "directory", leader.Directory, "prefix", leader.Prefix, "epoch", epoch, "extension", leader.Extension)
	return fmt.Sprintf("%s/%s%d%s", leader.Directory, leader.Prefix, epoch, leader.Extension)
}
