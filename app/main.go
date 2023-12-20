package exporter

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	csv "github.com/gocarina/gocsv"
	log "github.com/sirupsen/logrus"
)

type dayOneTimestamps struct {
	Created  time.Time
	Modified time.Time
}

func convertToDayOne(entry *DaylioEntry, gList *dayOneGenerators) (*DayOneEntry, error) {
	return nil, nil
}

func convertToDayOneExport(inCSV string, generators dayOneGenerators) (*DayOneExport, error) {
	var entries []DaylioEntry
	if err := csv.UnmarshalString(inCSV, &entries); err != nil {
		return nil, err
	}
	out := DayOneExport{
		Metadata: DayOneMetadata{
			Version: "1.0",
		},
		Entries: []DayOneEntry{},
	}
	for idx := 0; idx < len(entries); idx++ {
		daylioEntry := entries[idx]
		dayOneEntry := NewEmptyDayOneEntry()
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
		dayOneEntry.RichText = rt
		dayOneEntry.UUID = id
		dayOneEntry.Tags = tags
		dayOneEntry.Location = loc
		out.Entries = append(out.Entries, *dayOneEntry)
	}
	return &out, nil
}

func createRichTextNote(entry *DaylioEntry) string {
	noteParts := make([]string, 2)
	if entry.NoteTitle != "" {
		noteParts[0] = entry.NoteTitle
	} else {
		noteParts[0] = "Note"
	}
	noteParts[1] = entry.Note
	return fmt.Sprintf("%s\n\n%s", noteParts[0], noteParts[1])
}

func generateDayOneRichText(entry *DaylioEntry, gen DayOneEntryUUIDGenerator) (string, error) {
	uuid, err := gen.GenerateUUID(entry)
	if err != nil {
		return "", err
	}
	rt := DayOneRichTextObjectData{
		Meta: DayOneRichTextObjectDataMetadata{
			Version: 1,
		},
		Contents: []DayOneRichTextObject{
			DayOneRichTextObject{
				Text: createRichTextNote(entry),
				Attributes: DayOneRichTextObjectAttributes{
					Line: DayOneRichTextLineObject{
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

func generateTagsFromDaylioActivities(entry *DaylioEntry) ([]string, error) {
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

func generateLocationFromDaylioActivities(entry *DaylioEntry) (DayOneEntryLocation, error) {
	if os.Getenv("NO_AUTO_HOME_LOCATION") != "" {
		return DayOneEntryLocation{}, nil
	}
	if os.Getenv("HOME_ADDRESS_JSON") == "" {
		log.Warn("Auto home location quirk is on but HOME_ADDRESS_JSON is empty")
		return DayOneEntryLocation{}, nil
	}
	var out DayOneEntryLocation
	if err := json.Unmarshal([]byte(os.Getenv("HOME_ADDRESS_JSON")), &out); err != nil {
		return DayOneEntryLocation{}, err
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

func createTimestamps(entry *DaylioEntry, g DayOneEntryModifiedTimestamper) (dayOneTimestamps, error) {
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
		Created:  created,
		Modified: modified,
	}, nil
}
