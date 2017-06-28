package littlefly

// CreateWorkers ...
func CreateWorkers(influxDB *InfluxDBConnection, jobs <-chan *MQTTMessage) {
	for w := 1; w <= GetConfig().Littlefly.Workers; w++ {
		createWorker(w, influxDB, jobs)
	}
}

func createWorker(id int, influxDB *InfluxDBConnection, jobs <-chan *MQTTMessage) {
	for j := range jobs {
		influxDB.Forward(j)
	}
}
