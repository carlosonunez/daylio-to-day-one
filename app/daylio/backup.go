package daylio

import (
	"archive/zip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var DaylioMoodIDs = map[int]string{
	1: "rad",
	2: "good",
	3: "ok",
	4: "bad",
	5: "awful",
}

// BackupFileTraverser lists Daylio backups.
type BackupFileTraverser interface {
	// Dir() provides the name of the directory being traversed.
	Dir() string
	// ListBackups lists Daylio backups.
	ListBackups() ([]fs.DirEntry, error)
}

type defaultBFT struct{}

func (t *defaultBFT) Dir() string {
	return filepath.Join(os.Getenv("HOME"), "Library", "MobileDocuments", "com~apple~CloudDocs", "Downloads")
}

func (t *defaultBFT) ListBackups() ([]fs.DirEntry, error) {
	log.Debugf("Searching for Daylio backups here: %s", t.Dir())
	fl, err := os.ReadDir(t.Dir())
	if err != nil {
		return nil, err
	}
	var daylioFiles []fs.DirEntry
	for _, f := range fl {
		if strings.Contains(f.Name(), "ios_backup") {
			log.Debugf("Found backup file: %s", f.Name())
			daylioFiles = append(daylioFiles, f)
		}
	}
	return daylioFiles, nil
}

// GetEntriesFromBackupFile retrieves entries from a backup file.
func GetEntriesFromBackupFile(providedFile string) ([]Entry, error) {
	var fpath string
	var err error
	if providedFile == "" {
		fpath, err = resolveDaylioBackupLocationMacOS(&defaultBFT{})
		if err != nil {
			return nil, err
		}
	} else {
		fpath = providedFile
	}
	json, err := extractJSONFromDaylioBackupFile(fpath)
	if err != nil {
		return nil, err
	}
	backup, err := backupFromJSON(json)
	if err != nil {
		return nil, err
	}
	entries, err := simpleEntriesFromBackup(backup)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func simpleEntriesFromBackup(b *Backup) ([]Entry, error) {
	el := []Entry{}
	for _, d := range b.DayEntries {
		e, err := dayEntryToEntry(&d, b.Tags)
		if err != nil {
			return nil, err
		}
		el = append(el, *e)
	}
	return el, nil
}

func extractJSONFromDaylioBackupFile(fpath string) ([]byte, error) {
	reader, err := zip.OpenReader(fpath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	for _, f := range reader.File {
		if f.FileHeader.Name == "backup.daylio" {
			jsonEnc, err := getEncodedDaylioJSON(f)
			if err != nil {
				return nil, err
			}
			json, err := decodeDaylioJSON(jsonEnc)
			if err != nil {
				return nil, err
			}
			return json, nil
		}
	}
	return nil, fmt.Errorf("No Daylio backup JSONs found in file: %s", fpath)
}

func decodeDaylioJSON(b []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(strings.ReplaceAll(string(b), "\r\n", ""))
}

func getEncodedDaylioJSON(f *zip.File) ([]byte, error) {
	fReader, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer fReader.Close()
	encJSON, err := io.ReadAll(fReader)
	if err != nil {
		return nil, err
	}
	return encJSON, nil
}

func backupFromJSON(b []byte) (*Backup, error) {
	var backup Backup
	if err := json.Unmarshal(b, &backup); err != nil {
		return nil, err
	}
	return &backup, nil
}

func resolveDaylioBackupLocationMacOS(t BackupFileTraverser) (string, error) {
	backupList, err := t.ListBackups()
	if err != nil {
		return "", err
	}
	if len(backupList) == 0 {
		return "", fmt.Errorf("No backups found in '%s'", t.Dir())
	}
	sort.Slice(backupList, func(i, j int) bool {
		iInfo, err := backupList[i].Info()
		if err != nil {
			panic(err)
		}
		jInfo, err := backupList[j].Info()
		if err != nil {
			panic(err)
		}
		return iInfo.ModTime().Unix() > jInfo.ModTime().Unix()
	})
	return filepath.Join(t.Dir(), backupList[0].Name()), nil
}

func dayEntryToEntry(d *DayEntry, tags []Tag) (*Entry, error) {
	activities, err := exportTagsFromIDs(d.TagIDs, tags)
	if err != nil {
		return nil, err
	}
	mood, err := resolveMood(d.Mood)
	if err != nil {
		return nil, err
	}
	eTime := time.Unix(d.TimeUNIX, 0).UTC()
	return &Entry{
		FullDate:       eTime.Format("2006-01-02"),
		Date:           eTime.Format("Jan 02"),
		Weekday:        eTime.Format("Monday"),
		Time:           eTime.Format("15:04"),
		Mood:           mood,
		ActivitiesList: activities,
		NoteTitle:      d.Title,
		Note:           d.Note,
	}, nil
}

func resolveMood(mID int) (string, error) {
	mName, ok := DaylioMoodIDs[mID]
	if ok {
		return mName, nil
	}
	return "", fmt.Errorf("Not a valid Daylio mood ID: %d", mID)
}
