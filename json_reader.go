package main

import (
	"encoding/json"
	"strings"
)

// jsonReader implements the sourceReader interface.
type jsonReader struct{}

// Unmarshal decodes the byte slice into a go readable format.
// Unmarshal implements the sourceReader interface.
func (*jsonReader) Unmarshal(bts []byte) (interface{}, error) {
	var jsonData interface{}
	err := json.Unmarshal(bts, &jsonData)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

// IsFile checks whether str ends with .json.
// IsFile implements the sourceReader interface.
func (*jsonReader) IsFile(str string) bool {
	return strings.HasSuffix(str, ".json")
}
