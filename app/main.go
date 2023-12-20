package exporter

import (
	"encoding/json"
	"fmt"
)

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

func convertToDayOne(entry *DaylioEntry, gList *dayOneGenerators) (*DayOneEntry, error) {
	return nil, nil
}
