package logged

import (
	"runtime"
	"strings"
	"time"
)

const (
	// Info log level
	Info = "info"
	// Debug log level
	Debug = "debug"
)

// Log provides functions to write messages and data to an output
type Log interface {
	Info(message string, data map[string]string) error
	Debug(message string, data map[string]string) error

	IsDebug() bool
	New(data map[string]string) Log
}

// Opts is a struct used provide optional values to the log
type Opts struct {
	DebugPackages []string
	Defaults      map[string]string
}

// New creates a new Log
func New(s Serializer) Log {
	return &log{serializer: s}
}

// NewOpts creates a new log with options
func NewOpts(s Serializer, o Opts) Log {
	return &log{
		serializer:    s,
		debugPackages: o.DebugPackages,
		defaults:      o.Defaults,
	}
}

type log struct {
	serializer    Serializer
	defaults      map[string]string
	debugPackages []string
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

func (l *log) New(data map[string]string) Log {
	return &log{
		serializer:    l.serializer,
		defaults:      l.mergedData(data),
		debugPackages: l.debugPackages,
	}
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
