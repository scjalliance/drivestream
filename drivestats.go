package drivestream

import (
	"fmt"

	"code.cloudfoundry.org/bytefmt"
)

// DriveStats hold statistics about a drive.
type DriveStats struct {
	Count           int64
	TotalBytes      int64
	Collections     int64
	CollectionBytes int64
	Commits         int64
	CommitBytes     int64
	Versions        int64
	VersionBytes    int64
	ViewCommits     int64
	ViewBytes       int64
	Files           FileStats
}

// Summary returns a slice of strings summarizing the statistics.
func (ds DriveStats) Summary() []string {
	var output []string
	output = append(output, fmt.Sprintf("Drives: %d (%s)", ds.Count, bytefmt.ByteSize(uint64(ds.TotalBytes))))
	output = append(output, fmt.Sprintf("  Collections: %d (%s)", ds.Collections, bytefmt.ByteSize(uint64(ds.CollectionBytes))))
	output = append(output, fmt.Sprintf("  Commits: %d (%s)", ds.Commits, bytefmt.ByteSize(uint64(ds.CommitBytes))))
	output = append(output, fmt.Sprintf("  Drive Versions: %d (%s)", ds.Versions, bytefmt.ByteSize(uint64(ds.VersionBytes))))
	output = append(output, fmt.Sprintf("  Drive View Commits: %d (%s)", ds.ViewCommits, bytefmt.ByteSize(uint64(ds.ViewBytes))))
	output = append(output, fmt.Sprintf("Files: %d (%s)", ds.Files.Count, bytefmt.ByteSize(uint64(ds.Files.TotalBytes))))
	output = append(output, fmt.Sprintf("  File Versions: %d (%s)", ds.Files.Versions, bytefmt.ByteSize(uint64(ds.Files.VersionBytes))))
	output = append(output, fmt.Sprintf("  File View Commits: %d (%s)", ds.Files.ViewCommits, bytefmt.ByteSize(uint64(ds.Files.ViewBytes))))
	return output
}
