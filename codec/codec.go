package codec

type Codec[T any, PT *T] interface {
	MarshalBinary(v PT) ([]byte, error)
	UnmarshalBinary(data []byte) (PT, error)
}

type FileCodec[T any, PT *T] interface {
	MarshalFile(path string, v PT) error
	UnmarshalFile(path string) (PT, error)
}
