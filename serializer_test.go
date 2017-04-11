package logged

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSerializer(t *testing.T) {
	var (
		buf bytes.Buffer
		s   = newSerializer(&buf)
		e   = &entry{Timestamp: "rightnow", Level: "somelevel", Message: "test123", Data: Data{"test": "123", "test2": "345"}}
		out entry
	)

	assert.NoError(t, s.write(e))
	assert.NoError(t, json.NewDecoder(&buf).Decode(&out))
	assert.Equal(t, e, &out)
}

func TestSerializerNoData(t *testing.T) {
	var (
		buf bytes.Buffer
		s   = newSerializer(&buf)
		e   = &entry{Timestamp: "sometimelater", Level: "otherlevel", Message: "345test"}
		out entry
	)

	assert.NoError(t, s.write(e))
	assert.NoError(t, json.NewDecoder(&buf).Decode(&out))
	assert.Equal(t, e, &out)
}

func TestSerializerBadChars(t *testing.T) {
	var (
		buf bytes.Buffer
		s   = newSerializer(&buf)
		e   = &entry{Timestamp: "sometimelater", Level: "otherlevel", Message: "\" \\ \b \f \n \r \t"}
		out entry
	)

	assert.NoError(t, s.write(e))
	assert.NoError(t, json.NewDecoder(&buf).Decode(&out))
	assert.Equal(t, e, &out)
}

func TestSerializerExtended(t *testing.T) {
	var (
		buf bytes.Buffer
		s   = newSerializer(&buf)
		e   = &entry{Timestamp: "sometimelater", Level: "otherlevel", Message: "Hello, 世界 √ 😂 \ufffd"}
		out entry
	)

	assert.NoError(t, s.write(e))
	assert.NoError(t, json.NewDecoder(&buf).Decode(&out))
	assert.Equal(t, e, &out)
}

func BenchmarkSerializer(b *testing.B) {
	var (
		s = newSerializer(os.Stdout)
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
