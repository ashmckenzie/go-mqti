littlefly
=========

[![Build Status](https://travis-ci.org/ashmckenzie/go-litlefly.svg?branch=master)](https://travis-ci.org/ashmckenzie/go-litlefly)

MQTT subscriber that pumps data into InfluxDB.

```
<USAGE>
```

Configuration
-------------

Configuration is handled through a `config.toml` file.  Example:

```
[mqtt]
host = 'localhost'
port = 1883
topic = 'mqtt'
client_id = 'littlefly'

[influxdb]
host = 'localhost'
port = 8086
database = 'mqtt'
```

Install
-------

`go get github.com/ashmckenzie/go-littlefly/littlefly`

or download a release:

https://github.com/ashmckenzie/go-littlefly/releases

License
-------

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
