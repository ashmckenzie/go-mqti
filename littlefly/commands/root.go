package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/ashmckenzie/go-littlefly/littlefly"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var configFile string
var debug, showVersion bool

// RootCmd ...
var RootCmd = &cobra.Command{
	Use:           "littlefly",
	Short:         "MQTT subscriber that pumps data into InfluxDB",
	Long:          `MQTT subscriber that pumps data into InfluxDB`,
	SilenceErrors: true,
	// PreRunE: func(cmd *cobra.Command, args []string) error {
	//   return CheckRequiredFlags(cmd.Flags())
	// },
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := validateInput(); err != nil {
			return err
		}

		littlefly.EnableDebugging(debug)

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			fmt.Println(littlefly.Version)
		}
	},
}

// CheckRequiredFlags ...
func CheckRequiredFlags(flags *pflag.FlagSet) error {
	requiredError := false
	flagName := ""

	flags.VisitAll(func(flag *pflag.Flag) {
		requiredAnnotation := flag.Annotations[cobra.BashCompOneRequiredFlag]
		if len(requiredAnnotation) == 0 {
			return
		}

		flagRequired := requiredAnnotation[0] == "true"

		if flagRequired && !flag.Changed {
			requiredError = true
			flagName = flag.Name
		}
	})

	if requiredError {
		return errors.New("Required flag `" + flagName + "` has not been set")
	}

	return nil
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

	RootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is config.toml)")
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
		littlefly.Log.Fatal("Can't read config:", err)
		os.Exit(1)
	}
}
