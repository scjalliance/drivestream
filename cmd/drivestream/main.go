package main

import (
	"os"
	"syscall"

	"github.com/gentlemanautomaton/signaler"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		app            = kingpin.New("drivestream", "Collects and preserves team drive metadata.")
		dbType         = app.Flag("db", "database type").Default("bolt").Envar("DB_TYPE").String()
		dbPath         = app.Flag("file", "database file path").Default("drivestream.db").Envar("DB_PATH").String()
		includeStats   = app.Flag("stats", "include statistics in output").Envar("INCLUDE_STATISTICS").Bool()
		updateCommand  = app.Command("update", "Collects metadata and updates a drivestream database.")
		updateEmail    = updateCommand.Flag("email", "email address of group or account to use during collection").Envar("GOOGLE_ACCOUNT").Required().String()
		updateInterval = updateCommand.Flag("interval", "interval between updates").Short('i').Envar("INTERVAL").Duration()
		updateWanted   = updateCommand.Arg("wanted", "team drives to update (name or ID)").Strings()
		dumpCommand    = app.Command("dump", "Dumps team drive metadata currently stored within a drivestream database.")
		dumpKinds      = dumpCommand.Flag("kind", "kinds of data to dump").Short('k').Default("collections").Strings()
		dumpWanted     = dumpCommand.Arg("wanted", "team drives to dump (name or ID)").Strings()
	)

	shutdown := signaler.New().Capture(os.Interrupt, syscall.SIGTERM)
	defer shutdown.Trigger()
	ctx := shutdown.Context()

	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	db := NewDB(app, *dbType, *dbPath)
	defer db.Close()

	switch command {
	case updateCommand.FullCommand():
		update(ctx, app, db, *includeStats, *updateEmail, *updateInterval, *updateWanted)
	case dumpCommand.FullCommand():
		dump(ctx, app, db, *dumpKinds, *dumpWanted)
	}
}
