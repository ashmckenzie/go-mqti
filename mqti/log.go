package mqti

import (
	"fmt"
	"os"
	"runtime"

	"github.com/Sirupsen/logrus"
)

// Log ...
var Log = logrus.New()

// DiskLog ...
var DiskLog *logrus.Logger

// DiskLogFile ...
var DiskLogFile *os.File

// DEBUGDISKFILE ...
const DEBUGDISKFILE string = "/tmp/ofr-debug.log"

func init() {
	setupStderrLogging()
}

// EnableDebugging ...
func EnableDebugging(yes bool) {
	var err error

	if yes {
		Log.Infof("Debugging output will go to %s", DEBUGDISKFILE)
		DiskLog = logrus.New()
		setLogLevelFor(DiskLog, logrus.DebugLevel)
		DiskLog.Formatter = &logrus.JSONFormatter{}
		if DiskLogFile, err = os.OpenFile(DEBUGDISKFILE, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600); err != nil {
			Log.Panic(err)
		}
		DiskLog.Out = DiskLogFile
	}
}

func setupStderrLogging() {
	Log.Out = os.Stderr
	setLogLevelFor(Log, logrus.InfoLevel)
}

func setLogLevelFor(l *logrus.Logger, level logrus.Level) {
	l.Level = level
}

// DebugLog ...
func DebugLog(line ...interface{}) {
	if DiskLog != nil {
		logIt(DiskLog, logrus.DebugLevel, line)
	}
}

// LogMQTTMessage ...
func LogMQTTMessage(m *MQTTMessage) {
	fields := logrus.Fields{
		"topic":    m.Topic(),
		"mqtt":     m.MappingConfiguration.MQTT,
		"influxdb": m.MappingConfiguration.InfluxDB,
	}
	Log.WithFields(fields).Info(string(m.Payload()))
}

func logIt(l *logrus.Logger, level logrus.Level, msg ...interface{}) {
	pc, _, _, _ := runtime.Caller(2)
	details := runtime.FuncForPC(pc)
	fileFunc, lineFunc := details.FileLine(pc)
	location := fmt.Sprintf("%s:%d", fileFunc, lineFunc-2)
	msgAsString := fmt.Sprintf("%s", msg)
	fields := logrus.Fields{"location": location}

	switch level {
	case logrus.DebugLevel:
		l.WithFields(fields).Debug(msgAsString)
		break
	case logrus.InfoLevel:
		l.WithFields(fields).Info(msgAsString)
		break
	}
}
