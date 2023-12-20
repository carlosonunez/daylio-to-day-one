package exporter

import (
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

type dayOneGenerators struct {
	UUIDGenerator DayOneEntryUUIDGenerator
	IDGenerator   DayOneIDGenerator
	Timestamper   DayOneEntryModifiedTimestamper
}
