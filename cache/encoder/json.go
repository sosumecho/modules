package encoder

import jsoniter "github.com/json-iterator/go"

type JSONEncoder[T any] struct{}

func (J JSONEncoder[T]) Encode(t any) ([]byte, error) {
	rs, err := jsoniter.Marshal(t)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

func (J JSONEncoder[T]) Decode(rs []byte) (T, error) {
	var item T
	if err := jsoniter.Unmarshal(rs, &item); err != nil {
		return item, err
	}
	return item, nil
}

func NewJSONEncoder[T any]() *JSONEncoder[T] {
	return &JSONEncoder[T]{}
}
