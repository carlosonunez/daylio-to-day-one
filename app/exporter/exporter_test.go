package exporter

import (
	"bytes"
	"encoding/json"
	"exporter/daylio"
	"exporter/types"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	csv "github.com/gocarina/gocsv"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	FirstMockNoteUUID  = "5D73E4F9-491A-4DB4-BE24-D89CC8C52636"
	SecondMockNoteUUID = "D5265940-000C-465D-8FFB-602375CEA7AE"
	ThirdMockNoteUUID  = "E25B024A-708E-4143-B1B6-4F1AD8BF76E7"
	FirstMockNoteID    = "DFFVSK3JENO8ZHKQTDLOS8AEWFS7IPS3O"
	SecondMockNoteID   = "CQD40MRYMUIT02FYYI2FQTU81DELADNQZ"
	ThirdMockNoteID    = "UXQQ9CPNEYZ13DI59FPVQ5YQWXZOQBE7I"
)

type mockUUIDGenerator struct{ uuid string }

func (g *mockUUIDGenerator) GenerateUUID() (uuid.UUID, error) {
	return uuid.Parse(g.uuid)
}

type mockIDGenerator struct{ id string }

func (g *mockIDGenerator) CreateID() string {
	return g.id
}

type mockTimestamper struct{}

func (g *mockTimestamper) CreateModifiedTime() (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05Z", "2023-12-20T12:13:00Z")
}

func newMockUUIDGenerator(t *testing.T, entry *daylio.Entry) *mockUUIDGenerator {
	switch entry.Note {
	case "note text 1":
		return &mockUUIDGenerator{uuid: FirstMockNoteUUID}
	case "note text 2":
		return &mockUUIDGenerator{uuid: SecondMockNoteUUID}
	case "note text 3":
		return &mockUUIDGenerator{uuid: ThirdMockNoteUUID}
	default:
		return &mockUUIDGenerator{uuid: FirstMockNoteUUID}
	}
}

func newMockIDGenerator(t *testing.T, entry *daylio.Entry) *mockIDGenerator {
	switch entry.Note {
	case "note text 1":
		return &mockIDGenerator{id: FirstMockNoteID}
	case "note text 2":
		return &mockIDGenerator{id: SecondMockNoteID}
	case "note text 3":
		return &mockIDGenerator{id: ThirdMockNoteID}
	default:
		return &mockIDGenerator{id: FirstMockNoteID}
	}
}

func mustGetZuluTime(s string) time.Time {
	t, err := time.Parse("2006-01-02T15:04:05Z", s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestWritingDayOneExportsToDisk(t *testing.T) {
	var buf bytes.Buffer
	want := types.DayOneExport{
		Metadata: types.DayOneMetadata{Version: "1.0"},
		Entries: []types.DayOneEntry{{
			Text: "hello",
			Tags: []string{"tag 1", "tag 2", "tag 3"},
		},
		},
	}
	err := writeDayOneExport(&buf, &want)
	assert.NoError(t, err)
	var got types.DayOneExport
	err = json.Unmarshal(buf.Bytes(), &got)
	require.NoError(t, err)
	assert.Equal(t, want.Metadata.Version, got.Metadata.Version)
	assert.Equal(t, want.Entries[0].Text, got.Entries[0].Text)
}

func TestGenerateDayOneRichText(t *testing.T) {
	entry := daylio.Entry{
		NoteTitle: "Title",
		Note:      "note text 1",
	}
	uGen := newMockUUIDGenerator(t, &entry)
	got, err := generateDayOneRichText(&entry, uGen)
	assert.NoError(t, err)
	assert.Contains(t, got, fmt.Sprintf(`"identifier":"%s"`, strings.ToLower(FirstMockNoteUUID)))
	assert.Contains(t, got, fmt.Sprintf(`"text":"%s\n\n%s"`, entry.NoteTitle, entry.Note))
}

func TestGenerateDayOneRichTextWithoutTitle(t *testing.T) {
	entry := daylio.Entry{
		Note: "note text 1",
	}
	uGen := newMockUUIDGenerator(t, &entry)
	got, err := generateDayOneRichText(&entry, uGen)
	assert.NoError(t, err)
	assert.Contains(t, got, fmt.Sprintf(`"identifier":"%s"`, strings.ToLower(FirstMockNoteUUID)))
	assert.Contains(t, got, fmt.Sprintf(`"text":"Note\n\n%s"`, entry.Note))
}

func TestCreateTimestamps(t *testing.T) {
	entry := daylio.Entry{
		FullDate: "2023-12-17",
		Time:     "08:00",
	}
	want := dayOneTimestamps{
		Created:  types.DayOneDateTime(mustGetZuluTime("2023-12-17T08:00:00Z")),
		Modified: types.DayOneDateTime(mustGetZuluTime("2023-12-20T12:13:00Z")),
	}
	g := mockTimestamper{}
	got, err := createTimestamps(&entry, &g)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestConvertToDayOneSingle(t *testing.T) {
	t.Setenv("TZ", "America/Chicago")
	var export types.DayOneExport
	wantJSON, err := os.ReadFile("./fixtures/dayone.json")
	require.NoError(t, err)
	err = json.Unmarshal(wantJSON, &export)
	require.NoError(t, err)
	want := []types.DayOneEntry{export.Entries[0]}
	gGen := newMockUUIDGenerator(t, &daylio.Entry{})
	iGen := newMockIDGenerator(t, &daylio.Entry{})
	tGen := mockTimestamper{}
	generators := types.DayOneGenerators{
		UUIDGenerator: gGen,
		IDGenerator:   iGen,
		Timestamper:   &tGen,
	}
	var entries []daylio.Entry
	csvRaw := `full_date,date,weekday,time,mood,activities,note_title,note
2023-12-17,Dec 17,Sunday,08:00,good,activity 1 | activity 2 | activity 3,note title,note text 1`
	err = csv.UnmarshalString(csvRaw, &entries)
	require.NoError(t, err)
	got, err := convertToDayOneEntries(entries, generators)
	// NOTE: Ignore testing RichText, as this is covered by another test.  This
	// will always fail due to the keys in the underlying map being inserted in
	// random order.
	require.Len(t, got, 1)
	want[0].RichText = ""
	got[0].RichText = ""
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestConvertDaylioEntriesToDayOne(t *testing.T) {
	t.Setenv("TZ", "America/Chicago")
	var entries []daylio.Entry
	mockEntries, err := os.OpenFile("./fixtures/daylio.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	require.NoError(t, err)
	defer mockEntries.Close()
	err = csv.UnmarshalFile(mockEntries, &entries)
	require.NoError(t, err)
	var export types.DayOneExport
	wantJSON, err := os.ReadFile("./fixtures/dayone.json")
	require.NoError(t, err)
	err = json.Unmarshal(wantJSON, &export)
	require.NoError(t, err)
	want := export.Entries
	gGen := newMockUUIDGenerator(t, &daylio.Entry{})
	iGen := newMockIDGenerator(t, &daylio.Entry{})
	tGen := mockTimestamper{}
	generators := types.DayOneGenerators{
		UUIDGenerator: gGen,
		IDGenerator:   iGen,
		Timestamper:   &tGen,
	}
	got, err := convertToDayOneEntries(entries, generators)
	// NOTE: Ignore testing RichText and UUIOD, as this is covered by another test.
	for idx := 0; idx < len(want); idx++ {
		want[idx].RichText = ""
		want[idx].UUID = ""
	}
	for idx := 0; idx < len(got); idx++ {
		got[idx].RichText = ""
		got[idx].UUID = ""
	}
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestExportLocationQuirk(t *testing.T) {
	locJSON, err := os.ReadFile("./fixtures/home_location.json")
	require.NoError(t, err)
	t.Setenv("HOME_ADDRESS_JSON", string(locJSON))
	var want types.DayOneEntryLocation
	err = json.Unmarshal(locJSON, &want)
	require.NoError(t, err)
	got, err := generateLocationFromDaylioActivities([]string{
		"activity 1",
		"home",
		"activity 2",
	})
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
