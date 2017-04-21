package logged

// Entry represents a single line in the log.
type Entry struct {
	Timestamp string            `json:"timestamp"`
	Level     string            `json:"level"`
	Message   string            `json:"message"`
	Data      map[string]string `json:"data"`
}
