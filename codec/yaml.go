package codec

import (
	"os"

	"github.com/goccy/go-yaml"
)

type YamlCodec[T any, PT *T] struct {
}

func (codec YamlCodec[T, PT]) MarshalBinary(v PT) ([]byte, error) {
	return yaml.Marshal(v)
}

func (codec YamlCodec[T, PT]) UnmarshalBinary(data []byte) (PT, error) {
	var msg T
	msgPtr := PT(&msg)
	err := yaml.Unmarshal(data, msgPtr)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (codec YamlCodec[T, PT]) MarshalFile(path string, v PT) error {
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

func (codec YamlCodec[T, PT]) UnmarshalFile(path string) (PT, error) {
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
