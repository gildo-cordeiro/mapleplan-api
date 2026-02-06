package goal

type Phase string

const (
	PreDeparture Phase = "pre_departure"
	Arrival      Phase = "arrival"
	PostArrival  Phase = "post_arrival"
	UnknownPhase Phase = "unknown"
)

func IsValidPhase(phase Phase) bool {
	switch phase {
	case PreDeparture, Arrival, PostArrival:
		return true
	default:
		return false
	}
}

func GetAllPhases() []Phase {
	return []Phase{PreDeparture, Arrival, PostArrival}
}

func PhaseToString(phase Phase) string {
	return string(phase)
}

func StringToPhase(phaseStr string) (Phase, bool) {
	phase := Phase(phaseStr)
	if IsValidPhase(phase) {
		return phase, true
	}
	return UnknownPhase, false
}
