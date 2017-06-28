package commands

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	InfluxDBClient "github.com/influxdata/influxdb/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var watch bool
var cInflux chan MQTT.Message
var influxDBConn *InfluxDBClient.Client

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch MQTT messages",
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		influxDBConn, err = InfluxDBClient.NewClient(InfluxDBClient.Config{URL: *InfluxDBURI()})
		if err != nil {
			log.Fatal(err)
		}

		watchMessages()
	},
}

func init() {
	RootCmd.AddCommand(watchCmd)
}

// MQTTConfig ...
func MQTTConfig() map[string]interface{} {
	return viper.GetStringMap("mqtt")
}

// InfluxDBConfig ...
func InfluxDBConfig() map[string]interface{} {
	return viper.GetStringMap("influxdb")
}

// MQTTCBrokerURI ...
func MQTTCBrokerURI() string {
	return fmt.Sprintf("tcp://%s:%d", MQTTConfig()["host"], MQTTConfig()["port"])
}

// MQTTClientID ...
func MQTTClientID() string {
	return MQTTConfig()["client_id"].(string)
}

// MQTTUsername ...
func MQTTUsername() string {
	return MQTTConfig()["username"].(string)
}

// MQTTPassword ...
func MQTTPassword() string {
	return MQTTConfig()["password"].(string)
}

// MQTTTopic ...
func MQTTTopic() string {
	return MQTTConfig()["topic"].(string)
}

// InfluxDBURI ...
func InfluxDBURI() *url.URL {
	host, _ := url.Parse(fmt.Sprintf("http://%s:%d", InfluxDBConfig()["host"], InfluxDBConfig()["port"]))
	return host
}

// InfluxDBDatabase ...
func InfluxDBDatabase() string {
	return InfluxDBConfig()["database"].(string)
}

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	fmt.Printf("Received message on topic: %s\nMessage: %s\n", message.Topic(), message.Payload())

	pts := make([]InfluxDBClient.Point, 1)

	pts[0] = InfluxDBClient.Point{
		Measurement: "cpu_load_short",
		Tags:        map[string]string{},
		Fields:      map[string]interface{}{"value": string(message.Payload())},
		Time:        time.Now(),
	}

	bps := InfluxDBClient.BatchPoints{
		Points:   pts,
		Database: InfluxDBDatabase(),
	}

	_, err := influxDBConn.Write(bps)
	if err != nil {
		log.Fatal(err)
	}
}

func watchMessages() {
	var client MQTT.Client

	cs := make(chan os.Signal, 1)
	signal.Notify(cs, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-cs
		fmt.Println("signal received, exiting")
		os.Exit(0)
	}()

	opts := &MQTT.ClientOptions{
		ClientID:  MQTTClientID(),
		Username:  MQTTUsername(),
		Password:  MQTTPassword(),
		TLSConfig: tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert},
	}

	opts.AddBroker(MQTTCBrokerURI())

	opts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(MQTTTopic(), byte(0), onMessageReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	client = MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
