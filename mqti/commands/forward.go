package commands

import (
	"github.com/ashmckenzie/go-mqti/mqti"
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
	influxDB, _ := mqti.NewInfluxDBConnection()

	incoming := make(chan *mqti.MQTTMessage)
	forward := make(chan *mqti.MQTTMessage)

	go mqti.CreateWorkers(influxDB, forward)
	go mqti.MQTTSubscribe(incoming)

	for m := range incoming {
		mqti.DebugLogMQTTMessage(m)
		forward <- m
	}
}
