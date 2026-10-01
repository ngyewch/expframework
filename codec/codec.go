package codec

type Codec[T any, PT *T] interface {
	MarshalBinary(v PT) ([]byte, error)
	UnmarshalBinary(data []byte) (PT, error)
}
