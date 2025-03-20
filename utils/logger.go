package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Logger provides a simple logging interface
type Logger interface {
	Info(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
	Debug(msg string, keyvals ...interface{})
}

// SimpleLogger implements the Logger interface
type SimpleLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
}

// NewLogger creates a new SimpleLogger
func NewLogger() Logger {
	infoLogger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	debugLogger := log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime)

	return &SimpleLogger{
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
		debugLogger: debugLogger,
	}
}

// formatKeyvals formats key-value pairs for logging
func formatKeyvals(keyvals ...interface{}) string {
	result := ""
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			result += " " + string(keyvals[i].(string)) + "=" + formatValue(keyvals[i+1])
		}
	}
	return result
}

// formatValue converts a value to a string
func formatValue(value interface{}) string {
	switch v := value.(type) {
	case time.Duration:
		return v.String()
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Info logs info level messages
func (l *SimpleLogger) Info(msg string, keyvals ...interface{}) {
	l.infoLogger.Println(msg + formatKeyvals(keyvals...))
}

// Error logs error level messages
func (l *SimpleLogger) Error(msg string, keyvals ...interface{}) {
	l.errorLogger.Println(msg + formatKeyvals(keyvals...))
}

// Debug logs debug level messages
func (l *SimpleLogger) Debug(msg string, keyvals ...interface{}) {
	l.debugLogger.Println(msg + formatKeyvals(keyvals...))
}
