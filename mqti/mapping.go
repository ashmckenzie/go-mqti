package mqti

import "github.com/spf13/viper"

type mQTTMappingConfiguration struct {
	Topic   string
	LUAFile string `mapstructure:"lua_file"`
	Mungers struct {
		Filter FilterMungerConfiguration `mapstructure:"filter"`
	}
}

type influxDBMappingConfiguration struct {
	Database    string
	Measurement string
	Tags        map[string]string
	Mungers     struct {
		Tags    TagsMungerConfiguration
		Geohash GeohashMungerConfiguration
	}
}

// FilterMungerConfiguration ...
type FilterMungerConfiguration struct {
	JSON FilterJSONMungerConfiguration
}

// FilterJSONMungerConfiguration ...
type FilterJSONMungerConfiguration struct {
	And []map[string]string
	Or  []map[string]string
}

// TagsMungerConfiguration ...
type TagsMungerConfiguration struct {
	From []map[string]string
}

// GeohashMungerConfiguration ...
type GeohashMungerConfiguration struct {
	LatitudeField  string `mapstructure:"lat_field"`
	LongitudeField string `mapstructure:"lng_field"`
	ResultField    string `mapstructure:"result_field"`
}

// MappingConfiguration ...
type MappingConfiguration struct {
	Name     string
	MQTT     mQTTMappingConfiguration
	InfluxDB influxDBMappingConfiguration
}

// Config ...
type Config struct {
	MQti     mQtiConfiguration
	MQTT     mQTTConfiguration
	InfluxDB influxDBConfiguration
	Mappings []MappingConfiguration
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
