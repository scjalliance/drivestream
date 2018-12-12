package commit

// FileMap is an unordered map of file changes.
type FileMap interface {
	// Read returns the set of file changes for the commit, in unspecified
	// order.
	Read() (changes []FileChange, err error)

	// Add adds the given file changes to the map.
	// If two or more changes conflict, the last change added takes
	// precedence.
	Add(changes ...FileChange) error
}
