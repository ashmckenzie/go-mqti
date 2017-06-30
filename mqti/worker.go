package mqti

// CreateWorkers ...
func CreateWorkers(influxDB *InfluxDBConnection, jobs <-chan *MQTTMessage) {
	var err error
	var config *Config

	config, err = GetConfig()
	if err != nil {
		Log.Fatal(err)
	}

	for w := 1; w <= config.MQti.Workers; w++ {
		createWorker(w, influxDB, jobs)
	}
}

func createWorker(id int, influxDB *InfluxDBConnection, jobs <-chan *MQTTMessage) {
	var err error
	for j := range jobs {
		if err = influxDB.Forward(j); err != nil {
			Log.Error(err)
		}
	}
}
