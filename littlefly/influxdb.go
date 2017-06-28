package littlefly

import (
	"fmt"
	"net/url"
	"time"

	InfluxDBClient "github.com/influxdata/influxdb/client"
	"github.com/spf13/viper"
)

// InfluxDBConnection ...
type InfluxDBConnection struct {
	*InfluxDBClient.Client
}

// Forward ...
func (i InfluxDBConnection) Forward(m *MQTTMessage) error {
	pts := make([]InfluxDBClient.Point, 1)

	measurement := m.MappingConfiguration.InfluxDB.Measurement
	tags := m.MappingConfiguration.InfluxDB.Tags
	value := string(m.Payload())

	pts[0] = InfluxDBClient.Point{
		Measurement: measurement,
		Tags:        tags,
		Fields:      map[string]interface{}{"value": value},
		Time:        time.Now(),
	}

	bps := InfluxDBClient.BatchPoints{Points: pts, Database: m.MappingConfiguration.InfluxDB.Database}

	_, err := i.Write(bps)
	if err != nil {
		return err
	}

	return nil
}

func influxDBConfig() map[string]interface{} {
	return viper.GetStringMap("influxdb")
}

func influxDBURI() *url.URL {
	host, _ := url.Parse(fmt.Sprintf("http://%s:%s", influxDBConfig()["host"], influxDBConfig()["port"]))
	return host
}

// NewInfluxDBConnection ...
func NewInfluxDBConnection() (*InfluxDBConnection, error) {
	var err error
	var influxDBConn *InfluxDBClient.Client

	influxDBConn, err = InfluxDBClient.NewClient(InfluxDBClient.Config{URL: *influxDBURI()})
	if err != nil {
		return nil, err
	}

	return &InfluxDBConnection{influxDBConn}, nil
}
