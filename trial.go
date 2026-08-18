package expframework

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

type TrialContext[SessionProperties Properties, TrialProperties Properties] struct {
	parent        *SessionContext[SessionProperties, TrialProperties]
	info          TrialInfo[TrialProperties]
	baseDirectory string
	logger        *slog.Logger
	closers       []io.Closer
}

type TrialInfo[TrialProperties Properties] struct {
	TrialNumber int             `json:"trialNumber,omitempty"`
	Id          string          `json:"id,omitempty"`
	StartTime   time.Time       `json:"startTime,omitzero"`
	EndTime     time.Time       `json:"endTime,omitzero"`
	Properties  TrialProperties `json:"properties,omitempty"`
}

func newTrialContext[SessionProperties Properties, TrialProperties Properties](parent *SessionContext[SessionProperties, TrialProperties], trialNumber int, id string, baseDirectory string, properties TrialProperties, cfg ExperimentConfig) (*TrialContext[SessionProperties, TrialProperties], error) {
	logHandler, err := NewFileLogHandler(filepath.Join(baseDirectory, "logs.txt"))
	if err != nil {
		return nil, err
	}
	logger := slog.New(slog.NewMultiHandler(logHandler, parent.LogHandler(), slog.Default().Handler())).With(
		"sessionId", parent.Info().Id,
		"trialNumber", trialNumber,
		"trialId", id,
	)
	logger.Info("trial started", properties.LogArgs()...)
	var closers []io.Closer

	context := &TrialContext[SessionProperties, TrialProperties]{
		parent: parent,
		info: TrialInfo[TrialProperties]{
			TrialNumber: trialNumber,
			Id:          id,
			StartTime:   time.Now(),
			Properties:  properties,
		},
		baseDirectory: baseDirectory,
		logger:        logger,
		closers:       closers,
	}

	err = context.writeInfo()
	if err != nil {
		return nil, err
	}

	return context, nil
}

func (context *TrialContext[SessionProperties, TrialProperties]) Parent() *SessionContext[SessionProperties, TrialProperties] {
	return context.parent
}

func (context *TrialContext[SessionProperties, TrialProperties]) Info() TrialInfo[TrialProperties] {
	return context.info
}

func (context *TrialContext[SessionProperties, TrialProperties]) BaseDirectory() string {
	return context.baseDirectory
}

func (context *TrialContext[SessionProperties, TrialProperties]) Logger() *slog.Logger {
	return context.logger
}

func (context *TrialContext[SessionProperties, TrialProperties]) AddCloser(closer io.Closer) {
	context.closers = append(context.closers, closer)
}

func (context *TrialContext[SessionProperties, TrialProperties]) Close() error {
	for _, closer := range context.closers {
		err := closer.Close()
		if err != nil {
			context.logger.Error("error closing resource",
				slog.String("type", fmt.Sprintf("%T", closer)),
				slog.Any("err", err),
			)
		}
	}
	context.closers = nil
	if context.info.EndTime.IsZero() {
		context.info.EndTime = time.Now()
		context.logger.Info("trial ended")
	}

	err := context.writeInfo()
	if err != nil {
		context.logger.Error("error writing trial info",
			slog.Any("err", err),
		)
	}

	return nil
}

func (context *TrialContext[SessionProperties, TrialProperties]) writeInfo() error {
	outputFile := filepath.Join(context.baseDirectory, "info.json")
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	jsonEncoder := json.NewEncoder(f)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.SetIndent("", "  ")
	err = jsonEncoder.Encode(context.Info())
	if err != nil {
		return err
	}

	return nil
}
