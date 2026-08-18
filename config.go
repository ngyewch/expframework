package expframework

import (
	"github.com/ngyewch/expframework/template"
)

type ExperimentConfig struct {
	BaseDirectory           string
	TemplateHelper          *template.Helper
	SessionIdTemplateString string
	TrialIdTemplateString   string
}
