package commands

import (
	"github.com/ashmckenzie/go-littlefly/littlefly"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch MQTT messages",
	Run: func(cmd *cobra.Command, args []string) {
		watchMessages()
	},
}

func init() {
	RootCmd.AddCommand(watchCmd)
}

func watchMessages() {
	incoming := make(chan *littlefly.MQTTMessage)
	go littlefly.MQTTSubscribe(incoming)

	for m := range incoming {
		littlefly.LogMQTTMessage(m)
	}
}
