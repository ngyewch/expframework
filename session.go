package expframework

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/jacoblockett/gosan/v3"
)

type SessionContext[SessionProperties Properties, TrialProperties Properties] struct {
	cfg           ExperimentConfig
	info          SessionInfo[SessionProperties]
	baseDirectory string
	logHandler    slog.Handler
	logger        *slog.Logger
	trialContext  *TrialContext[SessionProperties, TrialProperties]
}

type SessionInfo[SessionProperties Properties] struct {
	Id         string            `json:"id,omitempty"`
	StartTime  time.Time         `json:"startTime,omitzero"`
	EndTime    time.Time         `json:"endTime,omitzero"`
	Properties SessionProperties `json:"properties,omitempty"`
}

func newSessionContext[SessionProperties Properties, TrialProperties Properties](id string, baseDirectory string, properties SessionProperties, cfg ExperimentConfig) (*SessionContext[SessionProperties, TrialProperties], error) {
	logHandler, err := NewFileLogHandler(filepath.Join(baseDirectory, "logs.txt"))
	if err != nil {
		return nil, err
	}
	logger := slog.New(slog.NewMultiHandler(logHandler, slog.Default().Handler())).With(
		"sessionId", id,
	)
	logger.Info("session started", properties.LogArgs()...)

	context := &SessionContext[SessionProperties, TrialProperties]{
		cfg:           cfg,
		baseDirectory: baseDirectory,
		logHandler:    logHandler,
		logger:        logger,
		info: SessionInfo[SessionProperties]{
			Id:         id,
			StartTime:  time.Now(),
			Properties: properties,
		},
	}

	err = context.writeInfo()
	if err != nil {
		return nil, err
	}

	return context, nil
}

func (context *SessionContext[SessionProperties, TrialProperties]) Info() SessionInfo[SessionProperties] {
	return context.info
}

func (context *SessionContext[SessionProperties, TrialProperties]) BaseDirectory() string {
	return context.baseDirectory
}

func (context *SessionContext[SessionProperties, TrialProperties]) LogHandler() slog.Handler {
	return context.logHandler
}

func (context *SessionContext[SessionProperties, TrialProperties]) Logger() *slog.Logger {
	return context.logger
}

func (context *SessionContext[SessionProperties, TrialProperties]) Close() error {
	if context.trialContext != nil {
		err := context.trialContext.Close()
		if err != nil {
			context.logger.Error("error closing trial context",
				slog.Any("err", err),
			)
		}
		context.trialContext = nil
	}
	if context.info.EndTime.IsZero() {
		context.info.EndTime = time.Now()
		context.logger.Info("session ended")
	}

	err := context.writeInfo()
	if err != nil {
		context.logger.Error("error writing session info",
			slog.Any("err", err),
		)
	}

	return nil
}

func (context *SessionContext[SessionProperties, TrialProperties]) CurrentTrial() *TrialContext[SessionProperties, TrialProperties] {
	return context.trialContext
}

func (context *SessionContext[SessionProperties, TrialProperties]) StartTrial(trialNumber int, properties TrialProperties) (*TrialContext[SessionProperties, TrialProperties], error) {
	if context.trialContext != nil {
		_ = context.trialContext.Close()
		context.trialContext = nil
	}

	trialInfo := &TrialInfo[TrialProperties]{
		TrialNumber: trialNumber,
		StartTime:   time.Now(),
		Properties:  properties,
	}
	trialId, err := context.cfg.TemplateHelper.ExecuteTemplate(context.cfg.TrialIdTemplateString, context.getWrappedTemplateData(trialInfo))
	if err != nil {
		return nil, err
	}
	trialInfo.Id = trialId

	dirName, err := gosan.Filename(trialId, nil)
	if err != nil {
		return nil, err
	}
	baseDirectory := filepath.Join(context.baseDirectory, dirName)
	_, err = os.Stat(baseDirectory)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("trial directory '%s' already exists", baseDirectory)
	}
	err = os.MkdirAll(baseDirectory, 0755)
	if err != nil {
		return nil, err
	}
	trialContext, err := newTrialContext[SessionProperties, TrialProperties](context, trialNumber, trialId, baseDirectory, properties, context.cfg)
	if err != nil {
		return nil, err
	}
	context.trialContext = trialContext
	return trialContext, nil
}

func (context *SessionContext[SessionProperties, TrialProperties]) getTemplateData(trialInfo *TrialInfo[TrialProperties]) TemplateData[SessionProperties, TrialProperties] {
	sessionInfo := context.Info()
	return TemplateData[SessionProperties, TrialProperties]{
		SessionInfo: &sessionInfo,
		TrialInfo:   trialInfo,
	}
}

func (context *SessionContext[SessionProperties, TrialProperties]) getWrappedTemplateData(trialInfo *TrialInfo[TrialProperties]) any {
	return map[string]any{
		"Context": context.getTemplateData(trialInfo),
	}
}

func (context *SessionContext[SessionProperties, TrialProperties]) writeInfo() error {
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
