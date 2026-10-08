package ipcamera

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ngyewch/go-mjpeg"
)

type MJPEGStream struct {
	cfg           Config
	streamCfg     StreamConfig
	streamUrl     string
	streamDecoder *mjpeg.RetryingURLStreamDecoder
	streamProxy   *mjpeg.Proxy
	started       bool
}

func NewMJPEGStream(cfg Config, streamCfg StreamConfig, httpClient *http.Client) (*MJPEGStream, error) {
	mjpegStreamUrl, err := streamCfg.ToURL(cfg, false)
	if err != nil {
		return nil, err
	}
	credentials := streamCfg.ResolveCredentials(cfg)
	var disconnectedFrame []byte
	if (streamCfg.Dimensions != nil) && (streamCfg.Dimensions.Width > 0) && (streamCfg.Dimensions.Height > 0) {
		var err error
		disconnectedFrame, err = getDisconnectedFrame(nil, streamCfg.Dimensions.Width, streamCfg.Dimensions.Height)
		if err != nil {
			return nil, err
		}
	}
	mjpegStreamDecoder := mjpeg.NewRetryingURLStreamDecoder(
		mjpegStreamUrl.String(),
		httpClient,
		func(httpRequest *http.Request) {
			if credentials != nil {
				httpRequest.SetBasicAuth(credentials.Username, credentials.Password)
			}
		},
		streamCfg.Timeout,
		disconnectedFrame,
	)
	mjpegStreamProxy := mjpeg.NewProxy(mjpegStreamDecoder)
	return &MJPEGStream{
		cfg:           cfg,
		streamCfg:     streamCfg,
		streamUrl:     mjpegStreamUrl.String(),
		streamDecoder: mjpegStreamDecoder,
		streamProxy:   mjpegStreamProxy,
	}, nil
}

func (mjpegStream *MJPEGStream) Close() error {
	_ = mjpegStream.streamDecoder.Close()
	return nil
}

func (mjpegStream *MJPEGStream) Start(ctx context.Context) error {
	if mjpegStream.started {
		return fmt.Errorf("mjpeg stream already started")
	}
	mjpegStream.started = true
	slog.Info("mjpeg stream started",
		slog.String("cameraId", mjpegStream.cfg.Id),
		slog.String("streamId", mjpegStream.streamCfg.Id),
		slog.String("streamUrl", mjpegStream.streamUrl),
	)
	defer func() {
		slog.Info("mjpeg stream stopped",
			slog.String("cameraId", mjpegStream.cfg.Id),
			slog.String("streamId", mjpegStream.streamCfg.Id),
			slog.String("streamUrl", mjpegStream.streamUrl),
		)
	}()
	err := mjpegStream.streamProxy.Run(ctx)
	if err != nil {
		slog.Error("mjpeg stream error",
			slog.Any("err", err),
		)
		return err
	}
	return nil
}

func (mjpegStream *MJPEGStream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mjpegStream.streamProxy.ServeHTTP(w, r)
}
