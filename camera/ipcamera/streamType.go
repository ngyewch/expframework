package ipcamera

import (
	"fmt"
	"strings"
)

type StreamType string

const (
	StreamTypeRTSPOverUDP          StreamType = "rtsp/udp"
	StreamTypeRTSPOverTCP          StreamType = "rtsp/tcp"
	StreamTypeRTSPOverUDPMulticast StreamType = "rtsp/udp-multicast"
	StreamTypeRTSPOverHTTP         StreamType = "rtsp/http"
	StreamTypeMJPEGOverHTTP        StreamType = "mjpeg/http"
)

func (streamType StreamType) String() string {
	return string(streamType)
}

func (streamType *StreamType) UnmarshalText(b []byte) error {
	s := strings.ToLower(string(b))
	switch s {
	case "rtsp/udp":
		*streamType = StreamTypeRTSPOverUDP
		return nil
	case "rtsp/tcp":
		*streamType = StreamTypeRTSPOverTCP
		return nil
	case "rtsp/udp-multicast":
		*streamType = StreamTypeRTSPOverUDPMulticast
		return nil
	case "rtsp/http":
		*streamType = StreamTypeRTSPOverHTTP
		return nil
	case "mjpeg/http":
		*streamType = StreamTypeMJPEGOverHTTP
		return nil
	default:
		return fmt.Errorf("unknown stream type: %s", string(b))
	}
}

func (streamType StreamType) MarshalText() ([]byte, error) {
	return []byte(streamType.String()), nil
}
