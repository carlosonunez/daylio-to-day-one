package exporter

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUUIDGenerator struct{}

func (g *mockUUIDGenerator) GenerateUUID(entry *DaylioEntry) (uuid.UUID, error) {
	switch entry.Note {
	case "note text 1":
		return uuid.FromBytes([]byte("5D73E4F9-491A-4DB4-BE24-D89CC8C52636"))
	case "note text 2":
		return uuid.FromBytes([]byte("D5265940-000C-465D-8FFB-602375CEA7AE"))
	case "note text 3":
		return uuid.FromBytes([]byte("E25B024A-708E-4143-B1B6-4F1AD8BF76E7"))
	default:
		return uuid.UUID{}, fmt.Errorf("Invalid test Daylio Note: %s", entry.Note)
	}
}

type mockIDGenerator struct{}

func (g *mockIDGenerator) CreateID(entry *DaylioEntry) string {
	switch entry.Note {
	case "note text 1":
		return "DFFVSK3JENO8ZHKQTDLOS8AEWFS7IPS3O"
	case "note text 2":
		return "CQD40MRYMUIT02FYYI2FQTU81DELADNQZ"
	case "note text 3":
		return "UXQQ9CPNEYZ13DI59FPVQ5YQWXZOQBE7I"
	default:
		return "INVALID-POST"
	}
}

type mockTimestamper struct{}

func (g *mockTimestamper) CreateModifiedTime(entry *DaylioEntry) (time.Time, error) {
	return time.Parse("2006-01-02T15:04Z", "2023-12-20T12:13Z")
}

func TestConvertToDayOneEntryComplete(t *testing.T) {
	var export DayOneExport
	wantJSON, err := os.ReadFile("./fixtures/dayone.json")
	require.NoError(t, err)
	err = json.Unmarshal(wantJSON, &export)
	require.NoError(t, err)
	want := export.Entries[0]
	gGen := mockUUIDGenerator{}
	iGen := mockIDGenerator{}
	tGen := mockTimestamper{}
	generators := dayOneGenerators{
		UUIDGenerator: &gGen,
		IDGenerator:   &iGen,
		Timestamper:   &tGen,
	}
	entry := DaylioEntry{
		FullDate:   "2023-12-17",
		Date:       "Dec 17",
		Weekday:    "Sunday",
		Time:       "08:00",
		Mood:       DaylioMoodRad,
		Activities: "activity 1 | activity 2 | activity 3",
		NoteTitle:  "Title",
		Note:       "note text 1",
	}
	got, err := convertToDayOne(&entry, &generators)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
