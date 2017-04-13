package logged

import (
	"bufio"
	"io"
	"sync"
)

func NewTextSerializer(w io.Writer) Serializer {
	return &textSerializer{
		w: bufio.NewWriter(w),
	}
}

type textSerializer struct {
	w  *bufio.Writer
	mu sync.Mutex
}

func (s *textSerializer) Write(e *Entry) error {
	s.mu.Lock()

	s.w.WriteString(e.Timestamp)

	s.w.WriteRune(' ')
	s.w.WriteRune('[')
	s.w.WriteString(e.Level)
	s.w.WriteRune(']')
	s.w.WriteRune(' ')

	s.w.WriteString(e.Message)

	s.w.WriteRune(' ')

	for k, v := range e.Data {
		s.w.WriteString(k)
		s.w.WriteRune('=')
		s.w.WriteString(v)
		s.w.WriteRune(' ')
	}

	s.w.WriteRune('\n')

	err := s.w.Flush()

	s.mu.Unlock()

	return err
}
