package collection

import "fmt"

// SeqNum is a collection sequence number.
type SeqNum int64

// String returns a string representation of the collection sequence number.
func (number SeqNum) String() string {
	return fmt.Sprintf("C%10d", number)
}
