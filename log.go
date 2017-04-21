package logged

import (
	"runtime"
	"strings"
	"time"
)

const (
	Error = "error"
	Info  = "info"
	Debug = "debug"
)

type Log interface {
	Error(err error, data map[string]string) error
	Info(message string, data map[string]string) error
	Debug(message string, data map[string]string) error
	IsDebug() bool
}

type Config struct {
	Serializer    Serializer
	DebugPackages []string
	Defaults      map[string]string
}

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

func (l *log) Error(err error, data map[string]string) error {
	return l.write(Error, err.Error(), data)
}

func (l *log) Info(message string, data map[string]string) error {
	return l.write(Info, message, data)
}

func (l *log) Debug(message string, data map[string]string) error {
	if l.IsDebug() {
		return l.write(Debug, message, data)
	}
	return nil
}

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
