package littlefly

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

const debugDiskFile string = "/tmp/littlefly-debug.log"

func init() {
	setupStderrLogging()
}

// EnableDebugging ...
func EnableDebugging(yes bool) {
	var err error

	if yes {
		Log.Warnf("Debugging output will go to %s", debugDiskFile)
		DiskLog = logrus.New()
		DiskLog.Formatter = &logrus.JSONFormatter{}
		setLogLevelFor(DiskLog, logrus.DebugLevel)
		if DiskLogFile, err = os.OpenFile(debugDiskFile, os.O_CREATE|os.O_WRONLY, 0640); err != nil {
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
		debugLogIt(DiskLog, logrus.DebugLevel, line)
		// DiskLogFile.Sync()
	}
}

// LogMQTTMessage ...
func LogMQTTMessage(m *MQTTMessage) {
	fields := logrus.Fields{
		"topic":    m.Topic(),
		"mqtt":     m.Mapping.MQTT,
		"influxdb": m.Mapping.InfluxDB,
	}
	Log.WithFields(fields).Info(string(m.Payload()))
}

func debugLogIt(l *logrus.Logger, level logrus.Level, msg ...interface{}) {
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
