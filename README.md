drivestream [![GoDoc](https://godoc.org/github.com/scjalliance/drivestream?status.svg)](https://godoc.org/github.com/scjalliance/drivestream)
====

A library for collecting and preserving Google Team Drive history.

This library is a work in progress and is subject to breaking changes.

## Quick Start

```
// Prepare a context
ctx := context.Background()

// Prepare a Google Drive API client
driveService, _ := drive.New(oauthHttpClient)

// Determine the ID of the team drive to be used
teamDriveID := resource.ID("TEAMDRIVEID")

// Prepare an in-memory repository
repo := memrepo.New()

// Prepare a collector that will query the drive service
collector := driveapicollector.New(driveService, teamDriveID)

// Create a stream
stream := drivestream.New(repo, teamDriveID, drivestream.WithLogging(os.Stdout))

// Initialize or update the repository
stream.Update(ctx, collector) // Will perform first-run initialization and resume incomplete updates automatically

// Update the database every minute
timer := time.NewTimer(time.Minute)
for {
    select {
    case <-timer.C:
        stream.Update(ctx, collector)
    case <-ctx.Done():
        return
    }
}
```

## Stream

The `drivestream` library is intended to be accessed through a `Stream`
type, which wraps a `Repository` implementation that the caller provides.

A stream provides two capabilities:

1. Collection and persistence of data through calls to `Update()`
2. Access to collected data through a `Cursor` (this is not yet implemented)

## Repository

Data collected from a Team Drive is preserved in a repository. The repository
implementation is pluggable, but drivestream was designed with key/value
storage in mind.

The drivestream project supplies the following implementations of `Repository`:

* `memrepo`: An in-memory repository useful for testing.
* `boltrepo`: A repository backed by a bolt database.

TODO: Add support for `badger`.

## Collection

Data is brought into a drivestream repository through a series of collections.
Each collection is assigned a monotonically increasing sequence number. The
initial collection is assigned sequence number zero.

During collection a data source is queried to retrieve the following metadata:

1. Information about the team drive, such as its name
2. The complete set of files within the team drive
3. Subsequent changes to the team drive or its contents

A collection is designed to be resumable in case of interruption or service
failure. Data is progressively added to a collection as a series of data
pages. A collection can be resumed from the last page that was
successfully written.

Collection relies upon a `Collector` implementation, which is responsible
for querying a particular kind of data source.

The drivestream project supplies the following implementations of `Collector`:

* `driveapicollector`: A collector that queries Google Drive API version 3

Once collected, data is processed and reformulated into a series of commits.

Once finished with its commit processing, a collection moves to moves to a
`Finalized` state. Collections in a finalized state cannot be modified.

## Commits

A commit represents a consistent view of the entire team drive at a point in
time. Like collections, each commit is assigned a monotonically increasing
sequence number. The initial commit is assigned sequence number zero.

A collection is processed into one or more commits. The first commit includes
all files present in the team drive when the first collection was performed.
Subsequent commits typically contain only a single change.

Commits are constructed in two phases:

1. Version Processing
2. Tree Processing

Each phase is resumable in case of interruption or service failure.

Once a commit has completed it moves to a `Finalized` state. Commits in a
finalized state cannot be modified.

## Commit Version Processing

TODO: Design this process.

## Commit Tree Processing

TODO: Design this process.

## Database Schema

A work-in-progress key/value database schema:

```
Key                                                                Value
----------------------------------------------------------------   ----------------
/schema                                                            "{SCHEMA_VERSION}"

/drive/{DRIVE_ID}/collection/{COLLECTION_NUM}/data                 "{JSON(COLLECTION_DATA)}"
/drive/{DRIVE_ID}/collection/{COLLECTION_NUM}/state/{STATE_NUM}    "{JSON(COLLECTION_STATE)}"
/drive/{DRIVE_ID}/collection/{COLLECTION_NUM}/page/{PAGE_NUM}      "{JSON(PAGE_DATA)}"

/drive/{DRIVE_ID}/commit/{COMMIT_NUM}/data                         "{JSON(COMMIT_DATA)}"
/drive/{DRIVE_ID}/commit/{COMMIT_NUM}/state/{STATE_NUM}            "{JSON(COMMIT_STATE)}"
/drive/{DRIVE_ID}/commit/{COMMIT_NUM}/file/{FILE_ID}               "{VERSION}"
/drive/{DRIVE_ID}/commit/{COMMIT_NUM}/tree/{PARENT_ID}/{CHILD_ID}  "{ACTION}"

/drive/{DRIVE_ID}/time/collection/{TIME}                           "{COLLECTION_NUM}"
/drive/{DRIVE_ID}/time/commit/{TIME}                               "{COMMIT_NUM}"

/drive/{DRIVE_ID}/version/{VERSION}                                "{JSON(DRIVE_DATA)}"
/drive/{DRIVE_ID}/view/{COMMIT_NUM}                                "{VERSION}"
/drive/{DRIVE_ID}/tree/{COMMIT_NUM}                                "{HASH(FILE_LIST|CHUNK_LIST)}"

/file/{FILE_ID}/version/{VERSION}                                  "{JSON(FILE_DATA)}"
/file/{FILE_ID}/view/{DRIVE_ID}/{COMMIT_NUM}                       "{VERSION}"
/file/{FILE_ID}/tree/{DRIVE_ID}/{COMMIT_NUM}                       "{HASH(FILE_LIST|CHUNK_LIST)}"
/file/{FILE_ID}/time/{DRIVE_ID}/{TIME}                             "{COMMIT_NUM}"

/tree/hash/{HASH(FILE_LIST|CHUNK_LIST)}                            "{BINARY(FILE_LIST|CHUNK_LIST)}"
```
