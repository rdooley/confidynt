package main

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
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
		fmt.Printf("READING FROM %s\n", *table)
		fmt.Printf("%s=%s\n", *readKey, *readValue)

	// Write
	case write.FullCommand():
		fmt.Printf("WRITING TO %s\n", *table)
		fmt.Printf("%s is written\n", *writeFile)
	}
}
