package logged

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJSONSerializer(t *testing.T) {
	var (
		buf bytes.Buffer
		s   = newJSONSerializer(&buf)
		e   = &entry{Timestamp: "rightnow", Level: "somelevel", Message: "test123", Data: Data{"test": "123", "test2": "345"}}
		out entry
	)

	assert.NoError(t, s.write(e))
	assert.NoError(t, json.NewDecoder(&buf).Decode(&out))
	assert.Equal(t, e, &out)
}

func TestJSONSerializerNoData(t *testing.T) {
	var (
		buf bytes.Buffer
		s   = newJSONSerializer(&buf)
		e   = &entry{Timestamp: "sometimelater", Level: "otherlevel", Message: "345test"}
		out entry
	)

	assert.NoError(t, s.write(e))
	assert.NoError(t, json.NewDecoder(&buf).Decode(&out))
	assert.Equal(t, e, &out)
}

func TestJSONSerializerBadChars(t *testing.T) {
	var (
		buf bytes.Buffer
		s   = newJSONSerializer(&buf)
		e   = &entry{Timestamp: "sometimelater", Level: "otherlevel", Message: "\" \\ \b \f \n \r \t"}
		out entry
	)

	assert.NoError(t, s.write(e))
	assert.NoError(t, json.NewDecoder(&buf).Decode(&out))
	assert.Equal(t, e, &out)
}

func TestJSONSerializerExtended(t *testing.T) {
	var (
		buf bytes.Buffer
		s   = newJSONSerializer(&buf)
		e   = &entry{Timestamp: "sometimelater", Level: "otherlevel", Message: "Hello, 世界 √ 😂 \ufffd"}
		out entry
	)

	assert.NoError(t, s.write(e))
	assert.NoError(t, json.NewDecoder(&buf).Decode(&out))
	assert.Equal(t, e, &out)
}

func BenchmarkJSONSerializer(b *testing.B) {
	var (
		s = newJSONSerializer(os.Stdout)
		e = &entry{
			Level:     "debug",
			Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
			Message:   "this is a test of the serializer for a message",
		}
	)

	for n := 0; n < b.N; n++ {
		s.write(e)
	}
}
