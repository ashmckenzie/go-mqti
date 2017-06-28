package littlefly

import (
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
)

// MQTTMessage ...
type MQTTMessage struct {
	MQTT.Message
	Mapping
}

// MQTTConfig ...
func MQTTConfig() map[string]interface{} {
	return viper.GetStringMap("mqtt")
}

// MQTTCBrokerURI ...
func MQTTCBrokerURI() string {
	return fmt.Sprintf("tcp://%s:%s", MQTTConfig()["host"], MQTTConfig()["port"])
}

// MQTTClientID ...
func MQTTClientID() string {
	return MQTTConfig()["client_id"].(string)
}

// MQTTUsername ...
func MQTTUsername() string {
	var u interface{}
	if u = MQTTConfig()["username"]; err != nil {
		return u.(string)
	}
	return ""
}

// MQTTPassword ...
func MQTTPassword() string {
	var p interface{}
	if p = MQTTConfig()["password"]; err != nil {
		return p.(string)
	}
	return ""
}

// MQTTTopic ...
func MQTTTopic() string {
	return MQTTConfig()["topic"].(string)
}

var outgoing chan *MQTTMessage

// MQTTSubscribe ...
func MQTTSubscribe(incoming chan *MQTTMessage) {
	outgoing = incoming

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
		for _, mapping := range GetConfig().Mappings {
			m := mapping
			var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
				outgoing <- &MQTTMessage{msg, m}
			}
			c.AddRoute(mapping.MQTT.Topic, f)
		}
	}

	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
