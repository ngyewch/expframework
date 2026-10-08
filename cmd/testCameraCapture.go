package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ngyewch/expframework/camera/ipcamera"
	"github.com/ngyewch/expframework/config"
	"github.com/urfave/cli/v3"
	"golang.org/x/sync/errgroup"
)

type TestCameraCaptureConfig struct {
	Cameras []ipcamera.Config `json:"cameras" validate:"dive"`
}

func (config TestCameraCaptureConfig) Camera(id string) (*ipcamera.Config, error) {
	for _, cfg := range config.Cameras {
		if cfg.Id == id {
			return &cfg, nil
		}
	}
	return nil, fmt.Errorf("camera not found: %s", id)
}

func doTestCameraCapture(ctx context.Context, cmd *cli.Command) error {
	configFile := cmd.String(configFileFlag.Name)
	if configFile == "" {
		return fmt.Errorf("no config file specified")
	}

	cameraId := cmd.Args().Get(0)
	streamId := cmd.Args().Get(1)

	if cameraId == "" {
		return fmt.Errorf("no camera id specified")
	}
	if streamId == "" {
		return fmt.Errorf("no stream id specified")
	}

	cfg, err := config.LoadYAMLConfig[TestCameraCaptureConfig, *TestCameraCaptureConfig](configFile)
	if err != nil {
		return err
	}

	cameraCfg, err := cfg.Camera(cameraId)
	if err != nil {
		return err
	}

	cameraInstance := ipcamera.New(*cameraCfg)
	recording, err := cameraInstance.Capture(streamId, "workspace/output.mkv")
	if err != nil {
		return err
	}
	defer func(recording ipcamera.Recording) {
		_ = recording.Close()
	}(recording)

	errGroup, errGroupCtx := errgroup.WithContext(ctx)
	errGroup.Go(func() error {
		go func() {
			time.Sleep(10 * time.Second)
			err := recording.Close()
			if err != nil {
				slog.Error("error stopping video recording",
					slog.Any("err", err),
				)
			}
		}()

		err := recording.Start(errGroupCtx)
		if err != nil {
			return err
		}

		time.Sleep(10 * time.Second)

		return nil
	})

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-errGroupCtx.Done():
		return errGroupCtx.Err()
	case <-time.After(30 * time.Second):
	}

	return nil
}
