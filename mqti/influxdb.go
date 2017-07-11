package mqti

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

// NewInfluxDBConnection ...
func NewInfluxDBConnection() (*InfluxDBConnection, error) {
	var err error
	var influxDBConn *InfluxDBClient.Client

	opts := InfluxDBClient.Config{URL: *influxDBURI()}

	if influxDBUsername() != "" && influxDBPassword() != "" {
		opts.Username = influxDBUsername()
		opts.Password = influxDBPassword()
	}

	influxDBConn, err = InfluxDBClient.NewClient(opts)
	if err != nil {
		return nil, err
	}

	return &InfluxDBConnection{influxDBConn}, nil
}

// Forward ...
func (i InfluxDBConnection) Forward(m *MQTTMessage) error {
	var err error
	var fields map[string]interface{}

	config := m.Mapping.InfluxDB

	fields, err = m.PayloadAsJSON()
	if err != nil {
		Log.Error(err)
		fields = map[string]interface{}{"value": m.PayloadAsString()}
	}

	p := InfluxDBClient.Point{
		Measurement: config.Measurement,
		Tags:        m.Tags(),
		Fields:      fields,
		Time:        time.Now(),
	}

	Log.Info(p)

	_, err = i.Write(InfluxDBClient.BatchPoints{
		Points:   []InfluxDBClient.Point{p},
		Database: m.Mapping.InfluxDB.Database,
	})

	return err
}

func influxDBConfig() map[string]interface{} {
	return viper.GetStringMap("influxdb")
}

func influxDBURI() *url.URL {
	host, _ := url.Parse(fmt.Sprintf("%s://%s:%s", influxDBProtocol(), influxDBConfig()["host"], influxDBConfig()["port"]))
	return host
}

func influxDBProtocol() string {
	t := influxDBConfig()["tls"]
	if t != nil && t.(bool) {
		return "https"
	}
	return "http"
}

func influxDBUsername() string {
	u := influxDBConfig()["username"]
	if u != nil {
		return u.(string)
	}
	return ""
}

func influxDBPassword() string {
	p := influxDBConfig()["password"]
	if p != nil {
		return p.(string)
	}
	return ""
}
