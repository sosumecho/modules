package encoder

import (
	"bytes"
	"encoding/gob"
)

type GobEncoder[T any] struct{}

func (g GobEncoder[T]) Encode(t any) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(t)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (g GobEncoder[T]) Decode(rs []byte) (T, error) {
	var item T
	decoder := gob.NewDecoder(bytes.NewBuffer([]byte(rs)))
	err := decoder.Decode(&item)
	if err != nil {
		return item, err
	}
	return item, nil
}

func NewGobEncoder[T any]() *GobEncoder[T] {
	return &GobEncoder[T]{}
}
