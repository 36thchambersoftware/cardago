package config

import (
	"cardano/cardago/internal/cardano"
	"cardano/cardago/internal/discord"
	"cardano/cardago/internal/log"
	"cardano/cardago/pkg/blockfrost"

	"github.com/spf13/viper"
)

type Config struct {
	Cardano    cardano.Config    `yaml:"cardano"`
	Leader     cardano.Leader    `yaml:"leader"`
	Discord    discord.Config    `yaml:"discord"`
	Blockfrost blockfrost.Config `yaml:"blockfrost"`
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
		log.Errorw("CARDAGO", "PACKAGE", "CONFIG", "ERROR", err)
		return err
	}

	// Unmarshal the configuration into the Config struct.
	err = viper.Unmarshal(cfg)
	if err != nil {
		log.Errorw("CARDAGO", "PACKAGE", "CONFIG", "ERROR", err)
		return err
	}
	log.Debugw("CARDAGO", "PACKAGE", "CONFIG", "viper config", cfg)

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
		log.Errorw("CARDAGO", "PACKAGE", "CONFIG", "ERROR", err)
	}

	return C
}
