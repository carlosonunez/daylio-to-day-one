package daylio

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExportTagsNoQuirks(t *testing.T) {
	tags := []Tag{
		{ID: 0, Name: "activity 1"},
		{ID: 1, Name: "activity 2"},
		{ID: 2, Name: "activity 3"},
	}
	want := []string{"activity 1", "activity 2", "activity 3"}
	got, err := exportTagsFromIDs([]int{0, 1, 2}, tags)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestExportTagsAloneTimeQuirk(t *testing.T) {
	tags := []Tag{
		{ID: 0, Name: "No"},
		{ID: 1, Name: "A Little Bit"},
		{ID: 2, Name: "Yes!"},
	}
	want := []string{
		"alone score: 0",
		"alone score: 1",
		"alone score: 2",
	}
	got, err := exportTagsFromIDs([]int{0, 1, 2}, tags)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestExportTagsAloneTimeQuirkWhenDisabled(t *testing.T) {
	t.Setenv("NO_ALONE_TIME_SCORING", "anything")
	tags := []Tag{
		{ID: 0, Name: "No"},
		{ID: 1, Name: "A Little Bit"},
		{ID: 2, Name: "Yes!"},
	}
	want := []string{
		"No",
		"A Little Bit",
		"Yes!",
	}
	got, err := exportTagsFromIDs([]int{0, 1, 2}, tags)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
