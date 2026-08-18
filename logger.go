package expframework

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/phsym/console-slog"
)

func NewFileLogHandler(path string) (slog.Handler, error) {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return nil, err
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	logLevel := slog.LevelInfo
	_ = logLevel.UnmarshalText([]byte(os.Getenv("LOG_LEVEL")))
	addSource := os.Getenv("LOG_ADD_SOURCE") == "true"

	return console.NewHandler(f, &console.HandlerOptions{
		Level:     logLevel,
		AddSource: addSource,
		NoColor:   true,
	}), nil
}
