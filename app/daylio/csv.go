package daylio

import (
	"os"

	csv "github.com/gocarina/gocsv"
)

func GetEntriesFromCSVFile(csvFile string) ([]Entry, error) {
	f, err := os.OpenFile(csvFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var entries []Entry
	if err := csv.UnmarshalFile(f, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}
