package config

import (
	"os"

	"cardano/cardago/internal/cardano"
	"cardano/cardago/internal/discord"

	"github.com/sagikazarmark/slog-shim"
	"github.com/spf13/viper"
)

type Config struct {
	Cardano cardano.Config `yaml:"cardano"`
	Logs    cardano.Logs   `yaml:"logs"`
	Discord discord.Config `yaml:"discord"`
}

func (cfg *Config) LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../.")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CONFIG", "ERROR", err)
		return err
	}

	err = viper.Unmarshal(cfg)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CONFIG", "ERROR", err)
		return err
	}
	slog.Info("CARDAGO", "PACKAGE", "CONFIG", "viper config", cfg)

	_, err = os.Stat(cfg.Logs.Leader.Directory)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CONFIG", "Log.Leader.Directory", err)
		slog.Error("CARDAGO", "PACKAGE", "CONFIG", "Log.Leader.Directory", cfg.Logs)

		return err
	}

	_, err = os.Stat(cfg.Cardano.NodeCertPath)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CONFIG", "Cardano.NodeCertPath", err)
		return err
	}

	_, err = os.Stat(cfg.Cardano.ShelleyGenesisFilePath)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CONFIG", "Cardano.ShelleyGenesisFilePath", err)
		return err
	}

	_, err = os.Stat(cfg.Cardano.VRFSKeyFilePath)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CONFIG", "Cardano.VRFSKeyFilePath", err)
		return err
	}

	return err
}

func Get() *Config {
	C := &Config{}

	err := viper.Unmarshal(C)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CONFIG", "ERROR", err)
	}

	return C
}
