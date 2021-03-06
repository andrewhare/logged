package logged

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJSONSerializer(t *testing.T) {
	var (
		buf bytes.Buffer
		s   = NewJSONSerializer(&buf)
		e   = &Entry{Timestamp: "rightnow", Level: "somelevel", Message: "test123", Data: map[string]string{"test": "123", "test2": "345"}}
		out Entry
	)

	assert.NoError(t, s.Write(e))
	assert.NoError(t, json.NewDecoder(&buf).Decode(&out))
	assert.Equal(t, e, &out)
}

func TestJSONSerializerNoData(t *testing.T) {
	var (
		buf bytes.Buffer
		s   = NewJSONSerializer(&buf)
		e   = &Entry{Timestamp: "sometimelater", Level: "otherlevel", Message: "345test"}
		out Entry
	)

	assert.NoError(t, s.Write(e))
	assert.NoError(t, json.NewDecoder(&buf).Decode(&out))
	assert.Equal(t, e, &out)
}

func TestJSONSerializerBadChars(t *testing.T) {
	var (
		buf bytes.Buffer
		s   = NewJSONSerializer(&buf)
		e   = &Entry{Timestamp: "sometimelater", Level: "otherlevel", Message: "\" \\ \b \f \n \r \t"}
		out Entry
	)

	assert.NoError(t, s.Write(e))
	assert.NoError(t, json.NewDecoder(&buf).Decode(&out))
	assert.Equal(t, e, &out)
}

func TestJSONSerializerExtended(t *testing.T) {
	var (
		buf bytes.Buffer
		s   = NewJSONSerializer(&buf)
		e   = &Entry{Timestamp: "sometimelater", Level: "otherlevel", Message: "Hello, 世界 √ 😂 \ufffd"}
		out Entry
	)

	assert.NoError(t, s.Write(e))
	assert.NoError(t, json.NewDecoder(&buf).Decode(&out))
	assert.Equal(t, e, &out)
}

func BenchmarkJSONSerializer(b *testing.B) {
	var (
		buf bytes.Buffer
		s   = NewJSONSerializer(&buf)
		e   = &Entry{
			Level:     "debug",
			Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
			Message:   "this is a test of the serializer for a message",
		}
	)

	for n := 0; n < b.N; n++ {
		s.Write(e)
	}
}

type stdlibJSONSerializer struct{ e *json.Encoder }

func (s *stdlibJSONSerializer) Write(e *Entry) error { return s.e.Encode(e) }

func BenchmarkStdlibJSONSerializer(b *testing.B) {
	var (
		buf bytes.Buffer
		s   = &stdlibJSONSerializer{json.NewEncoder(&buf)}
		e   = &Entry{
			Level:     "debug",
			Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
			Message:   "this is a test of the serializer for a message",
		}
	)

	for n := 0; n < b.N; n++ {
		s.Write(e)
	}
}
