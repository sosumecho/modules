package encoder

type Encoder[T any] interface {
	Encode(t any) ([]byte, error)
	Decode(rs []byte) (T, error)
}

func New[T any](typ string) Encoder[T] {
	switch typ {
	case "json":
		return NewJSONEncoder[T]()
	case "gob":
		return NewGobEncoder[T]()
	}
	return nil
}
