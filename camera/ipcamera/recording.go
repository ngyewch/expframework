package ipcamera

import (
	"context"
	"io"
)

type Recording interface {
	io.Closer

	OutputPath() string
	Start(ctx context.Context) error
}
