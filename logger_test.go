package logged

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfo(t *testing.T) {
	var (
		buf bytes.Buffer
		l   = New(&buf, nil)
		out entry

		msg  = "a test 123"
		data = Data{"test": "123", "test2": []interface{}{float64(1111)}}
	)

	l.Info(msg, data)

	assert.NoError(t, json.NewDecoder(&buf).Decode(&out))
	assert.Equal(t, msg, out.Message)
	assert.Equal(t, Info, out.Level)
	assert.NotEmpty(t, out.Timestamp)
	assert.Equal(t, data, out.Data)
}

func TestInfoFailedEncode(t *testing.T) {
	var (
		buf  bytes.Buffer
		l    = New(&buf, nil)
		out  failure
		data = Data{"test": make(chan int)}
	)

	l.Info("", data)

	assert.NoError(t, json.NewDecoder(&buf).Decode(&out))
	assert.Contains(t, out.Error, "json: unsupported type: chan int")
	assert.Contains(t, out.Data, "(chan int)")
}

func TestDebug(t *testing.T) {
	var (
		buf bytes.Buffer
		l   = New(&buf, nil)
		out entry

		msg  = "a test 123"
		data = Data{"test": "123", "test2": []interface{}{float64(1111)}}
	)

	l.Debug(msg, data)

	assert.NoError(t, json.NewDecoder(&buf).Decode(&out))
	assert.Equal(t, msg, out.Message)
	assert.Equal(t, Debug, out.Level)
	assert.NotEmpty(t, out.Timestamp)
	assert.Equal(t, data, out.Data)
}

func TestDebugFailedEncode(t *testing.T) {
	var (
		buf  bytes.Buffer
		l    = New(&buf, nil)
		out  failure
		data = Data{"test": make(chan int)}
	)

	l.Debug("", data)

	assert.NoError(t, json.NewDecoder(&buf).Decode(&out))
	assert.Contains(t, out.Error, "json: unsupported type: chan int")
	assert.Contains(t, out.Data, "(chan int)")
}

func TestMergedDataNoDefaults(t *testing.T) {
	var (
		data = Data{"one": "1"}
		l    = &logger{}
	)

	assert.Nil(t, l.mergedData(nil))
	assert.Equal(t, data, l.mergedData(data))
}

func TestMergedDataNoData(t *testing.T) {
	var (
		defaults = Data{"one": "1"}
		l        = &logger{defaults: defaults}
	)

	assert.Equal(t, defaults, l.mergedData(nil))
}

func TestMergedData(t *testing.T) {
	var (
		defaults = Data{"one": "1"}
		data     = Data{"two": "2"}
		out      = Data{"one": "1", "two": "2"}
		l        = &logger{defaults: defaults}
	)

	assert.Equal(t, out, l.mergedData(data))
}

func TestMergedDataOverwriteDefaults(t *testing.T) {
	var (
		defaults = Data{"one": "1"}
		data     = Data{"one": "111"}
		l        = &logger{defaults: defaults}
	)

	assert.Equal(t, data, l.mergedData(data))
}
