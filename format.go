package tournaments

type Phase string

// type Format map[Phase]map[string]Phase

const (
	PhaseRegistration   = "registration"
	PhaseInitialization = "initialization"
	PhaseDraft          = "draft"
	PhaseRounds         = "rounds"
	PhaseEnded          = "ended"
)

// const (
// 	FormatPauperCube = "pauper-cube"
// )

type Format interface {
	Type() string
	On(action string, p Phase) Phase
}

type FormatPauperCube struct {
	Type string
}

func (pc *FormatPauperCube) On(action string, p Phase) Phase {
	switch p {
	case PhaseInitialization:
		switch action {
		case ActionEndPhase:
			return PhaseRegistration
		}
	case PhaseRegistration:
		switch action {
		case ActionEndPhase:
			return PhaseDraft
		}
	case PhaseDraft:
		switch action {
		case ActionEndPhase:
			return PhaseRounds
		}
	case PhaseRounds:
		switch action {
		case ActionEndPhase:
			return PhaseEnded
		}
	}
	return p
}

// var (
// 	FormatPauperCube = Format{
// 		PhaseInitialization: map[string]Phase{
// 			ActionEndPhase: PhaseRegistration,
// 		},
// 		PhaseRegistration: map[string]Phase{
// 			ActionEndPhase: PhaseDraft,
// 		},
// 		PhaseDraft: map[string]Phase{
// 			ActionEndPhase: PhaseRounds,
// 		},
// 		PhaseRounds: map[string]Phase{
// 			ActionEndPhase: PhaseEnded,
// 		},
// 	}
// )
