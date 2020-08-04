package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goldap",
	Short: "CLI to interact with LDAP servers",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Execute root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Root command flags
	rootCmd.PersistentFlags().StringVarP(&ConfigFile, "config-file", "c", "", "Path to the config file")
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Enable debug logs")

	// Bind all command flags to Viper config
	if err := viper.BindPFlags(rootCmd.Flags()); err != nil {
		Logger.DPanicf("Error binding config flags with viper: %s", err)
		os.Exit(1)
	}

	// Reconfigure Logger based on --debug flag and config
	// If debug flag is passed, enable development mode and set log level

	// loggingFromConfig := viper.Get

	// logLevel :=

	// if Debug {
	//	viper.Set("logging.development", true)
	// 	viper.Set("logging.level", "debug")
	// } else {
	// 	if !viper.IsSet("logging.development") {
	// 		viper.Set("logging.development", false)
	// 	}
	// 	if !viper.IsSet("logging.level") {
	// 		viper.Set("logging.level", "info")
	// 	}
	// }

	// atomicLevel :=

	// // Logging config
	// logger, err := logging.NewLoggerFromConfig(&logging.ConfigInput{
	// 	Development:       viper.GetBool("logging.development"),
	// 	DisableCaller:     viper.GetBool("logging.disableCaller"),
	// 	DisableStacktrace: viper.GetBool("logging.disableStacktrace"),
	// 	ErrorOutputPaths:  viper.GetStringSlice("logging.errorOutputPaths"),
	// 	InitialFields:     viper.GetStringMap("logging.initialFields"),
	// 	Level:             viper.GetString("logging.level"),
	// 	OutputPaths:       viper.GetStringSlice("logging.outputPaths"),
	// })

	// if err != nil {
	// 	fmt.Sprintf("Failed to initialize logger: %s", err)
	// }

	// Logger = logger
}
