package ipcamera

import (
	"fmt"
	"log/slog"
	"slices"
)

type CaptureOption func(*captureOptions)

func WithLogger(logger *slog.Logger) CaptureOption {
	return func(options *captureOptions) {
		options.logger = logger
	}
}

type captureOptions struct {
	logger *slog.Logger
}

func (options captureOptions) Logger() *slog.Logger {
	if options.logger == nil {
		return slog.Default()
	}
	return options.logger
}

type IPCamera struct {
	cfg Config
}

func New(cfg Config) *IPCamera {
	return &IPCamera{
		cfg: cfg,
	}
}

func (camera *IPCamera) HasStream(streamId string) bool {
	if camera.cfg.Disabled {
		return false
	}
	streamConfig := camera.cfg.Stream(streamId)
	if streamConfig == nil {
		return false
	}
	if streamConfig.Disabled {
		return false
	}
	return true
}

func (camera *IPCamera) StreamIdsWithTag(tag string) []string {
	if camera.cfg.Disabled {
		return nil
	}
	var streamIds []string
	for _, streamConfig := range camera.cfg.Streams {
		if streamConfig.Disabled {
			continue
		}
		if slices.Contains(streamConfig.Tags, tag) {
			streamIds = append(streamIds, streamConfig.Id)
		}
	}
	return streamIds
}

func (camera *IPCamera) Capture(streamId string, outputPath string, options ...CaptureOption) (Recording, error) {
	if camera.cfg.Disabled {
		return nil, fmt.Errorf("camera disabled")
	}
	streamConfig := camera.cfg.Stream(streamId)
	if streamConfig == nil {
		return nil, fmt.Errorf("stream not found: %s", streamId)
	}
	if streamConfig.Disabled {
		return nil, fmt.Errorf("camera stream disabled")
	}
	captureOptions := &captureOptions{}
	for _, option := range options {
		option(captureOptions)
	}
	return newFFMPEGRecording(camera.cfg, *streamConfig, outputPath, captureOptions), nil
}
