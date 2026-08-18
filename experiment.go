package expframework

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/jacoblockett/gosan/v3"
)

type Properties interface {
	LogAttrs() []slog.Attr
	LogArgs() []any
}

type ExperimentContext[SessionProperties Properties, TrialProperties Properties] struct {
	cfg            ExperimentConfig
	sessionContext *SessionContext[SessionProperties, TrialProperties]
}

func NewExperimentContext[SessionProperties Properties, TrialProperties Properties](cfg ExperimentConfig) *ExperimentContext[SessionProperties, TrialProperties] {
	return &ExperimentContext[SessionProperties, TrialProperties]{
		cfg: cfg,
	}
}

func (context *ExperimentContext[SessionProperties, TrialProperties]) CurrentSession() *SessionContext[SessionProperties, TrialProperties] {
	return context.sessionContext
}

func (context *ExperimentContext[SessionProperties, TrialProperties]) StartSession(properties SessionProperties) (*SessionContext[SessionProperties, TrialProperties], error) {
	if context.sessionContext != nil {
		_ = context.sessionContext.Close()
		context.sessionContext = nil
	}

	sessionInfo := &SessionInfo[SessionProperties]{
		StartTime:  time.Now(),
		Properties: properties,
	}
	sessionId, err := context.cfg.TemplateHelper.ExecuteTemplate(context.cfg.SessionIdTemplateString, context.getWrappedTemplateData(sessionInfo))
	if err != nil {
		return nil, err
	}
	sessionInfo.Id = sessionId

	dirName, err := gosan.Filename(sessionId, nil)
	if err != nil {
		return nil, err
	}
	baseDirectory := filepath.Join(context.cfg.BaseDirectory, dirName)
	_, err = os.Stat(baseDirectory)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("session directory '%s' already exists", baseDirectory)
	}
	err = os.MkdirAll(baseDirectory, 0755)
	if err != nil {
		return nil, err
	}
	sessionContext, err := newSessionContext[SessionProperties, TrialProperties](sessionId, baseDirectory, properties, context.cfg)
	if err != nil {
		return nil, err
	}
	context.sessionContext = sessionContext
	return sessionContext, nil
}

func (context *ExperimentContext[SessionProperties, TrialProperties]) getTemplateData(sessionInfo *SessionInfo[SessionProperties]) TemplateData[SessionProperties, TrialProperties] {
	return TemplateData[SessionProperties, TrialProperties]{
		SessionInfo: sessionInfo,
	}
}

func (context *ExperimentContext[SessionProperties, TrialProperties]) getWrappedTemplateData(sessionInfo *SessionInfo[SessionProperties]) any {
	return map[string]any{
		"Context": context.getTemplateData(sessionInfo),
	}
}
