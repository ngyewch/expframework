package ipcamera

import (
	"context"
	"fmt"
	"net/http"
)

type MJPEGServer struct {
	cfg     Config
	streams []*MJPEGStream
	mux     *http.ServeMux
	started bool
}

func NewMJPEGServer(cfg Config) (*MJPEGServer, error) {
	mux := http.NewServeMux()
	var mjpegStreams []*MJPEGStream
	for _, streamCfg := range cfg.Streams {
		if streamCfg.Disabled || (streamCfg.Type != StreamTypeMJPEGOverHTTP) {
			continue
		}
		mjpegStream, err := NewMJPEGStream(cfg, streamCfg, nil)
		if err != nil {
			return nil, err
		}
		mjpegStreams = append(mjpegStreams, mjpegStream)
		mux.Handle(fmt.Sprintf("/%s", streamCfg.Id), mjpegStream)
	}
	return &MJPEGServer{
		cfg:     cfg,
		streams: mjpegStreams,
		mux:     mux,
	}, nil
}

func (mjpegServer *MJPEGServer) Close() error {
	for _, stream := range mjpegServer.streams {
		_ = stream.Close()
	}
	return nil
}

func (mjpegServer *MJPEGServer) ServeMux() *http.ServeMux {
	return mjpegServer.mux
}

func (mjpegServer *MJPEGServer) Start(ctx context.Context) error {
	if mjpegServer.started {
		return fmt.Errorf("mjpeg server already started")
	}
	for _, stream := range mjpegServer.streams {
		go func() {
			_ = stream.Start(ctx)
		}()
	}
	mjpegServer.started = true
	return nil
}
