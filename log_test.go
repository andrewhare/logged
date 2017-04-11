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
		l   = New(&Config{Writer: &buf})
		out entry

		msg  = "a test 123"
		data = Data{"test": "123", "test2": "1111"}
	)

	l.InfoEx(msg, data)

	assert.NoError(t, json.NewDecoder(&buf).Decode(&out))
	assert.Equal(t, msg, out.Message)
	assert.Equal(t, Info, out.Level)
	assert.NotEmpty(t, out.Timestamp)
	assert.Equal(t, data, out.Data)
}

func TestDebug(t *testing.T) {
	var (
		buf bytes.Buffer
		l   = New(&Config{Writer: &buf})
		out entry

		msg  = "a test 123"
		data = Data{"test": "123", "test2": "1111"}
	)

	l.DebugEx(msg, data)

	assert.NoError(t, json.NewDecoder(&buf).Decode(&out))
	assert.Equal(t, msg, out.Message)
	assert.Equal(t, Debug, out.Level)
	assert.NotEmpty(t, out.Timestamp)
	assert.Equal(t, data, out.Data)
}

func TestMergedDataNoDefaults(t *testing.T) {
	var (
		data = Data{"one": "1"}
		l    = &log{}
	)

	assert.Nil(t, l.mergedData(nil))
	assert.Equal(t, data, l.mergedData(data))
}

func TestMergedDataNoData(t *testing.T) {
	var (
		defaults = Data{"one": "1"}
		l        = &log{defaults: defaults}
	)

	assert.Equal(t, defaults, l.mergedData(nil))
}

func TestMergedData(t *testing.T) {
	var (
		defaults = Data{"one": "1"}
		data     = Data{"two": "2"}
		out      = Data{"one": "1", "two": "2"}
		l        = &log{defaults: defaults}
	)

	assert.Equal(t, out, l.mergedData(data))
}

func TestMergedDataOverwriteDefaults(t *testing.T) {
	var (
		defaults = Data{"one": "1"}
		data     = Data{"one": "111"}
		l        = &log{defaults: defaults}
	)

	assert.Equal(t, data, l.mergedData(data))
}
