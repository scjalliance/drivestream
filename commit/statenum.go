package commit

import "fmt"

// StateNum is a commit state number.
type StateNum int64

// String returns a string representation of the commit state number.
func (number StateNum) String() string {
	return fmt.Sprintf("S%10d", number)
}
