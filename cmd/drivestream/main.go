package main

import (
	"os"
	"syscall"

	"github.com/gentlemanautomaton/signaler"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		app             = kingpin.New("drivestream", "Collects and preserves team drive metadata.")
		dbType          = app.Flag("db", "database type").Default("bolt").Envar("DB_TYPE").String()
		dbPath          = app.Flag("file", "database file path").Default("drivestream.db").Envar("DB_PATH").String()
		includeMemStats = app.Flag("memstats", "include memory statistics in output").Envar("INCLUDE_MEMORY_STATS").Bool()
		updateCommand   = app.Command("update", "Collects metadata and updates a drivestream database.")
		updateEmail     = updateCommand.Flag("email", "email address of group or account to use during collection").Envar("GOOGLE_ACCOUNT").Required().String()
		updateInterval  = updateCommand.Flag("interval", "interval between updates").Short('i').Envar("INTERVAL").Duration()
		updateWanted    = updateCommand.Arg("wanted", "team drives to update (name or ID)").Strings()
		statsCommand    = app.Command("stats", "Reports statistics about a drivestream database.")
		statsSelections = statsCommand.Flag("select", "statistics to select").Short('s').Default("collections", "commits").Strings()
		statsWanted     = statsCommand.Arg("wanted", "team drives to report statistics for (name or ID)").Strings()
		dumpCommand     = app.Command("dump", "Dumps team drive metadata currently stored within a drivestream database.")
		dumpSelections  = dumpCommand.Flag("selection", "kinds of data to dump").Short('s').Default("collections", "commits").Strings()
		dumpWanted      = dumpCommand.Arg("wanted", "team drives to dump (name or ID)").Strings()
	)

	shutdown := signaler.New().Capture(os.Interrupt, syscall.SIGTERM)
	defer shutdown.Trigger()
	ctx := shutdown.Context()

	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	repo, repoClose := NewRepository(app, *dbType, *dbPath)
	defer repoClose()

	switch command {
	case updateCommand.FullCommand():
		update(ctx, app, repo, *includeMemStats, *updateEmail, *updateInterval, *updateWanted)
	case statsCommand.FullCommand():
		stats(ctx, app, repo, *statsSelections, *statsWanted)
	case dumpCommand.FullCommand():
		dump(ctx, app, repo, *dumpSelections, *dumpWanted)
	}
}
