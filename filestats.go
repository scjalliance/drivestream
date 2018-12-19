package drivestream

// FileStats hold statistics about one or more files.
type FileStats struct {
	Count        int64
	TotalBytes   int64
	Versions     int64
	VersionBytes int64
	Views        int64
	ViewCommits  int64
	ViewBytes    int64
}
