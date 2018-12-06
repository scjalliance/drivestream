package collection

import "fmt"

// StateNum is a collection state number.
type StateNum int64

// String returns a string representation of the collection state number.
func (number StateNum) String() string {
	return fmt.Sprintf("S%10d", number)
}
