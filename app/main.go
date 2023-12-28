package main

import (
	"exporter/exporter"
	"exporter/types"
	"fmt"
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

const (
	USAGE = `Usage: daylio-to-day-one [FILE]
Exports entries in a Daylio CSV to a Day One JSON ZIP file.

OPTIONS

	FILE			The path to the Daylio export.
`
)

func printSuccessMessage(r *types.DayOneExportResult) {
	zf, err := filepath.Abs(r.ZipFile)
	if err != nil {
		panic(err)
	}
	log.Infof(`Your Day One JSON ZIP file is ready! Do the following on this computer to finish \
importing your Daylio entries into Day One:

1. Open the Day One app.
2. Click on 'File', then 'Import', then 'JSON ZIP File'.
3. Browse to this folder: %s
4. Click on this file, then on Open: %s

Your journal entries will appear in a new Day One journal called "%s". You can leave them there
or move them into your desired journal.
`, path.Dir(zf), filepath.Base(zf), r.JournalName)
}

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
	result, err := exporter.WriteDayOneExports(dayOneExports)
	if err != nil {
		log.Errorf("Something went wrong while writing the exports: %s", err.Error())
		os.Exit(1)
	}
	printSuccessMessage(result)
}
