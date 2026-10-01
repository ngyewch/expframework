package codec

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type ProtoJsonCodec[T any, PT interface {
	*T
	proto.Message
}] struct {
}

func (codec ProtoJsonCodec[T, PT]) MarshalBinary(v *T) ([]byte, error) {
	marshalOptions := protojson.MarshalOptions{
		Multiline: true,
		Indent:    " ",
	}
	msgPtr := PT(v)
	return marshalOptions.Marshal(msgPtr)
}

func (codec ProtoJsonCodec[T, PT]) UnmarshalBinary(data []byte) (*T, error) {
	var msg T
	msgPtr := PT(&msg)
	err := protojson.Unmarshal(data, msgPtr)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}
