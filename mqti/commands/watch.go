package commands

import (
	"github.com/ashmckenzie/go-mqti/mqti"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch MQTT messages",
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()
		watchMessages()
	},
}

func init() {
	RootCmd.AddCommand(watchCmd)
}

func watchMessages() {
	incoming := make(chan *mqti.MQTTMessage)
	go mqti.MQTTSubscribe(incoming)

	for m := range incoming {
		mqti.LogMQTTMessage(m)
	}
}
