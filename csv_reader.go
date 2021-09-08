package main

import (
	"fmt"
	"regexp"

	"github.com/programmfabrik/go-csvx"
)

// csvReader implements the sourceReader interface.
type csvReader struct{}

// Unmarshal decodes the byte slice into a go readable format.
// Unmarshal implements the sourceReader interface.
func (*csvReader) Unmarshal(bts []byte) (interface{}, error) {
	csvp := csvx.CSVParser{
		Comma:            ',',
		Comment:          '#',
		TrimLeadingSpace: true,
		SkipEmptyColumns: true,
	}
	csvData, err := csvp.Untyped(bts)
	if err != nil {
		return nil, fmt.Errorf("unable to parse bytes into CSV format: %w", err)
	}
	return csvData, nil
}

// IsFile checks whether str ends with .csv.
// IsFile implements the sourceReader interface.
func (*csvReader) IsFile(str string) (bool, error) {
	return regexp.MatchString("(.csv)$", str)
}
