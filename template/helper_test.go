package template

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CurrentContext struct {
	SessionInfo *SessionInfo
	TrialInfo   *TrialInfo
}

type SessionInfo struct {
	StartTime *timestamppb.Timestamp
}

type TrialInfo struct {
	TrialNumber int32
	StartTime   *timestamppb.Timestamp
}

func TestTemplateHelper(t *testing.T) {
	timeLocation, err := time.LoadLocation("Asia/Singapore")
	if assert.NoError(t, err) {
		templateHelper, err := NewHelper(timeLocation, "%Y-%m-%d_%H.%M.%S")
		if assert.NoError(t, err) {
			{
				data := map[string]any{
					"Date": time.Date(2026, 3, 24, 17, 45, 05, 0, timeLocation),
				}
				s, err := templateHelper.ExecuteTemplate("{{ .Date | standardDate }}", data)
				if assert.NoError(t, err) {
					assert.Equal(t, "2026-03-24_17.45.05", s)
				}
			}

			{
				data := map[string]any{
					"Context": &CurrentContext{
						SessionInfo: &SessionInfo{
							StartTime: timestamppb.New(time.Date(2026, 3, 24, 17, 45, 05, 0, timeLocation)),
						},
						TrialInfo: &TrialInfo{
							TrialNumber: 2,
							StartTime:   timestamppb.New(time.Date(2026, 3, 24, 18, 30, 10, 0, timeLocation)),
						},
					},
				}
				s, err := templateHelper.ExecuteTemplate("{{ if .Context.SessionInfo }}{{ .Context.SessionInfo.StartTime | standardDate }}{{ else }}???{{ end }}", data)
				if assert.NoError(t, err) {
					assert.Equal(t, "2026-03-24_17.45.05", s)
				}
				s, err = templateHelper.ExecuteTemplate("{{ if .Context.TrialInfo }}{{ .Context.TrialInfo.StartTime | standardDate }}_T{{ .Context.TrialInfo.TrialNumber | sprintf \"%03d\" }}{{ else }}???{{ end }}", data)
				if assert.NoError(t, err) {
					assert.Equal(t, "2026-03-24_18.30.10_T002", s)
				}
			}
		}
	}
}
