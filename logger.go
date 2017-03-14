package logged

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

const (
	Info  = "info"
	Debug = "debug"
)

type Data map[string]interface{}

type Logger interface {
	Info(message string, data Data)
	Debug(message string, data Data)
}

func New(w io.Writer, defaults Data) Logger {
	return &logger{
		Encoder:  json.NewEncoder(w),
		defaults: defaults,
	}
}

type logger struct {
	*json.Encoder
	defaults Data
}

type entry struct {
	Level     string    `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Data      Data      `json:"data,omitempty"`
}

type failure struct {
	Error string `json:"logger_error"`
	Data  string `json:"data"`
}

func (l *logger) Info(message string, data Data) {
	l.write(Info, message, data)
}

func (l *logger) Debug(message string, data Data) {
	l.write(Debug, message, data)
}

func (l *logger) write(level, message string, data Data) {
	e := entry{
		Level:     level,
		Timestamp: time.Now().UTC(),
		Message:   message,
		Data:      l.mergedData(data),
	}

	if err := l.Encode(e); err != nil {
		l.Encode(failure{
			Error: err.Error(),
			Data:  fmt.Sprintf("%#v", e.Data),
		})
	}
}

func (l *logger) mergedData(data Data) Data {
	if l.defaults == nil || len(l.defaults) == 0 {
		return data
	}

	if data == nil || len(data) == 0 {
		return l.defaults
	}

	mergedData := make(Data)

	for k, v := range l.defaults {
		mergedData[k] = v
	}

	for k, v := range data {
		mergedData[k] = v
	}

	return mergedData
}
