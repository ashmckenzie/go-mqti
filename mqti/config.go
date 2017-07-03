package mqti

type mQtiConfiguration struct {
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
