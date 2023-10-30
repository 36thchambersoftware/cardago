package config

import (
	"os"

	"cardano/cardago/internal/discord"

	"github.com/sagikazarmark/slog-shim"
	"github.com/spf13/viper"
)

type Config struct {
	NodeCertPath       string                `yaml:"nodeCertPath"`
	LeaderLogDirectory string                `yaml:"leaderLogDirectory"`
	LeaderLogPrefix    string                `yaml:"leaderLogPrefix"`
	LeaderLogExtension string                `yaml:"leaderLogExtension"`
	Discord            discord.DiscordConfig `yaml:"discord"`
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

	_, err = os.Stat(cfg.NodeCertPath)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CONFIG", "ERROR", err)
		return err
	}

	_, err = os.Stat(cfg.LeaderLogDirectory)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CONFIG", "ERROR", err)
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
