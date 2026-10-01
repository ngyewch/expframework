package codec

import (
	"google.golang.org/protobuf/proto"
)

type ProtoCodec[T any, PT interface {
	*T
	proto.Message
}] struct {
}

func (codec ProtoCodec[T, PT]) MarshalBinary(v PT) ([]byte, error) {
	msgPtr := PT(v)
	return proto.Marshal(msgPtr)
}

func (codec ProtoCodec[T, PT]) UnmarshalBinary(data []byte) (PT, error) {
	var msg T
	msgPtr := PT(&msg)
	err := proto.Unmarshal(data, msgPtr)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}
