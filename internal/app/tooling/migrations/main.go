package main

import (
	"expvar"
	"flag"
	"fmt"
	"github.com/AgeroFlynn/crud/foundation/config"
	"os"

	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
)

const usageText = `This program runs command on the db. Supported commands are:
  - init - creates version info table in the database
  - up - runs all available migrations.
  - up [target] - runs available migrations up to the target one.
  - down - reverts last migration.
  - reset - reverts all migrations.
  - version - prints current db version.
  - set_version [version] - sets db version without running migrations.
Usage:
  go run *.go <command> [args]
`

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func main() {

	// =========================================================================
	// Configuration

	expvar.NewString("build").Set(build)

	cfg, err := config.NewConfigFromFile(build)
	if err != nil {
		fmt.Errorf("parsing config: %w", err)
		os.Exit(1)
	}
	// =========================================================================

	flag.Usage = usage
	flag.Parse()

	db := pg.Connect(&pg.Options{
		User:     cfg.DB.User,
		Database: cfg.DB.Name,
	})

	oldVersion, newVersion, err := migrations.Run(db, flag.Args()...)
	if err != nil {
		exitf(err.Error())
	}
	if newVersion != oldVersion {
		fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		fmt.Printf("version is %d\n", oldVersion)
	}
}

func usage() {
	fmt.Print(usageText)
	flag.PrintDefaults()
	os.Exit(2)
}

func errorf(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s+"\n", args...)
}

func exitf(s string, args ...interface{}) {
	errorf(s, args...)
	os.Exit(1)
}
