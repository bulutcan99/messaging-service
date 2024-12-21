package logger

import (
	"messaging-service/utils"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

type LogConfig struct {
	Level utils.LogLevel
}

var (
	log  *logrus.Logger
	once sync.Once
)

// InitLogger initializes the logger with given configuration
func InitLogger(logLevel utils.LogLevel) *logrus.Logger {
	once.Do(func() {
		log = logrus.New()

		// Set logger output format
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})

		// Set output to stdout
		log.SetOutput(os.Stdout)

		// Set log level based on config
		switch logLevel {
		case utils.DebugLevel:
			log.SetLevel(logrus.DebugLevel)
		case utils.InfoLevel:
			log.SetLevel(logrus.InfoLevel)
		case utils.WarningLevel:
			log.SetLevel(logrus.WarnLevel)
		case utils.ErrorLevel:
			log.SetLevel(logrus.ErrorLevel)
		case utils.FatalLevel:
			log.SetLevel(logrus.FatalLevel)
		default:
			log.SetLevel(logrus.InfoLevel)
		}
	})
	return log
}

// GetLogger returns the singleton logger instance
func GetLogger() *logrus.Logger {
	if log == nil {
		panic("Logger not initialized. Call InitLogger first")
	}
	return log
}

// Helper functions for logging
func Debug(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

func Info(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

func Warning(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

func Error(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

func Fatal(format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}

// WithFields returns a new entry with fields
func WithFields(fields map[string]interface{}) *logrus.Entry {
	return GetLogger().WithFields(logrus.Fields(fields))
}
