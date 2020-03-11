package tournaments

type Phase string

const (
	PhaseRegistration   = "registration"
	PhaseInitialization = "initialization"
	PhaseDraft          = "draft"
	PhaseRounds         = "rounds"
	PhaseEnded          = "ended"
)
