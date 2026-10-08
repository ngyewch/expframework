package main

import (
	"context"
	"fmt"

	"github.com/ngyewch/expframework/camera/ipcamera"
	"github.com/ngyewch/expframework/config"
	"github.com/urfave/cli/v3"
)

type TestCameraConfig struct {
	Cameras []ipcamera.Config `json:"cameras" validate:"dive"`
}

func (config TestCameraConfig) Camera(id string) (*ipcamera.Config, error) {
	for _, cfg := range config.Cameras {
		if cfg.Id == id {
			return &cfg, nil
		}
	}
	return nil, fmt.Errorf("camera not found: %s", id)
}

func getTestCameraConfig(ctx context.Context, cmd *cli.Command) (*TestCameraConfig, error) {
	configFile := cmd.String(configFileFlag.Name)
	if configFile == "" {
		return nil, fmt.Errorf("no config file specified")
	}

	cfg, err := config.LoadYAMLConfig[TestCameraConfig, *TestCameraConfig](configFile)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
