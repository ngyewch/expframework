package ipcamera

import (
	"context"
	"io"
)

type Recording interface {
	io.Closer

	Config() Config
	StreamConfig() StreamConfig
	OutputPath() string
	Start(ctx context.Context) error
}
