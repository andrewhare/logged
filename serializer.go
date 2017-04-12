package logged

import (
	"bufio"
	"io"
	"sync"
	"unicode/utf8"
)

var hex = "0123456789abcdef"

type entry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Data      Data   `json:"data"`
}

func newSerializer(w io.Writer) *serializer {
	return &serializer{
		Writer: bufio.NewWriter(w),
	}
}

type serializer struct {
	*bufio.Writer
	mu sync.Mutex
}

func (s *serializer) write(e *entry) error {
	s.mu.Lock()

	s.WriteString(`{"timestamp":"`)
	s.WriteString(e.Timestamp)
	s.WriteString(`","level":"`)
	s.WriteString(e.Level)
	s.WriteString(`","message":`)
	s.writeJSONString(e.Message)

	if len(e.Data) > 0 {
		s.WriteString(`,"data":{`)
		first := true
		for k, v := range e.Data {
			if !first {
				s.WriteRune(',')
			}
			first = false
			s.writeJSONString(k)
			s.WriteRune(':')
			s.writeJSONString(v)
		}
		s.WriteRune('}')
	}

	s.WriteRune('}')
	s.WriteRune('\n')

	err := s.Flush()

	s.mu.Unlock()

	return err
}

func (s *serializer) writeJSONString(str string) {
	s.WriteByte('"')
	start := 0
	for i := 0; i < len(str); {
		if b := str[i]; b < utf8.RuneSelf {
			if safeSet[b] {
				i++
				continue
			}
			if start < i {
				s.WriteString(str[start:i])
			}
			switch b {
			case '\\', '"':
				s.WriteByte('\\')
				s.WriteByte(b)
			case '\n':
				s.WriteByte('\\')
				s.WriteByte('n')
			case '\r':
				s.WriteByte('\\')
				s.WriteByte('r')
			case '\t':
				s.WriteByte('\\')
				s.WriteByte('t')
			default:
				s.WriteString(`\u00`)
				s.WriteByte(hex[b>>4])
				s.WriteByte(hex[b&0xf])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(str[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				s.WriteString(str[start:i])
			}
			s.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(str) {
		s.WriteString(str[start:])
	}
	s.WriteByte('"')
}
