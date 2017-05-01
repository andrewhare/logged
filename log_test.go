package logged

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfo(t *testing.T) {
	var (
		buf bytes.Buffer
		l   = New(NewJSONSerializer(&buf))
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
		l   = NewOpts(NewJSONSerializer(&buf), Opts{
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

	assert.Empty(t, l.mergedData(nil))
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

func TestLogNew(t *testing.T) {
	var (
		l           = &log{defaults: map[string]string{"one": "1"}}
		newDefaults = map[string]string{"two": "2"}
		data        = map[string]string{"three": "3"}
		out         = map[string]string{"one": "1", "two": "2", "three": "3"}
	)

	assert.Equal(t, out, l.New(newDefaults).(*log).mergedData(data))
}

func TestAPI(t *testing.T) {
	ser := NewJSONSerializer(os.Stdout)
	log := NewOpts(ser, Opts{
		Defaults:      map[string]string{"a": "1"},
		DebugPackages: []string{"github.com/andrewhare/logged"},
	})

	log.Info("test")
	log.Info("test", map[string]string{"b": "2"})
	log.Info("test", map[string]string{"b": "2"}, map[string]string{"c": "3"})

	log.InfoError(fmt.Errorf("boom"))
	log.InfoError(fmt.Errorf("boom"), map[string]string{"b": "2"})
	log.InfoError(fmt.Errorf("boom"), map[string]string{"b": "2"}, map[string]string{"c": "3"})

	log.Debug("test")
	log.Debug("test", map[string]string{"b": "2"})
	log.Debug("test", map[string]string{"b": "2"}, map[string]string{"c": "3"})

	log.DebugError(fmt.Errorf("boom"))
	log.DebugError(fmt.Errorf("boom"), map[string]string{"b": "2"})
	log.DebugError(fmt.Errorf("boom"), map[string]string{"b": "2"}, map[string]string{"c": "3"})
}
