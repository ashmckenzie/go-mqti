package mqti

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"github.com/yuin/gluamapper"
	"github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

const mQTTDefaultPort string = "1883"

// MQTTMessage ...
type MQTTMessage struct {
	Message MQTT.Message
	Mapping MappingConfiguration
	payload string
}

// Message ...
type Message struct {
	Topic   string
	Message string
}

// MessageAsStruct ...
func (m MQTTMessage) MessageAsStruct() Message {
	msg := Message{
		m.Message.Topic(),
		string(m.Message.Payload()),
	}
	return msg
}

// MessageAsJSONString ...
func (m MQTTMessage) MessageAsJSONString() (string, error) {
	var b []byte
	var err error

	if b, err = json.Marshal(m.MessageAsStruct()); err != nil {
		return "", err
	}

	return string(b), nil
}

// SetPayload ...
func (m *MQTTMessage) SetPayload(s string) {
	m.payload = s
}

// PayloadAsString ...
func (m MQTTMessage) PayloadAsString() string {
	// return string(m.Message.Payload())
	return m.payload
}

// PayloadAsJSON ...
func (m MQTTMessage) PayloadAsJSON() (map[string]interface{}, error) {
	var fields map[string]interface{}

	// err := json.Unmarshal(m.Message.Payload(), &fields)
	err := json.Unmarshal([]byte(m.payload), &fields)

	return fields, err
}

// MQTTSubscribe ...
func MQTTSubscribe(incoming chan *MQTTMessage) {
	var files []string
	var outgoing chan *MQTTMessage
	outgoing = incoming

	cs := make(chan os.Signal, 1)
	signal.Notify(cs, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-cs
		Log.Error("signal received, exiting")
		os.Exit(0)
	}()

	opts := MQTT.NewClientOptions()

	opts.ClientID = mQTTClientID()
	opts.Username = mQTTUsername()
	opts.Password = mQTTPassword()
	opts.CleanSession = mQTTCleanSession()

	if mQTTTLSDefined() {
		opts.TLSConfig = mQTTTLSConfig()
	}

	opts.AddBroker(mQTTBrokerURI())

	opts.OnConnect = func(c MQTT.Client) {
		var err error
		var config *Config

		config, err = GetConfig()
		if err != nil {
			Log.Fatal(err)
		}

		for _, mapping := range config.Mappings {
			m := mapping
			var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
				var ok bool
				match := false

				mQTTMessage := &MQTTMessage{msg, m, ""}
				message := mQTTMessage.MessageAsStruct()
				mQTTMessage.SetPayload(message.Message)

				if files, ok = mQTTMessage.luaFiles(); ok {
					for _, f := range files {
						message, ok = mQTTMessage.runLuaFile(f)
						mQTTMessage.SetPayload(message.Message)
						if ok && !match {
							match = true
						}
					}
				} else {
					match = true
				}

				Log.Debugf("match=[%v], message=[%s]", match, mQTTMessage)

				if match {
					outgoing <- mQTTMessage
				}
			}

			c.Subscribe(mapping.MQTT.Topic, 0, f)
		}
	}

	opts.OnConnectionLost = func(c MQTT.Client, e error) {
		Log.Error(e)
	}

	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		Log.Panic(token.Error())
	}

	for {
		time.Sleep(1 * time.Second)
	}
}

func (m MQTTMessage) luaFiles() ([]string, bool) {
	if len(m.Mapping.MQTT.LuaFiles) > 0 {
		return m.Mapping.MQTT.LuaFiles, true
	}
	return nil, false
}

func (m MQTTMessage) runLuaFile(f string) (Message, bool) {
	L := lua.NewState()
	luajson.Preload(L)
	defer L.Close()

	if err := L.DoFile(f); err != nil {
		panic(err)
	}

	str, _ := m.MessageAsJSONString()
	msg := m.MessageAsStruct()

	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("process"),
		NRet:    1,
		Protect: true,
	}, lua.LString(str)); err != nil {
		panic(err)
	}

	lv := L.Get(-1)
	if v, ok := lv.(*lua.LTable); ok {
		if err := gluamapper.Map(v, &msg); err != nil {
			Log.Error(err)
			return msg, false
		}
		Log.Debug(msg)
		return msg, true
	}
	return msg, false
}

func mQTTConfig() map[string]interface{} {
	return viper.GetStringMap("mqtt")
}

func mQTTBrokerURI() string {
	return fmt.Sprintf("%s://%s:%s", mQTTProtocol(), mQTTConfig()["host"], mQTTPort())
}

func mQTTPort() string {
	var port string
	if p := mQTTConfig()["port"]; p != nil {
		port = p.(string)
	} else {
		port = mQTTDefaultPort
	}
	return port
}

func mQTTProtocol() string {
	if p := mQTTConfig()["protocol"]; p != nil {
		return p.(string)
	}
	if mQTTTLSDefined() {
		return "ssl"
	}
	return "tcp"
}

func mQTTClientID() string {
	return mQTTConfig()["client_id"].(string)
}

func mQTTUsername() string {
	u := mQTTConfig()["username"]
	if u != nil {
		return u.(string)
	}
	return ""
}

func mQTTPassword() string {
	p := mQTTConfig()["password"]
	if p != nil {
		return p.(string)
	}
	return ""
}

func mQTTTLSDefined() bool {
	return mQTTConfig()["tls_cert"] != nil && mQTTConfig()["tls_private_key"] != nil
}

func mQTTTLSConfig() tls.Config {
	return *NewTLSConfig(mQTTConfig()["tls_cert"].(string), mQTTConfig()["tls_private_key"].(string))
}

func mQTTCleanSession() bool {
	return mQTTConfig()["clean_session"] != nil && (mQTTConfig()["clean_session"].(bool) == true)
}
