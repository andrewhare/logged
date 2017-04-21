package logged

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	var (
		buf bytes.Buffer
		l   = New(&Config{Serializer: NewJSONSerializer(&buf)})
		out Entry

		msg  = "a test 123"
		err  = errors.New(msg)
		data = map[string]string{"test": "123", "test2": "1111"}
	)

	l.Error(err, data)

	assert.NoError(t, json.NewDecoder(&buf).Decode(&out))
	assert.Equal(t, msg, out.Message)
	assert.Equal(t, Error, out.Level)
	assert.NotEmpty(t, out.Timestamp)
	assert.Equal(t, data, out.Data)
}

func TestInfo(t *testing.T) {
	var (
		buf bytes.Buffer
		l   = New(&Config{Serializer: NewJSONSerializer(&buf)})
		out Entry

		msg  = "a test 123"
		data = map[string]string{"test": "123", "test2": "1111"}
	)

	l.Info(msg, data)

	assert.NoError(t, json.NewDecoder(&buf).Decode(&out))
	assert.Equal(t, msg, out.Message)
	assert.Equal(t, Info, out.Level)
	assert.NotEmpty(t, out.Timestamp)
	assert.Equal(t, data, out.Data)
}

func TestDebug(t *testing.T) {
	var (
		buf bytes.Buffer
		l   = New(&Config{
			Serializer:    NewJSONSerializer(&buf),
			DebugPackages: []string{"github.com/andrewhare/logged"},
		})
		out Entry

		msg  = "a test 123"
		data = map[string]string{"test": "123", "test2": "1111"}
	)

	l.Debug(msg, data)

	assert.NoError(t, json.NewDecoder(&buf).Decode(&out))
	assert.Equal(t, msg, out.Message)
	assert.Equal(t, Debug, out.Level)
	assert.NotEmpty(t, out.Timestamp)
	assert.Equal(t, data, out.Data)
}

func TestMergedDataNoDefaults(t *testing.T) {
	var (
		data = map[string]string{"one": "1"}
		l    = &log{}
	)

	assert.Nil(t, l.mergedData(nil))
	assert.Equal(t, data, l.mergedData(data))
}

func TestMergedDataNoData(t *testing.T) {
	var (
		defaults = map[string]string{"one": "1"}
		l        = &log{defaults: defaults}
	)

	assert.Equal(t, defaults, l.mergedData(nil))
}

func TestMergedData(t *testing.T) {
	var (
		defaults = map[string]string{"one": "1"}
		data     = map[string]string{"two": "2"}
		out      = map[string]string{"one": "1", "two": "2"}
		l        = &log{defaults: defaults}
	)

	assert.Equal(t, out, l.mergedData(data))
}

func TestMergedDataOverwriteDefaults(t *testing.T) {
	var (
		defaults = map[string]string{"one": "1"}
		data     = map[string]string{"one": "111"}
		l        = &log{defaults: defaults}
	)

	assert.Equal(t, data, l.mergedData(data))
}
