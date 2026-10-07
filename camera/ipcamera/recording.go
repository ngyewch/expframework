package ipcamera

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type Recording struct {
	config       Config
	streamConfig StreamConfig
	outputPath   string
	options      *captureOptions
	stdout       *os.File
	stderr       *os.File
	cancel       func()
	cmd          *exec.Cmd
	interrupted  bool
}

func newFFMPEGRecording(config Config, streamConfig StreamConfig, outputPath string, options *captureOptions) *Recording {
	return &Recording{
		config:       config,
		streamConfig: streamConfig,
		outputPath:   outputPath,
		options:      options,
	}
}

func (recording *Recording) Start(ctx context.Context) error {
	if recording.cmd != nil {
		return fmt.Errorf("recording already started")
	}

	args := []string{
		"-hide_banner",
	}
	switch recording.streamConfig.Type {
	case StreamTypeRTSPOverUDP:
		args = append(args, "-rtsp_transport", "udp")
	case StreamTypeRTSPOverTCP:
		args = append(args, "-rtsp_transport", "tcp")
	case StreamTypeRTSPOverUDPMulticast:
		args = append(args, "-rtsp_transport", "udp_multicast")
	case StreamTypeRTSPOverHTTP:
		args = append(args, "-rtsp_transport", "http")
	case StreamTypeMJPEGOverHTTP:
	default:
		return fmt.Errorf("unsupported stream type: %s", recording.streamConfig.Type)
	}
	if recording.streamConfig.Timeout > 0 {
		args = append(args, "-timeout", strconv.FormatInt(recording.streamConfig.Timeout.Microseconds(), 10))
	}
	inputUrl, err := recording.streamConfig.ToURL(recording.config, true)
	if err != nil {
		return err
	}
	args = append(args, "-fflags", "+genpts+igndts")
	args = append(args, "-i", inputUrl.String())
	args = append(args, "-codec", "copy")
	args = append(args, "-flush_packets", "1")
	args = append(args, recording.outputPath)

	stdout, err := os.Create(recording.outputPath + ".stdout.log")
	if err != nil {
		return err
	}
	stderr, err := os.Create(recording.outputPath + ".stderr.log")
	if err != nil {
		_ = stdout.Close()
		return err
	}

	recording.options.Logger().Info("video recording started",
		slog.String("path", recording.outputPath),
		slog.String("cmd", "ffmpeg"),
		slog.Any("args", args),
	)

	recordingCtx, cancel := context.WithCancel(ctx)
	recording.cancel = cancel

	cmd := exec.CommandContext(recordingCtx, "ffmpeg", args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Cancel = func() error {
		err := cmd.Process.Signal(os.Interrupt)
		if err != nil {
			return err
		}
		recording.interrupted = true
		return nil
	}
	cmd.WaitDelay = 2 * time.Second

	recording.cmd = cmd
	recording.stdout = stdout
	recording.stderr = stderr

	err = cmd.Run()
	if !recording.interrupted && (err != nil) {
		recording.options.Logger().Error("video recording error",
			slog.String("path", recording.outputPath),
			slog.Any("err", err),
		)
		return err
	}

	recording.options.Logger().Info("video recording stopped",
		slog.String("path", recording.outputPath),
	)
	return nil
}

func (recording *Recording) Close() error {
	if recording.cancel != nil {
		recording.cancel()
		recording.cancel = nil
	}
	if recording.stdout != nil {
		_ = recording.stdout.Close()
		recording.stdout = nil
	}
	if recording.stderr != nil {
		_ = recording.stderr.Close()
		recording.stderr = nil
	}
	return nil
}
