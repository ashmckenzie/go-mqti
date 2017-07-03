# mqti

[![Build Status](https://travis-ci.org/ashmckenzie/go-mqti.svg?branch=master)](https://travis-ci.org/ashmckenzie/go-mqti)
[![Go Report Card](https://goreportcard.com/badge/github.com/ashmckenzie/go-mqti)](https://goreportcard.com/report/github.com/ashmckenzie/go-mqti)

MQTT subscriber that pumps data into InfluxDB.

Pronounced 'm-cutey' :)

## Features

* MQTT 3.1.1 supported, TLS, username/password
* Consume MQTT messages and inspect (`watch`) or `forward` with the following abilities:
  * Filter messages with AND + OR
* Receive MQTT messages and write into InfluxDB, with the following abilities:
  * Add tags based on MQTT fields (when MQTT payload is JSON)
  * Geohash support (applicable when consuming MQTT messages from [Owntracks](http://owntracks.org/)
* Includes `docker-compose.yaml` to get a full setup up and running!

## Configuration

Configuration is handled through a `config.yaml` file.  The following example reads as:

* Setup four workers for incoming MQTT messages
* Consume message from MQTT server `tcp://localhost:1883` with the client ID of `mqti`
* When a MQTT message is consumed from the `temperature` topic, send write requests to InfluxDB server `http://localhost:8086`, into the `iot` database as measurement `temperature`

```yaml
---
mqti:
  workers: 4

mqtt:
  host: "localhost"
  port: "1883"
  client_id: "mqti"

influxdb:
  host: "localhost"
  port: "8086"

mappings:
  - mqtt:
      topic: "temperature"
    influxdb:
      database: "iot"
      measurement: "temperature"
```

## Install

`go get github.com/ashmckenzie/go-mqti/mqti`

or download a release:

[github.com/ashmckenzie/go-mqti/releases](https://github.com/ashmckenzie/go-mqti/releases)

## Usage

1. Ensure you have a `config.yaml` setup (see above)

### To consume MQTT messages only

1. `$GOPATH/bin/mqti watch`

### To consume MQTT messages *and* forward to InfluxDB

1. `$GOPATH/bin/mqti forward`

## Trying out with Docker

Ensure you're into the root directory and then type in `make run`

* InfluxDB UI http://localhost:8083/ (root/root)
* Grafana http://localhost:3000/ (admin/admin)

## Help

```shell
MQTT subscriber that pumps data into InfluxDB

Usage:
  mqti [flags]
  mqti [command]

Available Commands:
  forward     Forward MQTT messages on to InfluxDB
  watch       Watch MQTT messages
  help        Help about any command

Flags:
      --config string   config file (default is config.yaml)
      --debug           enable debugging
  -h, --help            help for mqti
  -v, --version         show version

Use "mqti [command] --help" for more information about a command.
```

## Building

1. `make`

## Thanks

* [github.com/eclipse/paho.mqtt.golang](https://github.com/eclipse/paho.mqtt.golang)
* [github.com/influxdata/influxdb/client](https://github.com/influxdata/influxdb/client)
* [github.com/spf13/cobra](https://github.com/spf13/cobra)
* [github.com/spf13/viper](https://github.com/spf13/viper)
* [github.com/Sirupsen/logrus](https://github.com/Sirupsen/logrus)

## License

MIT License

Copyright (c) 2017 Ash McKenzie

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
