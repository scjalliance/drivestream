package page

import "fmt"

// SeqNum is a page sequence number.
type SeqNum int64

// String returns a string representation of the page sequence number.
func (number SeqNum) String() string {
	return fmt.Sprintf("P%10d", number)
}
