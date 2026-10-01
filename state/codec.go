package state

type Codec[T any] interface {
	MarshalBinary(v *T) ([]byte, error)
	UnmarshalBinary(data []byte) (*T, error)
}
