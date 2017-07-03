package commands

import (
	"fmt"
	"os"

	"github.com/ashmckenzie/go-mqti/mqti"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string
var debug, showVersion bool

// RootCmd ...
var RootCmd = &cobra.Command{
	Use:           "mqti",
	Short:         "MQTT subscriber that pumps data into InfluxDB",
	Long:          `MQTT subscriber that pumps data into InfluxDB`,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := validateInput(); err != nil {
			return err
		}

		mqti.EnableDebugging(debug)

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			fmt.Println(mqti.Version)
		}
		cmd.Help()
	},
}

func validateInput() error {
	return nil
}

// Execute ...
func Execute() error {
	err := RootCmd.Execute()
	return err
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "show version")
	RootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debugging")

	RootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is config.yaml)")
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		mqti.Log.Fatal("Can't read config:", err)
		os.Exit(1)
	}
}
