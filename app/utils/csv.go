package utils

import (
	"encoding/csv"
	"log"
	"os"
)

// ReadCSV reads the csv and returns a 2 dimensional array:
// First being the row.
// Second being the column.
func ReadCSV(csvPath string) ([][]string, error) {
	if csvPath == "" {
		log.Printf("csvPath: %s, is not found.", csvPath)
		return nil, nil
	}

	f, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}
