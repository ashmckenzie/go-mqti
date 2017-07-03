package mqti

import (
	"fmt"
	"net/url"
	"time"

	InfluxDBClient "github.com/influxdata/influxdb/client"
	"github.com/mmcloughlin/geohash"
	"github.com/spf13/viper"
)

// InfluxDBConnection ...
type InfluxDBConnection struct {
	*InfluxDBClient.Client
}

func (i InfluxDBConnection) geoHashFieldsDefined(g GeohashMungerConfiguration) bool {
	return len(g.LatitudeField) > 0 && len(g.LongitudeField) > 0 && len(g.ResultField) > 0
}

func (i InfluxDBConnection) applyGeohashMunger(g GeohashMungerConfiguration, fields map[string]interface{}, tags map[string]string) error {
	if i.geoHashFieldsDefined(g) {
		tags[g.ResultField] = geohash.Encode(
			fields[g.LatitudeField].(float64),
			fields[g.LongitudeField].(float64))
	}

	return nil
}

func (i InfluxDBConnection) applyTagsMunger(t TagsMungerConfiguration, fields map[string]interface{}, tags map[string]string) error {
	for _, x := range t.From {
		for k, v := range x {
			if fields[k] != nil {
				tags[v] = fields[k].(string)
			}
		}
	}
	return nil
}

func (i InfluxDBConnection) applyMungers(m struct {
	Tags    TagsMungerConfiguration
	Geohash GeohashMungerConfiguration
}, fields map[string]interface{}, tags map[string]string) error {
	var err error

	if err = i.applyGeohashMunger(m.Geohash, fields, tags); err != nil {
		Log.Warn(err)
	}

	if err = i.applyTagsMunger(m.Tags, fields, tags); err != nil {
		Log.Warn(err)
	}

	return err
}

// Forward ...
func (i InfluxDBConnection) Forward(m *MQTTMessage) error {
	var err error
	var fields map[string]interface{}

	config := m.MappingConfiguration.InfluxDB

	tags := config.Tags
	if tags == nil {
		tags = make(map[string]string)
	}

	fields, err = m.PayloadAsJSON()
	if err == nil {
		mungers := m.MappingConfiguration.InfluxDB.Mungers
		if err = i.applyMungers(mungers, fields, tags); err != nil {
			Log.Warn(err)
		}
	} else {
		fields = map[string]interface{}{"value": m.PayloadAsString()}
	}

	p := InfluxDBClient.Point{
		Measurement: config.Measurement,
		Tags:        tags,
		Fields:      fields,
		Time:        time.Now(),
	}

	Log.Info(p)

	_, err = i.Write(InfluxDBClient.BatchPoints{
		Points:   []InfluxDBClient.Point{p},
		Database: m.MappingConfiguration.InfluxDB.Database,
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
