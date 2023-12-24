package types

const (
	DaylioMoodRad   = "rad"
	DaylioMoodGood  = "good"
	DaylioMoodMeh   = "meh"
	DaylioMoodBad   = "bad"
	DaylioMoodAwful = "awful"
)

// DaylioEntry is an entry in Daylio.
type DaylioEntry struct {
	FullDate   string `csv:"full_date"`
	Date       string `csv:"date"`
	Weekday    string `csv:"weekday"`
	Time       string `csv:"time"`
	Mood       string `csv:"mood"`
	Activities string `csv:"activities"`
	NoteTitle  string `csv:"note_title"`
	Note       string `csv:"note"`
}
