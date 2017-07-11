package mqti

import "github.com/spf13/viper"

type mQTTMappingConfiguration struct {
	Topic    string
	LuaFiles []string `mapstructure:"lua_files"`
}

type influxDBMappingConfiguration struct {
	Database    string
	Measurement string
	Tags        map[string]string
}

// Config ...
type Config struct {
	MQti     mQtiConfiguration
	MQTT     mQTTConfiguration
	InfluxDB influxDBConfiguration
	Mappings []MappingConfiguration
}

// MappingConfiguration ...
type MappingConfiguration struct {
	Name     string
	MQTT     mQTTMappingConfiguration
	InfluxDB influxDBMappingConfiguration
}

// GetConfig ...
func GetConfig() (*Config, error) {
	var err error
	var c Config

	if err = viper.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, err
}
