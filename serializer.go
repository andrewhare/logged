package logged

// Serializer provides a function for writing an entry
type Serializer interface {
	Write(e *Entry) error
}
