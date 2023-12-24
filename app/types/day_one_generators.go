package types

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// DayOneEntryUUIDGenerator produces Day One richText UUIDs.
type DayOneEntryUUIDGenerator interface {
	// GenerateUUID makes the UUID
	GenerateUUID(entry *DaylioEntry) (uuid.UUID, error)
}

// DayOneIDGenerator generates a Day One entry ID. It needs to be 33 chars long.
type DayOneIDGenerator interface {
	// GenerateID creates an ID.
	CreateID(entry *DaylioEntry) string
}

// Timestamper produces modifiedOn timestamps
type DayOneEntryModifiedTimestamper interface {
	CreateModifiedTime(entry *DaylioEntry) (time.Time, error)
}

// DayOneGenerators is used to store references to ID and timestamp generators
// used to create Day One entries from Daylio entries.
type DayOneGenerators struct {
	UUIDGenerator DayOneEntryUUIDGenerator
	IDGenerator   DayOneIDGenerator
	Timestamper   DayOneEntryModifiedTimestamper
}

type DefaultDayOneEntryUUIDGenerator struct{}

func (g *DefaultDayOneEntryUUIDGenerator) GenerateUUID(entry *DaylioEntry) (uuid.UUID, error) {
	return uuid.New(), nil
}

type DefaultDayOneIDGenerator struct{}

func (g *DefaultDayOneIDGenerator) CreateID(entry *DaylioEntry) string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

type DefaultDayOneEntryModifiedTimestamper struct{}

func (g *DefaultDayOneEntryModifiedTimestamper) CreateModifiedTime(entry *DaylioEntry) (time.Time, error) {
	return time.Now(), nil
}

func DefaultDayOneGenerators() DayOneGenerators {
	return DayOneGenerators{
		UUIDGenerator: &DefaultDayOneEntryUUIDGenerator{},
		IDGenerator:   &DefaultDayOneIDGenerator{},
		Timestamper:   &DefaultDayOneEntryModifiedTimestamper{},
	}
}
