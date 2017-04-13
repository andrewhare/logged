package logged

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTextSerializer(t *testing.T) {
	var (
		s = NewTextSerializer(os.Stdout)
		e = &Entry{Timestamp: "rightnow", Level: "somelevel", Message: "test123", Data: map[string]string{"test": "123", "test2": "345"}}
	)

	assert.NoError(t, s.Write(e))
}

func TestTextSerializerNoData(t *testing.T) {
	var (
		s = NewTextSerializer(os.Stdout)
		e = &Entry{Timestamp: "sometimelater", Level: "otherlevel", Message: "345test"}
	)

	assert.NoError(t, s.Write(e))
}

func BenchmarkTextSerializer(b *testing.B) {
	var (
		s = NewTextSerializer(os.Stdout)
		e = &Entry{
			Level:     "debug",
			Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
			Message:   "this is a test of the serializer for a message",
		}
	)

	for n := 0; n < b.N; n++ {
		s.Write(e)
	}
}
