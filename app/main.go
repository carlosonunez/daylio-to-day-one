package main

import (
	"exporter/exporter"
	"exporter/types"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	USAGE = `Usage: daylio-to-day-one [FILE]
Converts a Daylio CSV export to importable Day One JSON files

OPTIONS

	FILE			The path to the Daylio export.

NOTES

- This app converts 99 Daylio entries at a time. This seems to be a Day One limitation.
`
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Print(USAGE)
		os.Exit(1)
	}
	if os.Args[1] == "-v" || os.Args[1] == "--version" {
		exporter.Version()
		os.Exit(0)
	}
	daylioCSVFile := os.Args[1]
	if err := exporter.Initialize(); err != nil {
		log.Errorf("Something went wrong while initializing the exporter: %s", err.Error())
	}
	dayOneExports, err := exporter.ConvertToDayOneExport(daylioCSVFile, types.DefaultDayOneGenerators())
	if err != nil {
		log.Errorf("Something went wrong while performing the export: %s", err.Error())
		os.Exit(1)
	}
	fileList, err := exporter.WriteDayOneExports(dayOneExports)
	if err != nil {
		log.Errorf("Something went wrong while writing the exports: %s", err.Error())
		os.Exit(1)
	}
	log.Infof("Your Day One exports are ready. You can find them here: %s", strings.Join(fileList, ", "))
}
