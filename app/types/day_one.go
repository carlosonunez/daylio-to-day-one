package types

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

// DayOneExport represents an export of a Day One journal (with entries) sans
// audio and video attachments.
type DayOneExport struct {
	Metadata DayOneMetadata `json:"metadata"`
	Entries  []DayOneEntry
}

// DayOneEntry is a representation of a journal entry.
type DayOneEntry struct {
	Starred             bool                `json:"starred"`
	Location            DayOneEntryLocation `json:"location"`
	CreationDeviceType  string              `json:"creationDeviceType"`
	CreationOSName      string              `json:"creationOSName"`
	CreationOSVersion   string              `json:"creationOSVersion"`
	CreationDate        DayOneDateTime      `json:"creationDate"`
	TimeZone            string              `json:"timeZone"`
	Tags                []string            `json:"tags"`
	Duration            int                 `json:"duration"`
	CreationDeviceModel string              `json:"creationDeviceModel"`
	// UUID is not a real UUID.
	UUID           string                 `json:"uuid"`
	IsAllDay       bool                   `json:"isAllDay"`
	Weather        map[string]interface{} `json:"weather"`
	ModifiedDate   DayOneDateTime         `json:"modifiedDate"`
	RichText       string                 `json:"richText"`
	Text           string                 `json:"text"`
	IsPinned       bool                   `json:"isPinned"`
	CreationDevice string                 `json:"creationDevice"`
}

type DayOneRichTextObjectData struct {
	Contents []DayOneRichTextObject           `json:"contents"`
	Meta     DayOneRichTextObjectDataMetadata `json:"meta"`
}

type DayOneRichTextObjectDataMetadata struct {
	Version           int                                   `json:"version"`
	SmallLinesRemoved bool                                  `json:"small-lines-removed"`
	Created           DayOneRichTextObjectCreatedProperties `json:"created"`
}

type DayOneRichTextObjectCreatedProperties struct {
	Platform string
	Version  int
}

type DayOneRichTextObject struct {
	Text       string                         `json:"text"`
	Attributes DayOneRichTextObjectAttributes `json:"attributes"`
}

type DayOneRichTextObjectAttributes struct {
	Line DayOneRichTextLineObject `json:"line"`
}

type DayOneRichTextLineObject struct {
	Header     int       `json:"header"`
	Identifier uuid.UUID `json:"identifier"`
}

// DayOneEntryLocation provides location data for a post.
type DayOneEntryLocation struct {
	Location           DayOneEntryLocationDetails `json:"location"`
	LocalityName       string                     `json:"localityName"`
	Country            string                     `json:"country"`
	TimeZoneName       string                     `json:"timeZoneName"`
	AdministrativeArea string                     `json:"administrativeArea"`
	Longitude          float32                    `json:"longitude"`
	PlaceName          string                     `json:"placeName"`
	Latitude           float32                    `json:"latitude"`
}

// DayOneEntryLocationDetails gives you coordinates and stuff.
type DayOneEntryLocationDetails struct {
	Region DayOneEntryLocationRegion `json:"region"`
}

// DayOneEntryLocationRegion gives you coords and stuff
type DayOneEntryLocationRegion struct {
	Radius int                       `json:"radius"`
	Center DayOneEntryLocationCoords `json:"center"`
}

// DayOneEntryLocationCoords are approx. coordinates for a Day One entry
// location.
type DayOneEntryLocationCoords struct {
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}

// DayOneDateTime is needed to work with Daylio's non-standard time format.
type DayOneDateTime time.Time

func (d *DayOneDateTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02T15:04Z", s)
	if err != nil {
		return err
	}
	*d = DayOneDateTime(t)
	return nil
}

func (d DayOneDateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Format("2006-01-02T15:04Z"))
}

func (d DayOneDateTime) Format(s string) string {
	return time.Time(d).Format(s)
}

// DayOneMetadata defines the version of the journal export.
type DayOneMetadata struct {
	Version string
}

// NewEmptyDayOneEntry generates an empty DayOne entry that looks like it was
// created from a MacBook using the Day One macOS app.
func NewEmptyDayOneEntry() *DayOneEntry {
	return &DayOneEntry{
		Starred:             false,
		CreationDeviceType:  "Laptop",
		CreationOSName:      "macOS",
		CreationOSVersion:   "14.1.2",
		CreationDeviceModel: "Mac14,2",
		TimeZone:            os.Getenv("TZ"),
		IsAllDay:            false,
		Weather:             map[string]interface{}{},
		IsPinned:            false,
		CreationDevice:      "MacBook"}
}

// NewDayOneExport generates a v1.0 DayOne export.
func NewDayOneExport(entries []DayOneEntry) *DayOneExport {
	return &DayOneExport{
		Metadata: DayOneMetadata{Version: "1.0"},
		Entries:  entries,
	}
}
