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
	MappingConfiguration
}

func mQTTConfig() map[string]interface{} {
	return viper.GetStringMap("mqtt")
}

func mQTTCBrokerURI() string {
	return fmt.Sprintf("tcp://%s:%s", mQTTConfig()["host"], mQTTConfig()["port"])
}

func mQTTClientID() string {
	return mQTTConfig()["client_id"].(string)
}

func mQTTUsername() string {
	var u interface{}
	if u = mQTTConfig()["username"]; err != nil {
		return u.(string)
	}
	return ""
}

func mQTTPassword() string {
	var p interface{}
	if p = mQTTConfig()["password"]; err != nil {
		return p.(string)
	}
	return ""
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
		ClientID:  mQTTClientID(),
		Username:  mQTTUsername(),
		Password:  mQTTPassword(),
		TLSConfig: tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert},
	}

	opts.AddBroker(mQTTCBrokerURI())

	opts.OnConnect = func(c MQTT.Client) {
		for _, mapping := range GetConfig().Mappings {
			m := mapping
			var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
				outgoing <- &MQTTMessage{msg, m}
			}
			c.Subscribe(mapping.MQTT.Topic, 0, f)
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
