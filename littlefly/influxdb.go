package littlefly

import (
	"fmt"
	"net/url"

	InfluxDBClient "github.com/influxdata/influxdb/client"
	"github.com/spf13/viper"
)

// InfluxDBConnection ...
type InfluxDBConnection struct {
	*InfluxDBClient.Client
}

// InfluxDBPoint ...
type InfluxDBPoint struct {
	// InfluxDBClient.Point
	Value       string
	Measurement string
	Tags        map[string]string
}

// Forward ...
func (i InfluxDBConnection) Forward(p *InfluxDBPoint) error {
	// pts := make([]InfluxDBClient.Point, 1)
	//
	// pts[0] := InfluxDBClient.Point{
	// 	Measurement: "cpu_load_short",
	// 	Tags:        map[string]string{},
	// 	Fields:      map[string]interface{}{"value": },
	// 	Time:        time.Now(),
	// }

	//
	// bps := InfluxDBClient.BatchPoints{
	// 	Points:   pts,
	// 	Database: InfluxDBDatabase(),
	// }
	//
	// _, err := i.Write(bps)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// InfluxDBConfig ...
func InfluxDBConfig() map[string]interface{} {
	return viper.GetStringMap("influxdb")
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

// InfluxDBTags ...
func InfluxDBTags() []string {
	return InfluxDBConfig()["tags"].([]string)
}

// NewInfluxDBConnection ...
func NewInfluxDBConnection() (*InfluxDBConnection, error) {
	var err error
	var influxDBConn *InfluxDBClient.Client

	influxDBConn, err = InfluxDBClient.NewClient(InfluxDBClient.Config{URL: *InfluxDBURI()})
	if err != nil {
		return nil, err
	}

	return &InfluxDBConnection{influxDBConn}, nil
}
