package oncall

import "github.com/sirupsen/logrus"

type LeveledLogger interface {
	WithField(key string, value interface{}) LeveledLogger
	Trace(...interface{})
	Tracef(format string, values ...interface{})
	Debug(...interface{})
	Debugf(format string, values ...interface{})
	Info(...interface{})
	Infof(format string, values ...interface{})
	Warn(...interface{})
	Warnf(format string, values ...interface{})
	Error(...interface{})
	Errorf(format string, values ...interface{})
	Fatal(...interface{})
	Fatalf(format string, values ...interface{})
}

type DefaultLogger struct {
	fields map[string]interface{}
}

var log LeveledLogger = DefaultLogger{}

func (l DefaultLogger) WithField(key string, value interface{}) LeveledLogger {
	if l.fields == nil {
		l.fields = make(map[string]interface{})
	}
	l.fields[key] = value
	return l
}

func (l DefaultLogger) Trace(a ...interface{}) {
	logrus.WithFields(logrus.Fields(l.fields)).Trace(a...)
}
func (l DefaultLogger) Tracef(format string, values ...interface{}) {
	logrus.WithFields(logrus.Fields(l.fields)).Tracef(format, values...)
}

func (l DefaultLogger) Debug(a ...interface{}) {
	logrus.WithFields(logrus.Fields(l.fields)).Debug(a...)
}
func (l DefaultLogger) Debugf(format string, values ...interface{}) {
	logrus.WithFields(logrus.Fields(l.fields)).Debugf(format, values...)
}
func (l DefaultLogger) Info(a ...interface{}) {
	logrus.WithFields(logrus.Fields(l.fields)).Info(a...)
}
func (l DefaultLogger) Infof(format string, values ...interface{}) {
	logrus.WithFields(logrus.Fields(l.fields)).Infof(format, values...)
}
func (l DefaultLogger) Warn(a ...interface{}) {
	logrus.WithFields(logrus.Fields(l.fields)).Warn(a...)
}
func (l DefaultLogger) Warnf(format string, values ...interface{}) {
	logrus.WithFields(logrus.Fields(l.fields)).Warnf(format, values...)
}
func (l DefaultLogger) Error(a ...interface{}) {
	logrus.WithFields(logrus.Fields(l.fields)).Error(a...)
}
func (l DefaultLogger) Errorf(format string, values ...interface{}) {
	logrus.WithFields(logrus.Fields(l.fields)).Errorf(format, values...)
}
func (l DefaultLogger) Fatal(a ...interface{}) {
	logrus.WithFields(logrus.Fields(l.fields)).Fatal(a...)
}
func (l DefaultLogger) Fatalf(format string, values ...interface{}) {
	logrus.WithFields(logrus.Fields(l.fields)).Fatalf(format, values...)
}
