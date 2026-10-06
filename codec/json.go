package codec

import (
	"encoding/json"
	"os"
)

type JsonCodec[T any, PT *T] struct {
}

func (codec JsonCodec[T, PT]) MarshalBinary(v PT) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

func (codec JsonCodec[T, PT]) UnmarshalBinary(data []byte) (PT, error) {
	var msg T
	msgPtr := PT(&msg)
	err := json.Unmarshal(data, msgPtr)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (codec JsonCodec[T, PT]) MarshalFile(path string, v PT) error {
	b, err := codec.MarshalBinary(v)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, b, 0755)
	if err != nil {
		return err
	}
	return nil
}

func (codec JsonCodec[T, PT]) UnmarshalFile(path string) (PT, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	v, err := codec.UnmarshalBinary(b)
	if err != nil {
		return nil, err
	}
	return v, nil
}
