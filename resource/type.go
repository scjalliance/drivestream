package resource

import (
	"fmt"
	"strings"
)

// Type is a type of resource.
type Type int

// Page types.
const (
	TypeDrive Type = 0
	TypeFile  Type = 1
)

// String returns a string representation of t.
func (t Type) String() string {
	switch t {
	case TypeDrive:
		return "drive#teamDrive"
	case TypeFile:
		return "drive#file"
	default:
		return fmt.Sprintf("resource type %d", t)
	}
}

// ParseType parses v as a page type.
func ParseType(v string) (Type, error) {
	switch strings.ToLower(v) {
	case "drive#teamDrive":
		return TypeDrive, nil
	case "drive#file":
		return TypeFile, nil
	default:
		return Type(0), fmt.Errorf("unknown resource type \"%s\"", v)
	}
}
