package daylio

const (
	DaylioMoodRad   = "rad"
	DaylioMoodGood  = "good"
	DaylioMoodMeh   = "meh"
	DaylioMoodBad   = "bad"
	DaylioMoodAwful = "awful"
)

// Entry is an entry in Daylio.
type Entry struct {
	FullDate       string   `csv:"full_date"`
	Date           string   `csv:"date"`
	Weekday        string   `csv:"weekday"`
	Time           string   `csv:"time"`
	Mood           string   `csv:"mood"`
	Activities     string   `csv:"activities"`
	ActivitiesList []string `csv:"activities_list,omitempty"`
	NoteTitle      string   `csv:"note_title"`
	Note           string   `csv:"note"`
}

// Backup is a full Daylio backup that can be used to restore Daylio from
// scratch.
type Backup struct {
	// Tags is a JSON representation of Daylio's tags database.
	Tags       []Tag      `json:"tags"`
	DayEntries []DayEntry `json:"dayEntries"`
}

// Tag is a tag within Daylio. There are more properties
// in the actual backup than exposed here; ID and Name are the only
// ones we care about.
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// DayEntry is a Daylio entry stored in a backup.
type DayEntry struct {
	Note     string `json:"note"`
	Title    string `json:"note_title"`
	TimeUNIX int64  `json:"datetime"`
	TagIDs   []int  `json:"tags"`
	Mood     int    `json:"mood"`
}
