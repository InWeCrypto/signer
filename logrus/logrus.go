package logrus

import (
	log "github.com/Sirupsen/logrus"
	"github.com/goany/slf4go"
)

// Logger facade for logrus
type Logger struct {
	name  string
	entry *log.Entry
}

func newLogger(name string, logger *log.Logger) *Logger {
	result := &Logger{}
	result.name = name
	result.entry = log.NewEntry(logger).WithField("name", name)
	return result
}

// GetName .
func (logger *Logger) GetName() string {
	return logger.name
}

// Trace .
func (logger *Logger) Trace(args ...interface{}) {
	// forward to Debug
	logger.Debug(args)
}

// TraceF .
func (logger *Logger) TraceF(format string, args ...interface{}) {
	// forward to Debug
	logger.DebugF(format, args)
}

// Debug .
func (logger *Logger) Debug(args ...interface{}) {
	logger.entry.Debugln(args)
}

// DebugF .
func (logger *Logger) DebugF(format string, args ...interface{}) {
	logger.entry.Debugf(format, args)
}

// Info .
func (logger *Logger) Info(args ...interface{}) {
	logger.entry.Infoln(args)
}

// InfoF .
func (logger *Logger) InfoF(format string, args ...interface{}) {
	logger.entry.Infof(format, args)
}

// Warn .
func (logger *Logger) Warn(args ...interface{}) {
	logger.entry.Warnln(args)
}

// WarnF .
func (logger *Logger) WarnF(format string, args ...interface{}) {
	logger.entry.Warnf(format, args)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.entry.Errorln(args)
}

// ErrorF .
func (logger *Logger) ErrorF(format string, args ...interface{}) {
	logger.entry.Errorf(format, args)
}

// Fatal .
func (logger *Logger) Fatal(args ...interface{}) {
	logger.entry.Fatalln(args)
}

// FatalF .
func (logger *Logger) FatalF(format string, args ...interface{}) {
	logger.entry.Fatalf(format, args)
}

// LoggerFactory for logrus
type LoggerFactory struct {
	logger *log.Logger
}

// GetLogger .
func (factory *LoggerFactory) GetLogger(name string) slf4go.Logger {
	return newLogger(name, factory.logger)
}

// NewLoggerFactory .
func NewLoggerFactory(logger *log.Logger) *LoggerFactory {
	factory := &LoggerFactory{}
	factory.logger = logger
	return factory
}
