package resource

// An ID is a file or team drive identifier.
type ID string

// String returns id as a string.
func (id ID) String() string {
	return string(id)
}
