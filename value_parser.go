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
	formatJSON format = "json"
	formatCSV  format = "csv"
	formatVar  format = "var"
)

type valueParser struct {
	FormatType format
}

// Unmarshal unmarshals the data byte slice into the desired format.
//   formatJSON: returns a map[string]interface{}
//   formatCSV: returns a []map[string]interface{}
func (v valueParser) Unmarshal(data []byte) (interface{}, error) {
	switch v.FormatType {
	case formatCSV:
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
	case formatJSON:
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
