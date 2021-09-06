package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/programmfabrik/go-csvx"
)

var (
	ErrUnsupportedFormatType = errors.New("unsupported format type")
)

type format string

const (
	FormatJSON format = "json"
	FormatCSV  format = "csv"
	FormatVar  format = "var"
)

type ValueParser struct {
	FormatType format
}

// Unmarshal unmarshals the data byte slice into the desired format.
//   FormatJSON: returns a map[string]interface{}
//   FormatCSV: returns a []map[string]interface{}
func (v ValueParser) Unmarshal(data []byte) (interface{}, error) {
	switch v.FormatType {
	case FormatCSV:
		csvp := csvx.CSVParser{
			Comma:            ',',
			Comment:          '#',
			TrimLeadingSpace: true,
			SkipEmptyColumns: true,
		}
		csvData, err := csvp.Untyped(data)
		if err != nil {
			return nil, fmt.Errorf("unable to parse bytes into CSV format: %w", err)
		}
		return csvData, nil
	case FormatJSON:
		jsonData := map[string]interface{}{}
		err := json.Unmarshal(data, &jsonData)
		if err != nil {
			return nil, fmt.Errorf("unable to parse bytes into JSON format: %w", err)
		}
		return jsonData, nil
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedFormatType, v.FormatType)
	}
}
