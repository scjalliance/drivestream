package collection

import (
	"fmt"
	"strings"
)

// Phase is a phase of collection.
type Phase int

// Collection phases.
const (
	PhaseDriveCollection  Phase = 0
	PhaseFileCollection   Phase = 1
	PhaseChangeCollection Phase = 2
	PhaseCommitProcessing Phase = 3
	PhaseFinalized        Phase = 4
	PhaseAbandoned        Phase = 1000
)

// String returns a string representation of p.
func (p Phase) String() string {
	switch p {
	case PhaseDriveCollection:
		return "drive collection"
	case PhaseFileCollection:
		return "file collection"
	case PhaseChangeCollection:
		return "change collection"
	case PhaseCommitProcessing:
		return "commit processing"
	case PhaseFinalized:
		return "finalized"
	default:
		return fmt.Sprintf("collection phase %d", p)
	}
}

// ParsePhase parses v as a collection phase.
func ParsePhase(v string) (Phase, error) {
	switch strings.ToLower(v) {
	case "drive collection":
		return PhaseDriveCollection, nil
	case "file collection":
		return PhaseFileCollection, nil
	case "change collection":
		return PhaseChangeCollection, nil
	case "commit processing":
		return PhaseCommitProcessing, nil
	case "finalized":
		return PhaseFinalized, nil
	default:
		return Phase(0), fmt.Errorf("unknown collection phase \"%s\"", v)
	}
}
