package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var (
	ConfigFile string
	Debug      bool
)

// initConfig uses Viper's order of precedence to build the CLI configuration
func initConfig() {
	if ConfigFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(ConfigFile)
	} else {
		// Search for config using order of precedence:
		// - /etc/goldap/goldap.<EXT>
		// - $HOME/.goldap.<EXT>
		// - ./.goldap.<EXT>
		viper.SetConfigName(".goldap")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("/etc/goldap/")
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath(".")
	}

	// Set the env var prefix and read in matching vars
	viper.SetEnvPrefix("goldap")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Sprintf("Error parsing config file %s: %s", viper.ConfigFileUsed(), err)
	}

	// Bind all command flags to Viper config
	if err := viper.BindPFlags(rootCmd.Flags()); err != nil {
		fmt.Sprintf("Error binding config flags with viper: %s", err)
		os.Exit(1)
	}
}
