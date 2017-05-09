package logged

import (
	"fmt"
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
	Info(message string, data ...map[string]string)
	InfoError(err error, data ...map[string]string)
	Debug(message string, data ...map[string]string)
	DebugError(err error, data ...map[string]string)

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
func (l *log) Info(message string, data ...map[string]string) {
	l.write(Info, message, data...)
}

// InfoError writes an error log entry at the Info level
func (l *log) InfoError(err error, data ...map[string]string) {
	if err != nil {
		l.write(Info, err.Error(), data...)
	}
}

// Debug writes a log entry at the Debug level if IsDebug
func (l *log) Debug(message string, data ...map[string]string) {
	if l.IsDebug() {
		l.write(Debug, message, data...)
	}
}

// DebugError writes an error log entry at the Debug level
func (l *log) DebugError(err error, data ...map[string]string) {
	if err != nil {
		l.write(Debug, err.Error(), data...)
	}
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

func (l *log) write(level, message string, data ...map[string]string) {
	err := l.serializer.Write(&Entry{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Level:     level,
		Message:   message,
		Data:      l.mergedData(data...),
	})

	if err != nil {
		fmt.Printf("logged: failed to write: %s\n", err)
	}
}

func (l *log) mergedData(data ...map[string]string) map[string]string {
	if len(data) == 0 {
		return l.defaults
	}

	merged := make(map[string]string)
	if l.defaults != nil || len(l.defaults) > 0 {
		for k, v := range l.defaults {
			merged[k] = v
		}
	}
	for _, d := range data {
		for k, v := range d {
			merged[k] = v
		}
	}

	return merged
}
