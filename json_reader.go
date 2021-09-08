package main

import (
	"encoding/json"
	"regexp"
)

// jsonReader implements the sourceReader interface.
type jsonReader struct{}

// Unmarshal decodes the byte slice into a go readable format.
// Unmarshal implements the sourceReader interface.
func (*jsonReader) Unmarshal(bts []byte) (interface{}, error) {
	jsonData := map[string]interface{}{}
	err := json.Unmarshal(bts, &jsonData)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

// IsFile checks whether str ends with .json.
// IsFile implements the sourceReader interface.
func (*jsonReader) IsFile(str string) (bool, error) {
	return regexp.MatchString("(.json)$", str)
}
