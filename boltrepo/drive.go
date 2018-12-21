package boltrepo

import (
	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream"
	"github.com/scjalliance/drivestream/binpath"
	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/driveversion"
	"github.com/scjalliance/drivestream/driveview"
	"github.com/scjalliance/drivestream/resource"
)

var _ drivestream.DriveReference = (*Drive)(nil)

// Drive is a drivestream drive reference for a bolt repository.
type Drive struct {
	db    *bolt.DB
	drive resource.ID
}

// Path returns the path of the drive.
func (ref Drive) Path() binpath.Text {
	return binpath.Text{RootBucket, DriveBucket, ref.drive.String()}
}

// DriveID returns the resource ID of the drive.
func (ref Drive) DriveID() resource.ID {
	return ref.drive
}

// Exists returns true if the drive exists.
func (ref Drive) Exists() (exists bool, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		if driveBucket(tx, ref.drive) != nil {
			exists = true
		}
		return nil
	})
	return exists, err
}

// Collections returns the collection sequence for the drive.
func (ref Drive) Collections() collection.Sequence {
	return Collections{
		db:    ref.db,
		drive: ref.drive,
	}
}

// Collection returns a collection reference. Equivalent to Collections().Ref(s).
func (ref Drive) Collection(c collection.SeqNum) collection.Reference {
	return Collection{
		db:         ref.db,
		drive:      ref.drive,
		collection: c,
	}
}

// Commits returns the commit sequence for the drive.
func (ref Drive) Commits() commit.Sequence {
	return Commits{
		db:    ref.db,
		drive: ref.drive,
	}
}

// Commit returns a commit reference. Equivalent to Commits().Ref(s).
func (ref Drive) Commit(c commit.SeqNum) commit.Reference {
	return Commit{
		db:     ref.db,
		drive:  ref.drive,
		commit: c,
	}
}

// Versions returns the version sequence for the drive.
func (ref Drive) Versions() driveversion.Sequence {
	return DriveVersions{
		db:    ref.db,
		drive: ref.drive,
	}
}

// Version returns a drive version reference. Equivalent to Versions().Ref(s).
func (ref Drive) Version(v resource.Version) driveversion.Reference {
	return DriveVersion{
		db:      ref.db,
		drive:   ref.drive,
		version: v,
	}
}

// View returns a view of the drive.
func (ref Drive) View() driveview.Reference {
	return DriveView{
		db:    ref.db,
		drive: ref.drive,
	}
}

// At returns a version reference of the drive at a particular commit.
func (ref Drive) At(seqNum commit.SeqNum) (driveversion.Reference, error) {
	return ref.View().At(seqNum)
}

// Stats returns statistics about the drive.
func (ref Drive) Stats() (stats drivestream.DriveStats, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		drv := driveBucket(tx, ref.drive)
		if drv == nil {
			return nil
		}
		stats.Count++
		stats.TotalBytes += countBytes(drv)
		if collections := collectionsBucket(tx, ref.drive); collections != nil {
			cursor := collections.Cursor()
			for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
				stats.Collections++
				stats.CollectionBytes += int64(len(k)) + int64(len(v))
				stats.CollectionBytes += countBytes(collections.Bucket(k))
			}
		}
		if commits := commitsBucket(tx, ref.drive); commits != nil {
			cursor := commits.Cursor()
			for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
				stats.Commits++
				stats.CommitBytes += int64(len(k)) + int64(len(v))
				stats.CommitBytes += countBytes(commits.Bucket(k))
			}
		}
		if versions := driveVersionsBucket(tx, ref.drive); versions != nil {
			cursor := versions.Cursor()
			for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
				stats.Versions++
				stats.VersionBytes += int64(len(k)) + int64(len(v))
				stats.VersionBytes += countBytes(versions.Bucket(k))
			}
		}
		if view := driveViewBucket(tx, ref.drive); view != nil {
			cursor := view.Cursor()
			for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
				stats.ViewCommits++
				stats.ViewBytes += int64(len(k)) + int64(len(v))
				stats.ViewBytes += countBytes(view.Bucket(k))
			}
		}
		if files := filesBucket(tx); files != nil {
			cursor := files.Cursor()
			for f, _ := cursor.First(); f != nil; f, _ = cursor.Next() {
				file := files.Bucket(f)
				if file == nil {
					continue
				}
				views := file.Bucket([]byte(ViewBucket))
				if views == nil {
					continue
				}
				view := views.Bucket([]byte(ref.drive))
				if view == nil {
					continue
				}
				stats.Files.Count++
				stats.Files.TotalBytes += countBytes(file)
				stats.Files.Views++
				{
					cursor := view.Cursor()
					for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
						stats.Files.ViewCommits++
						stats.Files.ViewBytes += int64(len(k)) + int64(len(v))
						stats.Files.ViewBytes += countBytes(view.Bucket(k))
					}
				}
				if versions := file.Bucket([]byte(VersionBucket)); versions != nil {
					cursor := versions.Cursor()
					for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
						stats.Files.Versions++
						stats.Files.VersionBytes += int64(len(k)) + int64(len(v))
						stats.Files.VersionBytes += countBytes(versions.Bucket(k))
					}
				}
			}
		}
		return nil
	})
	return stats, err
}
