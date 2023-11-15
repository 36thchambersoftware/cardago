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
	Leader  cardano.Leader `yaml:"leader"`
	Discord discord.Config `yaml:"discord"`
}

/**
 * Loads the configuration from the configuration file.
 *
 * @return error Any error that occurred while loading the configuration.
 */
func (cfg *Config) LoadConfig() error {
	// Set the configuration name and type.
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Add the configuration search paths.
	viper.AddConfigPath("../../.")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/home/cardano/scripts/cardago/.")

	// Read the configuration file.
	err := viper.ReadInConfig()
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CONFIG", "ERROR", err)
		return err
	}

	// Unmarshal the configuration into the Config struct.
	err = viper.Unmarshal(cfg)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CONFIG", "ERROR", err)
		return err
	}
	slog.Info("CARDAGO", "PACKAGE", "CONFIG", "viper config", cfg)

	// Check if the leader log directory exists.
	_, err = os.Stat(cfg.Leader.Directory)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CONFIG", "Log.Leader.Directory", err)
		slog.Error("CARDAGO", "PACKAGE", "CONFIG", "Log.Leader.Directory", cfg.Leader)

		return err
	}

	// Check if the Cardano node certificate path exists.
	_, err = os.Stat(cfg.Cardano.NodeCertPath)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CONFIG", "Cardano.NodeCertPath", err)
		return err
	}

	// Check if the Cardano Shelley genesis file path exists.
	_, err = os.Stat(cfg.Cardano.ShelleyGenesisFilePath)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CONFIG", "Cardano.ShelleyGenesisFilePath", err)
		return err
	}

	// Check if the Cardano VRFS key file path exists.
	_, err = os.Stat(cfg.Cardano.VRFSKeyFilePath)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CONFIG", "Cardano.VRFSKeyFilePath", err)
		return err
	}

	// Return nil to indicate that the configuration was loaded successfully.
	return err
}

/**
 * Get the configuration from the configuration file.
 *
 * @return *Config The configuration.
 * @error error Any error that occurred while getting the configuration.
 */
func Get() *Config {
	C := &Config{}

	err := viper.Unmarshal(C)
	if err != nil {
		slog.Error("CARDAGO", "PACKAGE", "CONFIG", "ERROR", err)
	}

	return C
}
