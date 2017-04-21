package logged

import (
	"runtime"
	"strings"
	"time"
)

const (
	// Error log level
	Error = "error"
	// Info log level
	Info = "info"
	// Debug log level
	Debug = "debug"
)

// Log provides functions to write messages and data to an output
type Log interface {
	Error(err error, data map[string]string) error
	Info(message string, data map[string]string) error
	Debug(message string, data map[string]string) error
	IsDebug() bool
}

// Config is a struct used to initialize a Log
type Config struct {
	Serializer    Serializer
	DebugPackages []string
	Defaults      map[string]string
}

// New creates a new Log
func New(c *Config) Log {
	return &log{
		serializer:    c.Serializer,
		debugPackages: c.DebugPackages,
		defaults:      c.Defaults,
	}
}

type log struct {
	serializer    Serializer
	defaults      map[string]string
	debugPackages []string
}

// Error writes a log entry at the Error level
func (l *log) Error(err error, data map[string]string) error {
	return l.write(Error, err.Error(), data)
}

// Info writes a log entry at the Info level
func (l *log) Info(message string, data map[string]string) error {
	return l.write(Info, message, data)
}

// Debug writes a log entry at the Debug level if IsDebug
func (l *log) Debug(message string, data map[string]string) error {
	if l.IsDebug() {
		return l.write(Debug, message, data)
	}
	return nil
}

// IsDebug determines if the log is configured to write debug entries
func (l *log) IsDebug() bool {
	if len(l.debugPackages) == 0 {
		return false
	}

	if l.debugPackages[0] == "*" {
		return true
	}

	pc, _, _, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()

	for _, pkg := range l.debugPackages {
		if strings.HasPrefix(funcName, pkg) {
			return true
		}
	}

	return false
}

func (l *log) write(level, message string, data map[string]string) error {
	return l.serializer.Write(&Entry{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Level:     level,
		Message:   message,
		Data:      l.mergedData(data),
	})
}

func (l *log) mergedData(data map[string]string) map[string]string {
	if l.defaults == nil || len(l.defaults) == 0 {
		return data
	}

	if data == nil || len(data) == 0 {
		return l.defaults
	}

	merged := make(map[string]string)
	for k, v := range l.defaults {
		merged[k] = v
	}
	for k, v := range data {
		merged[k] = v
	}

	return merged
}
