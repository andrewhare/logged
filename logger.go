package logged

import (
	"bufio"
	"encoding/json"
	"io"
)

type Logger interface {
	Info(message string, data map[string]string)
	Debug(message string, data map[string]string)
}

type logger struct {
	*json.Encoder
	defaults map[string]string
}

type entry struct {
	Level   string            `json:"level"`
	Message string            `json:"message"`
	Data    map[string]string `json:"data,omitempty"`
}

func (l *logger) Info(message string, data map[string]string) {
	l.write("info", message, data)
}

func (l *logger) Debug(message string, data map[string]string) {
	l.write("debug", message, data)
}

func (l *logger) write(level, message string, data map[string]string) {
	e := &entry{
		Level:   level,
		Message: message,
		Data:    l.allData(data),
	}

	if err := l.Encode(e); err != nil {
		// what?
	}
}

func NewLogger(w io.Writer, defaults map[string]string) Logger {
	return &logger{
		Encoder:  json.NewEncoder(bufio.NewWriter(w)),
		defaults: defaults,
	}
}

func (l *logger) allData(data map[string]string) map[string]string {
	if l.defaults == nil || len(l.defaults) == 0 {
		return data
	}

	if data == nil || len(data) == 0 {
		return l.defaults
	}

	allData := make(map[string]string)

	for k, v := range l.defaults {
		allData[k] = v
	}

	for k, v := range data {
		allData[k] = v
	}

	return allData
}
