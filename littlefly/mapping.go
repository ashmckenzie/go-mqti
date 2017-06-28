package littlefly

import (
	"github.com/spf13/viper"
)

type littleflyConfiguration struct {
	Workers int
}

type mQTTConfiguration struct {
	Host     string
	Port     string
	ClientID string
}

type influxDBConfiguration struct {
	Host string
	Port string
}

type mQTTMappingConfiguration struct {
	Topic string
}

type influxDBMappingConfiguration struct {
	Database    string
	Measurement string
	Tags        map[string]string
}

// MappingConfiguration ...
type MappingConfiguration struct {
	Name     string
	MQTT     mQTTMappingConfiguration
	InfluxDB influxDBMappingConfiguration
}

// Config ...
type Config struct {
	Littlefly littleflyConfiguration
	MQTT      mQTTConfiguration
	InfluxDB  influxDBConfiguration
	Mappings  []MappingConfiguration
}

// GetConfig ...
func GetConfig() *Config {
	var c Config
	viper.Unmarshal(&c)

	return &c
}
