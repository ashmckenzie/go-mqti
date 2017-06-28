package littlefly

import (
	"github.com/spf13/viper"
)

// MQTTMapping ...
type MQTTMapping struct {
	Topic string
}

// InfluxDBMapping ...
type InfluxDBMapping struct {
	Database    string
	Measurement string
	Tags        []string
}

// Mapping ...
type Mapping struct {
	Name     string
	MQTT     MQTTMapping
	InfluxDB InfluxDBMapping
}

// Config ...
type Config struct {
	Mappings []Mapping
}

// GetConfig ...
func GetConfig() *Config {
	var c Config
	viper.Unmarshal(&c)

	return &c
}
