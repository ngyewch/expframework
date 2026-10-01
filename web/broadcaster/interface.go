package broadcaster

import "context"

type Broadcaster[T any, PT *T] interface {
	Publish(ctx context.Context, m PT) error
}
