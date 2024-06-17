// Package log logger
package log

import (
	"log"
)

func init() {
	// defaultLogger = log.DefaultLogger
	defaultLogger = &InnerLogger{}
}

// RegisterLogger Register the logger and print the log to the set logger
func RegisterLogger(logger Logger) {
	defaultLogger = logger
}

// SetLoggerLevel Set the logger level
func SetLoggerLevel(l Level) {
	loggerLevel = l
}

// Logger Methods provided
type Logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
}

var (
	defaultLogger Logger
	loggerLevel   Level = InfoLevel
)

// Level Log Level
type Level uint32

// const ...
const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	NotLogLevel
)

// Debug Debug mode printing
func Debug(args ...interface{}) {
	if loggerLevel > DebugLevel {
		return
	}
	defaultLogger.Debug(args...)
}

// Debugf printing
func Debugf(format string, args ...interface{}) {
	if loggerLevel > DebugLevel {
		return
	}
	defaultLogger.Debugf(format, args...)
}

// Info printing
func Info(args ...interface{}) {
	if loggerLevel > InfoLevel {
		return
	}
	defaultLogger.Info(args...)
}

// Infof printing
func Infof(format string, args ...interface{}) {
	if loggerLevel > InfoLevel {
		return
	}
	defaultLogger.Infof(format, args...)
}

// Warn printing
func Warn(args ...interface{}) {
	if loggerLevel > WarnLevel {
		return
	}
	defaultLogger.Warn(args...)
}

// Warnf printing
func Warnf(format string, args ...interface{}) {
	if loggerLevel > WarnLevel {
		return
	}
	defaultLogger.Warnf(format, args...)
}

// Error printing
func Error(args ...interface{}) {
	if loggerLevel > ErrorLevel {
		return
	}
	defaultLogger.Error(args...)
}

// Errorf printing
func Errorf(format string, args ...interface{}) {
	if loggerLevel > ErrorLevel {
		return
	}
	defaultLogger.Errorf(format, args...)
}

// InnerLogger Internal logger
type InnerLogger struct{}

// Warn printing
func (i *InnerLogger) Warn(args ...interface{}) {
	log.Println(args...)
}

// Warnf printing
func (i *InnerLogger) Warnf(format string, args ...interface{}) {
	log.Printf(format+"\n", args...)
}

// Info printing
func (i *InnerLogger) Info(args ...interface{}) {
	log.Println(args...)
}

// Infof printing
func (i *InnerLogger) Infof(format string, args ...interface{}) {
	log.Printf(format+"\n", args...)
}

// Error printing
func (i *InnerLogger) Error(args ...interface{}) {
	log.Println(args...)
}

// Errorf printing
func (i *InnerLogger) Errorf(format string, args ...interface{}) {
	log.Printf(format+"\n", args...)
}

// Debug printing
func (i *InnerLogger) Debug(args ...interface{}) {
	log.Println(args...)
}

// Debugf printing
func (i *InnerLogger) Debugf(format string, args ...interface{}) {
	log.Printf(format+"\n", args...)
}
