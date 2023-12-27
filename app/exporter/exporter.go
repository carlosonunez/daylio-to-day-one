package exporter

import (
	"encoding/json"
	"exporter/types"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	csv "github.com/gocarina/gocsv"
	log "github.com/sirupsen/logrus"
)

const (
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

// ConvertToDayOneExport converts entries within an exported CSV file from
// Daylio into a list of DayOne-compatible JSON import files.
//
// Each Day One export file contains at most 99 entries, as this seems to be the
// most entries Day One will process at a time.
func ConvertToDayOneExport(daylioCSVPath string, generators types.DayOneGenerators) ([]*types.DayOneExport, error) {
	f, err := os.OpenFile(daylioCSVPath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var entries []types.DaylioEntry
	if err := csv.UnmarshalFile(f, &entries); err != nil {
		return nil, err
	}
	var exports []*types.DayOneExport
	for idx := 0; idx < len(entries); idx++ {
		start := idx
		end := numEntriesInThisPage(entries, idx)
		entries_text := "Daylio Entries"
		if idx > 0 {
			entries_text = "more Daylio Entries"
		}
		log.Infof("Converting %d %s to Day One (%d remaining)", end-start, entries_text, len(entries)-end)
		batch := entries[start:end]
		dayOneEntries, err := convertToDayOneEntries(batch, generators)
		if err != nil {
			return nil, err
		}
		exports = append(exports, types.NewDayOneExport(dayOneEntries))
		idx = end
	}
	return exports, nil
}

// WriteDayOneExports commits a list of DayOneExports to disk.
func WriteDayOneExports(exportList []*types.DayOneExport) ([]string, error) {
	exportFnames := []string{}
	today := time.Now()
	for page, export := range exportList {
		log.Debugf("On page %d", page)
		fname := exportFileName(page, today)
		f, err := os.Create(fname)
		if err != nil {
			return []string{}, err
		}
		defer f.Close()
		if err := writeDayOneExport(f, export); err != nil {
			return []string{}, err
		}
		exportFnames = append(exportFnames, fname)
	}
	return exportFnames, nil
}

func writeDayOneExport(buf io.Writer, export *types.DayOneExport) error {
	json, err := json.Marshal(export)
	if err != nil {
		return err
	}
	_, err = buf.Write(json)
	return err
}

func numEntriesInThisPage(entries []types.DaylioEntry, idx int) int {
	if len(entries) <= DAY_ONE_MAX_ENTRIES_IN_SINGLE_EXPORT {
		return len(entries) - 1
	}
	numEntries := idx + DAY_ONE_MAX_ENTRIES_IN_SINGLE_EXPORT
	if numEntries > len(entries) {
		numEntries = len(entries)
	}
	return numEntries
}

func convertToDayOneEntries(entries []types.DaylioEntry, generators types.DayOneGenerators) ([]types.DayOneEntry, error) {
	outs := []types.DayOneEntry{}
	for idx := 0; idx < len(entries); idx++ {
		daylioEntry := entries[idx]
		dayOneEntry := types.NewEmptyDayOneEntry()
		id := generators.IDGenerator.CreateID(&daylioEntry)
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

func createDayOneText(entry *types.DaylioEntry) string {
	noteParts := make([]string, 2)
	if entry.NoteTitle != "" {
		noteParts[0] = entry.NoteTitle
	} else {
		noteParts[0] = "Note"
	}
	noteParts[1] = entry.Note
	return fmt.Sprintf("%s\n\n%s", noteParts[0], noteParts[1])
}

func generateDayOneRichText(entry *types.DaylioEntry, gen types.DayOneEntryUUIDGenerator) (string, error) {
	uuid, err := gen.GenerateUUID(entry)
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

func generateTagsFromDaylioActivities(entry *types.DaylioEntry) ([]string, error) {
	var out []string
	for _, activityRaw := range strings.Split(entry.Activities, "|") {
		activity := strings.Trim(activityRaw, " ")
		if aloneScore := generateAloneTimeScore(activity); aloneScore != "" {
			activity = aloneScore
		}
		out = append(out, activity)
	}
	return out, nil
}

func generateLocationFromDaylioActivities(entry *types.DaylioEntry) (types.DayOneEntryLocation, error) {
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

func generateAloneTimeScore(s string) string {
	if os.Getenv("NO_ALONE_TIME_SCORING") != "" {
		return ""
	}
	var score int
	switch strings.ToLower(s) {
	case "no":
		score = 0
	case "a little bit":
		score = 1
	case "yes!":
		score = 2
	default:
		return ""
	}
	return fmt.Sprintf("alone score: %d", score)
}

func createTimestamps(entry *types.DaylioEntry, g types.DayOneEntryModifiedTimestamper) (dayOneTimestamps, error) {
	createdRaw := fmt.Sprintf("%sT%sZ", entry.FullDate, entry.Time)
	created, err := time.Parse("2006-01-02T15:04Z", createdRaw)
	if err != nil {
		return dayOneTimestamps{}, err
	}
	modified, err := g.CreateModifiedTime(entry)
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

func exportFileName(page int, today time.Time) string {
	fname := fmt.Sprintf("%s-%s", BASE_FILE_NAME, today.Format("20060102"))
	thisPage := page + 1
	if thisPage > 1 {
		fname = fmt.Sprintf("%s-%d", fname, thisPage)
	}
	fname = fname + ".json"
	return path.Join(exportDirectory(), fname)
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
