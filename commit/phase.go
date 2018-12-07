package commit

import (
	"fmt"
	"strings"
)

// A Phase is a phase of commit processing.
type Phase int

// Commit phases.
const (
	PhaseSourceProcessing Phase = 0
	PhaseTreeProcessing   Phase = 1
	PhaseFinalized        Phase = 2
)

// String returns a string representation of p.
func (p Phase) String() string {
	switch p {
	case PhaseSourceProcessing:
		return "source processing"
	case PhaseTreeProcessing:
		return "tree processing"
	case PhaseFinalized:
		return "finalized"
	default:
		return fmt.Sprintf("commit phase %d", p)
	}
}

// ParsePhase parses v as a commit phase.
func ParsePhase(v string) (Phase, error) {
	switch strings.ToLower(v) {
	case "source processing":
		return PhaseSourceProcessing, nil
	case "tree processing":
		return PhaseTreeProcessing, nil
	case "finalized":
		return PhaseFinalized, nil
	default:
		return Phase(0), fmt.Errorf("unknown commit phase \"%s\"", v)
	}
}
