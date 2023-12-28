package daylio

import (
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockDirEntryInfo struct {
	name  string
	mtime time.Time
}

func (m *mockDirEntryInfo) Name() string {
	return m.name
}

func (m *mockDirEntryInfo) Mode() fs.FileMode {
	return fs.FileMode(0777)
}

func (m *mockDirEntryInfo) ModTime() time.Time {
	return m.mtime
}

func (m *mockDirEntryInfo) IsDir() bool {
	return false
}

func (m *mockDirEntryInfo) Sys() any {
	return nil
}

func (m *mockDirEntryInfo) Size() int64 {
	return 0
}

func newMockDirEntryInfo(name string, mtime time.Time) *mockDirEntryInfo {
	return &mockDirEntryInfo{
		name:  name,
		mtime: mtime,
	}
}

type mockDirEntry struct {
	name string
	info fs.FileInfo
}

func (m *mockDirEntry) Name() string {
	return m.name
}

func (m *mockDirEntry) IsDir() bool {
	return false
}

func (m *mockDirEntry) Type() fs.FileMode {
	return fs.FileMode(0777)
}

func (m *mockDirEntry) Info() (fs.FileInfo, error) {
	return m.info, nil
}

func newMockDirEntry(name string, mtime time.Time) *mockDirEntry {
	return &mockDirEntry{
		name: name,
		info: newMockDirEntryInfo(name, mtime),
	}
}

type mockEntryInfo struct {
	name  string
	mtime time.Time
}

func newMockEntryInfo(name string, modTime time.Time) mockEntryInfo {
	return mockEntryInfo{name: name, mtime: modTime}
}

func newMockDirEntryList(infos []mockEntryInfo) []fs.DirEntry {
	l := []fs.DirEntry{}
	for _, info := range infos {
		l = append(l, newMockDirEntry(info.name, info.mtime))
	}
	return l
}

func mustParseDaylioEntryTime(ts string) time.Time {
	t, err := time.Parse("2006-01-02", ts)
	if err != nil {
		panic(err)
	}
	return t
}

type mockTraverser struct{}

func (t *mockTraverser) Dir() string {
	return "/Users/foobar/Library/Mobile Documents/com~apple~CloudDocs/Downloads"
}

func (t *mockTraverser) ListBackups() ([]fs.DirEntry, error) {
	return newMockDirEntryList([]mockEntryInfo{
		newMockEntryInfo("ios_backup_2023_12_10.daylio", mustParseDaylioEntryTime("2023-12-10")),
		newMockEntryInfo("ios_backup_2023_12_24.daylio", mustParseDaylioEntryTime("2023-12-24")),
		newMockEntryInfo("ios_backup_2023_12_06.daylio", mustParseDaylioEntryTime("2023-12-06")),
	}), nil
}

type mockEmptyTraverser struct{}

func (t *mockEmptyTraverser) Dir() string {
	return "/Users/foobar/Library/Mobile Documents/com~apple~CloudDocs/Downloads"
}

func (t *mockEmptyTraverser) ListBackups() ([]fs.DirEntry, error) {
	return []fs.DirEntry{}, nil
}

func TestGettingLatestBackupLocation_iCloud(t *testing.T) {
	want := "/Users/foobar/Library/Mobile Documents/com~apple~CloudDocs/Downloads/ios_backup_2023_12_24.daylio"
	got, err := resolveDaylioBackupLocationMacOS(&mockTraverser{})
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestGettingLatestBackupLocationNoFiles_iCloud(t *testing.T) {
	_, err := resolveDaylioBackupLocationMacOS(&mockEmptyTraverser{})
	assert.Error(t, err, "No backups found in '/Users/foobar/Library/Mobile Documents/com~apple~CloudDocs/Downloads'")
}

func TestReadingBackupJSON(t *testing.T) {
	json, err := os.ReadFile("./fixtures/daylio.json")
	require.NoError(t, err)
	want := Backup{
		Tags: []Tag{
			{ID: 1, Name: "activity 1"},
			{ID: 2, Name: "activity 2"},
			{ID: 3, Name: "activity 3"},
		},
		DayEntries: []DayEntry{
			{
				Note:     "note text 1",
				Title:    "note title",
				TimeUNIX: 1702800000,
				TagIDs:   []int{1, 2, 3},
				Mood:     1,
			},
			{
				Note:     "note text 2",
				Title:    "",
				TimeUNIX: 1702713600,
				TagIDs:   []int{1, 2, 3},
				Mood:     2,
			},
			{
				Note:     "note text 3",
				Title:    "",
				TimeUNIX: 1702627200,
				TagIDs:   []int{1},
				Mood:     3,
			},
		},
	}
	got, err := backupFromJSON(json)
	assert.NoError(t, err)
	assert.Equal(t, &want, got)
}

func TestDecodingBackupJSON(t *testing.T) {
	want, err := os.ReadFile("./fixtures/daylio.json")
	require.NoError(t, err)
	data, err := os.ReadFile("./fixtures/daylio.json.encoded")
	require.NoError(t, err)
	got, err := decodeDaylioJSON(data)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestDayEntryToSimpleEntry(t *testing.T) {
	entry := DayEntry{
		Note:     "note text 1",
		Title:    "note title",
		TimeUNIX: 1702800000,
		TagIDs:   []int{1, 2, 3},
		Mood:     1,
	}
	tags := []Tag{
		{ID: 1, Name: "activity 1"},
		{ID: 2, Name: "activity 2"},
		{ID: 3, Name: "activity 3"},
	}
	want := Entry{
		FullDate:       "2023-12-17",
		Date:           "Dec 17",
		Weekday:        "Sunday",
		Time:           "08:00",
		Mood:           "rad",
		ActivitiesList: []string{"activity 1", "activity 2", "activity 3"},
		NoteTitle:      "note title",
		Note:           "note text 1",
	}
	got, err := dayEntryToEntry(&entry, tags)
	assert.NoError(t, err)
	assert.Equal(t, want, *got)
}
