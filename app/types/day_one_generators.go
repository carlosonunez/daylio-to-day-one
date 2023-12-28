package types

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// DayOneEntryUUIDGenerator produces Day One richText UUIDs.
type DayOneEntryUUIDGenerator interface {
	// GenerateUUID makes the UUID
	GenerateUUID() (uuid.UUID, error)
}

// DayOneIDGenerator generates a Day One entry ID. It needs to be 33 chars long.
type DayOneIDGenerator interface {
	// GenerateID creates an ID.
	CreateID() string
}

// Timestamper produces modifiedOn timestamps
type DayOneEntryModifiedTimestamper interface {
	CreateModifiedTime() (time.Time, error)
}

// DayOneGenerators is used to store references to ID and timestamp generators
// used to create Day One entries from Daylio entries.
type DayOneGenerators struct {
	UUIDGenerator DayOneEntryUUIDGenerator
	IDGenerator   DayOneIDGenerator
	Timestamper   DayOneEntryModifiedTimestamper
}

type DefaultDayOneEntryUUIDGenerator struct{}

func (g *DefaultDayOneEntryUUIDGenerator) GenerateUUID() (uuid.UUID, error) {
	return uuid.New(), nil
}

type DefaultDayOneIDGenerator struct{}

func (g *DefaultDayOneIDGenerator) CreateID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

type DefaultDayOneEntryModifiedTimestamper struct{}

func (g *DefaultDayOneEntryModifiedTimestamper) CreateModifiedTime() (time.Time, error) {
	return time.Now(), nil
}

func DefaultDayOneGenerators() DayOneGenerators {
	return DayOneGenerators{
		UUIDGenerator: &DefaultDayOneEntryUUIDGenerator{},
		IDGenerator:   &DefaultDayOneIDGenerator{},
		Timestamper:   &DefaultDayOneEntryModifiedTimestamper{},
	}
}
