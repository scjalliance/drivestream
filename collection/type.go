package collection

import (
	"fmt"
	"strings"
)

// Type is a type of collection.
type Type int

// Collection types.
const (
	Full        Type = 0
	Incremental Type = 1
)

// String returns a string representation of t.
func (t Type) String() string {
	switch t {
	case Full:
		return "full"
	case Incremental:
		return "incremental"
	default:
		return fmt.Sprintf("collection type %d", t)
	}
}

// ParseType parses v as a collection type.
func ParseType(v string) (Type, error) {
	switch strings.ToLower(v) {
	case "full":
		return Full, nil
	case "incremental":
		return Incremental, nil
	default:
		return Type(0), fmt.Errorf("unknown collection type \"%s\"", v)
	}
}
