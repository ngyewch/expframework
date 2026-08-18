package expframework

type TemplateData[SessionProperties Properties, TrialProperties Properties] struct {
	SessionInfo *SessionInfo[SessionProperties]
	TrialInfo   *TrialInfo[TrialProperties]
}
