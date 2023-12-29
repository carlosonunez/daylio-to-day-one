package daylio

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// we're assuming that no Daylio entries have tag IDs that aren't in the backup.
func exportTagsFromIDs(ids []int, tags []Tag) ([]string, error) {
	tagNames := []string{}
	tagHT := map[int]string{}
	for _, tag := range tags {
		tagHT[tag.ID] = tag.Name
	}
	log.Tracef("tags: %+v", tagHT)
	for _, id := range ids {
		log.Tracef("looking for tag id: '%d'", id)
		tagName, ok := tagHT[id]
		if !ok {
			return []string{}, fmt.Errorf("tag ID not in Daylio backup: %d", id)
		}
		if score := generateAloneTimeScore(tagName); score != "" {
			tagName = score
		}
		tagNames = append(tagNames, tagName)
	}
	return tagNames, nil
}

func generateAloneTimeScore(s string) string {
	if os.Getenv("NO_ALONE_TIME_SCORING") != "" {
		return ""
	}
	var score int
	switch strings.ToLower(s) {
	case "no":
		score = 0
	case "a little bit":
		score = 1
	case "yes!":
		score = 2
	default:
		return ""
	}
	return fmt.Sprintf("alone score: %d", score)
}
