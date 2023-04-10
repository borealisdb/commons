package logger

import (
	"github.com/sirupsen/logrus"
)

var logLevels = map[string]logrus.Level{
	"info":    logrus.InfoLevel,
	"debug":   logrus.DebugLevel,
	"warning": logrus.WarnLevel,
}

type DefaultLogger struct {
	logger *logrus.Entry
}

func NewDefaultLogger(logLevelRaw string, component string) *logrus.Entry {
	logLevel, ok := logLevels[logLevelRaw]
	if !ok {
		logLevel = logrus.InfoLevel
	}
	l := logrus.New()
	l.SetLevel(logLevel)
	logger := logrus.NewEntry(l)
	prefixedLogger := logger.WithField("component", component)
	return prefixedLogger
}
