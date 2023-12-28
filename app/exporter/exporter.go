package exporter

import (
	"archive/zip"
	"encoding/json"
	"exporter/daylio"
	"exporter/types"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	DEFAULT_DESTINATION_JOURNAL          = "From Daylio"
	DAY_ONE_MAX_ENTRIES_IN_SINGLE_EXPORT = 99
	DEFAULT_EXPORT_DIRECTORY             = "./exports"
	BASE_FILE_NAME                       = "export"
	VERSION                              = "%%VER_CHANGED_BY_MAKE%%"
	COMMIT_SHA                           = "%%SHA_CHANGED_BY_MAKE%%"
)

type dayOneTimestamps struct {
	Created  types.DayOneDateTime
	Modified types.DayOneDateTime
}

// Version prints this app's version
func Version() {
	fmt.Printf("exporter version %s, commit %s\n", VERSION, COMMIT_SHA)
}

// Initializes sets up an export job.
func Initialize() error {
	log.Info("Starting Daylio to Day One export")
	setLogLevel()
	if err := createExportDirectoryIfMissing(); err != nil {
		return err
	}
	return nil
}

// ConvertToDayOneExportFromBackup converts entries within a Daylio backup file
// into a list of DayOne-compatible JSON import files.
func ConvertToDayOneExportFromDaylioBackup(providedFile string, generators types.DayOneGenerators) (*types.DayOneExport, error) {
	entries, err := daylio.GetEntriesFromBackupFile(providedFile)
	if err != nil {
		return nil, err
	}
	dayOneEntries, err := convertToDayOneEntries(entries, generators)
	if err != nil {
		return nil, err
	}
	return types.NewDayOneExport(dayOneEntries), nil
}

// ConvertToDayOneExportFromDaylioCSV converts entries within an exported CSV file from
// Daylio into a list of DayOne-compatible JSON import files.
func ConvertToDayOneExportFromDaylioCSV(daylioCSVPath string, generators types.DayOneGenerators) (*types.DayOneExport, error) {
	entries, err := daylio.GetEntriesFromCSVFile(daylioCSVPath)
	if err != nil {
		return nil, err
	}
	dayOneEntries, err := convertToDayOneEntries(entries, generators)
	if err != nil {
		return nil, err
	}
	return types.NewDayOneExport(dayOneEntries), nil
}

// WriteDayOneExports zips a DayOne export JSON and writes it to disk.
func WriteDayOneExports(export *types.DayOneExport) (*types.DayOneExportResult, error) {
	r := types.DayOneExportResult{
		ZipFile:     exportZipFileName(),
		JournalName: exportDayOneJournalName(),
	}
	f, err := os.Create(r.ZipFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	zip := zip.NewWriter(f)
	defer zip.Close()
	fInZip, err := zip.Create(r.JournalName + ".json")
	if err != nil {
		return nil, err
	}
	if err := writeDayOneExport(fInZip, export); err != nil {
		return nil, err
	}
	return &r, nil
}

func writeDayOneExport(buf io.Writer, export *types.DayOneExport) error {
	json, err := json.Marshal(export)
	if err != nil {
		return err
	}
	_, err = buf.Write(json)
	return err
}

func numEntriesInThisPage(entries []daylio.Entry, idx int) int {
	if len(entries) <= DAY_ONE_MAX_ENTRIES_IN_SINGLE_EXPORT {
		return len(entries) - 1
	}
	numEntries := idx + DAY_ONE_MAX_ENTRIES_IN_SINGLE_EXPORT
	if numEntries > len(entries) {
		numEntries = len(entries)
	}
	return numEntries
}

func convertToDayOneEntries(entries []daylio.Entry, generators types.DayOneGenerators) ([]types.DayOneEntry, error) {
	outs := []types.DayOneEntry{}
	for idx := 0; idx < len(entries); idx++ {
		daylioEntry := entries[idx]
		dayOneEntry := types.NewEmptyDayOneEntry()
		id := generators.IDGenerator.CreateID()
		rt, err := generateDayOneRichText(&daylioEntry, generators.UUIDGenerator)
		if err != nil {
			return nil, err
		}
		tags, err := generateTagsFromDaylioActivities(&daylioEntry)
		if err != nil {
			return nil, err
		}
		loc, err := generateLocationFromDaylioActivities(&daylioEntry)
		if err != nil {
			return nil, err
		}
		ts, err := createTimestamps(&daylioEntry, generators.Timestamper)
		if err != nil {
			return nil, err
		}
		dayOneEntry.RichText = rt
		dayOneEntry.UUID = id
		dayOneEntry.Tags = tags
		dayOneEntry.Location = loc
		dayOneEntry.CreationDate = ts.Created
		dayOneEntry.ModifiedDate = ts.Modified
		dayOneEntry.Text = createDayOneText(&daylioEntry)
		outs = append(outs, *dayOneEntry)
	}
	return outs, nil
}

func createDayOneText(entry *daylio.Entry) string {
	noteParts := make([]string, 2)
	if entry.NoteTitle != "" {
		noteParts[0] = entry.NoteTitle
	} else {
		noteParts[0] = "Note"
	}
	noteParts[1] = entry.Note
	return fmt.Sprintf("%s\n\n%s", noteParts[0], noteParts[1])
}

func generateDayOneRichText(entry *daylio.Entry, gen types.DayOneEntryUUIDGenerator) (string, error) {
	uuid, err := gen.GenerateUUID()
	if err != nil {
		return "", err
	}
	rt := types.DayOneRichTextObjectData{
		Meta: types.DayOneRichTextObjectDataMetadata{
			Version:           1,
			SmallLinesRemoved: false,
			Created: types.DayOneRichTextObjectCreatedProperties{
				Version:  1527,
				Platform: "com.bloombuilt.dayone-mac",
			},
		},
		Contents: []types.DayOneRichTextObject{
			{
				Text: createDayOneText(entry),
				Attributes: types.DayOneRichTextObjectAttributes{
					Line: types.DayOneRichTextLineObject{
						Header:     1,
						Identifier: uuid,
					},
				},
			},
		},
	}
	out, err := json.Marshal(rt)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func generateLocationFromDaylioActivities(entry *daylio.Entry) (types.DayOneEntryLocation, error) {
	if os.Getenv("NO_AUTO_HOME_LOCATION") != "" {
		return types.DayOneEntryLocation{}, nil
	}
	if os.Getenv("HOME_ADDRESS_JSON") == "" {
		return types.DayOneEntryLocation{}, nil
	}
	var out types.DayOneEntryLocation
	if err := json.Unmarshal([]byte(os.Getenv("HOME_ADDRESS_JSON")), &out); err != nil {
		return types.DayOneEntryLocation{}, err
	}
	return out, nil

}

func createTimestamps(entry *daylio.Entry, g types.DayOneEntryModifiedTimestamper) (dayOneTimestamps, error) {
	createdRaw := fmt.Sprintf("%sT%s:00Z", entry.FullDate, entry.Time)
	created, err := time.Parse("2006-01-02T15:04:05Z", createdRaw)
	if err != nil {
		return dayOneTimestamps{}, err
	}
	modified, err := g.CreateModifiedTime()
	if err != nil {
		return dayOneTimestamps{}, err
	}
	return dayOneTimestamps{
		Created:  types.DayOneDateTime(created),
		Modified: types.DayOneDateTime(modified),
	}, nil
}

func exportDirectory() string {
	return DEFAULT_EXPORT_DIRECTORY
}

func exportDayOneJournalName() string {
	return DEFAULT_DESTINATION_JOURNAL
}

func exportZipFileName() string {
	return filepath.Join(exportDirectory(), fmt.Sprintf("export-%s.zip", time.Now().Format("20060102")))
}

func createExportDirectoryIfMissing() error {
	_, err := os.Stat(exportDirectory())
	if err == nil {
		return nil
	}
	if exists := os.IsExist(err); !exists {
		log.Debugf("Creating export directory: %s", exportDirectory())
		return os.Mkdir(exportDirectory(), 0755)
	}
	return err
}

func setLogLevel() {
	if os.Getenv("LOG_LEVEL") == "" {
		return
	}
	level, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.Warningf("Invalid log level, using default: %s", os.Getenv("LOG_LEVEL"))
	}
	log.SetLevel(level)
}
