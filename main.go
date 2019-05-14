package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/rdooley/confidynt/cli"
)

var (
	app    = kingpin.New("confidynt", "A command-line application for 12 factor config dynamo db management")
	dryrun = app.Flag("dryrun", "Enable dryrun mode.").Bool()
	table  = app.Flag("table", "Table name").Required().String()

	read      = app.Command("read", "Read a config from dynamo")
	readKey   = read.Arg("key", "Query key").Required().String()
	readValue = read.Arg("value", "Query value").Required().String()

	write     = app.Command("write", "Write a config file to dynamo")
	writeFile = write.Arg("config", "Config file to write").ExistingFile()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	// Read
	case read.FullCommand():
		cli.Read(*table, *readKey, *readValue)
	// Write
	case write.FullCommand():
		cli.Write(*table, *writeFile)
	}
}
