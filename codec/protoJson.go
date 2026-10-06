package codec

import (
	"os"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type ProtoJsonCodec[T any, PT interface {
	*T
	proto.Message
}] struct {
}

func (codec ProtoJsonCodec[T, PT]) MarshalBinary(v PT) ([]byte, error) {
	marshalOptions := protojson.MarshalOptions{
		Multiline: true,
		Indent:    " ",
	}
	return marshalOptions.Marshal(v)
}

func (codec ProtoJsonCodec[T, PT]) UnmarshalBinary(data []byte) (PT, error) {
	var msg T
	msgPtr := PT(&msg)
	err := protojson.Unmarshal(data, msgPtr)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (codec ProtoJsonCodec[T, PT]) MarshalFile(path string, v PT) error {
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

func (codec ProtoJsonCodec[T, PT]) UnmarshalFile(path string) (PT, error) {
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
