package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ngyewch/expframework/camera/ipcamera"
	"github.com/urfave/cli/v3"
)

func doTestCameraMJPEG(ctx context.Context, cmd *cli.Command) error {
	cfg, err := getTestCameraConfig(ctx, cmd)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	var mjpegServers []*ipcamera.MJPEGServer
	defer func() {
		for _, mjpegServer := range mjpegServers {
			_ = mjpegServer.Close()
		}
	}()
	for _, cameraCfg := range cfg.Cameras {
		if cameraCfg.Disabled {
			continue
		}
		mjpegServer, err := ipcamera.NewMJPEGServer(cameraCfg)
		if err != nil {
			return err
		}
		mjpegServers = append(mjpegServers, mjpegServer)
		err = mjpegServer.Start(ctx)
		if err != nil {
			return err
		}
		pathPrefix := fmt.Sprintf("/mjpeg/%s", cameraCfg.Id)
		mux.Handle(fmt.Sprintf("%s/", pathPrefix), http.StripPrefix(pathPrefix, mjpegServer.ServeMux()))
	}

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		return err
	}

	return nil
}
