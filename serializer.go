package logged

type Serializer interface {
	Write(e *Entry) error
}
