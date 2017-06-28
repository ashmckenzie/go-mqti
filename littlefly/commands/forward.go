package commands

import (
	"github.com/ashmckenzie/go-littlefly/littlefly"
	"github.com/spf13/cobra"
)

var forwardCmd = &cobra.Command{
	Use:   "forward",
	Short: "Forward MQTT messages on to InfluxDB",
	Run: func(cmd *cobra.Command, args []string) {
		forwardMessages()
	},
}

func init() {
	RootCmd.AddCommand(forwardCmd)
}

func forwardMessages() {
	influxDB, _ := littlefly.NewInfluxDBConnection()

	incoming := make(chan *littlefly.MQTTMessage)
	forward := make(chan *littlefly.MQTTMessage)

	go littlefly.CreateWorkers(influxDB, forward)
	go littlefly.MQTTSubscribe(incoming)

	for m := range incoming {
		littlefly.LogMQTTMessage(m)
		forward <- m
	}
}
