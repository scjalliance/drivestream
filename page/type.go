package page

import (
	"fmt"
	"strings"
)

// Type is a type of page.
type Type int

// Page types.
const (
	DriveList  Type = 0
	FileList   Type = 1
	ChangeList Type = 2
)

// String returns a string representation of t.
func (t Type) String() string {
	switch t {
	case DriveList:
		return "drive list"
	case FileList:
		return "file list"
	case ChangeList:
		return "change list"
	default:
		return fmt.Sprintf("page type %d", t)
	}
}

// ParseType parses v as a page type.
func ParseType(v string) (Type, error) {
	switch strings.ToLower(v) {
	case "drive list":
		return DriveList, nil
	case "file list":
		return FileList, nil
	case "change list":
		return ChangeList, nil
	default:
		return Type(0), fmt.Errorf("unknown page type \"%s\"", v)
	}
}
