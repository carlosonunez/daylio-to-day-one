package main

import (
	"exporter/exporter"
	"exporter/types"
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

const (
	USAGE = `Usage: daylio-to-day-one [FILE]
Exports entries in a Daylio backup file to a Day One JSON ZIP file.

OPTIONS

	FILE			The path to the Daylio backup file. Optional if
						iCloud Backup is enabled within Daylio.

GENERATING DAYLIO EXPORT FILES

Do the following to generate a Daylio backup file and provide it to the Daylio to Day One Exporter:

  * Open Daylio,
  * Tap the "(...) More" button on the far right,
  * Tap "Backup & Restore"
  * Tap "Advanced Options"
  * Tap "Export". Save the file somewhere convenient, like
	* "Downloads/daylio.backup"
	* Copy this file to the computer running this program.
	* Provide the backup file to Exporter:  "daylio-to-day-one Downloads/daylio.backup"
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
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		exporter.Version()
		os.Exit(0)
	}
	if err := exporter.Initialize(); err != nil {
		log.Errorf("Something went wrong while initializing the exporter: %s", err.Error())
	}
	providedBackupFile := ""
	if len(os.Args) == 2 {
		providedBackupFile = os.Args[1]
	}
	dayOneExports, err := exporter.ConvertToDayOneExportFromDaylioBackup(providedBackupFile, types.DefaultDayOneGenerators())
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
