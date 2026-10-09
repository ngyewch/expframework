package ipcamera

import (
	"fmt"
	"net/netip"
	"net/url"
	"time"
)

type Config struct {
	Id          string                            `json:"id" validate:"required"`
	Disabled    bool                              `json:"disabled"`
	Address     netip.Addr                        `json:"address" validate:"required"`
	Credentials *UsernamePasswordCredentialConfig `json:"credentials"`
	Streams     []StreamConfig                    `json:"streams" validate:"dive"`
}

func (cfg Config) Stream(id string) *StreamConfig {
	for _, stream := range cfg.Streams {
		if stream.Id == id {
			return &stream
		}
	}
	return nil
}

type StreamConfig struct {
	Id          string                            `json:"id" validate:"required"`
	Disabled    bool                              `json:"disabled"`
	Tags        []string                          `json:"tags"`
	Credentials *UsernamePasswordCredentialConfig `json:"credentials"`
	Type        StreamType                        `json:"type" validate:"required"`
	Port        uint16                            `json:"port"`
	Path        string                            `json:"path"`
	Timeout     time.Duration                     `json:"timeout"`
	Dimensions  *Dimensions                       `json:"dimensions"`
}

func (streamConfig StreamConfig) ToURL(cfg Config, includeCredentials bool) (*url.URL, error) {
	var u url.URL
	u.Host = streamConfig.ResolveURLHost(cfg)
	u.Path = streamConfig.Path
	if includeCredentials {
		credentials := streamConfig.ResolveCredentials(cfg)
		if credentials != nil {
			u.User = url.UserPassword(credentials.Username, credentials.Password)
		}
	}
	switch streamConfig.Type {
	case StreamTypeRTSPOverUDP:
		u.Scheme = "rtsp"
	case StreamTypeRTSPOverTCP:
		u.Scheme = "rtsp"
	case StreamTypeRTSPOverUDPMulticast:
		u.Scheme = "rtsp"
	case StreamTypeRTSPOverHTTP:
		u.Scheme = "rtsp"
	case StreamTypeMJPEGOverHTTP:
		u.Scheme = "http"
	default:
		return nil, fmt.Errorf("unknown stream type: %s", streamConfig.Type)
	}
	return &u, nil
}

func (streamConfig StreamConfig) ResolveCredentials(cfg Config) *UsernamePasswordCredentialConfig {
	if streamConfig.Credentials != nil {
		return streamConfig.Credentials
	}
	return cfg.Credentials
}

func (streamConfig StreamConfig) ResolveURLHost(cfg Config) string {
	if streamConfig.Port > 0 {
		return netip.AddrPortFrom(cfg.Address, streamConfig.Port).String()
	}
	return cfg.Address.String()
}

type UsernamePasswordCredentialConfig struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Dimensions struct {
	Width  int `json:"width" validate:"required"`
	Height int `json:"height" validate:"required"`
}
