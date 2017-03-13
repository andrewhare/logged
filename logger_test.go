package logged

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllDataNoDefaults(t *testing.T) {
	var (
		data = map[string]string{"one": "1"}
		l    = &logger{}
	)

	assert.Nil(t, l.allData(nil))
	assert.Equal(t, data, l.allData(data))
}

func TestAllDataNoData(t *testing.T) {
	var (
		defaults = map[string]string{"one": "1"}
		l        = &logger{defaults: defaults}
	)

	assert.Equal(t, defaults, l.allData(nil))
}

func TestAllData(t *testing.T) {
	var (
		defaults = map[string]string{"one": "1"}
		data     = map[string]string{"two": "2"}
		out      = map[string]string{"one": "1", "two": "2"}
		l        = &logger{defaults: defaults}
	)

	assert.Equal(t, out, l.allData(data))
}

func TestAllDataOverwriteDefaults(t *testing.T) {
	var (
		defaults = map[string]string{"one": "1"}
		data     = map[string]string{"one": "111"}
		l        = &logger{defaults: defaults}
	)

	assert.Equal(t, data, l.allData(data))
}
