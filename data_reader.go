package main

import (
	"io/ioutil"
)

type sourceReader interface {
	// Unmarshal unmarshals bts to the desired format
	Unmarshal(bts []byte) (interface{}, error)
	// IsFile checks whether str contains the file ending
	IsFile(str string) bool
}

// readData checks if inputStrSlice provides key=value pairs and validates them using the following input techniques:
//   * file
//   * inline
//   * stdin
//
// The "tformat" parameter defines the parser we should use to extract the data.
func readData(inputStrSlice []string, sr sourceReader) (map[string]interface{}, error) {
	retData := make(map[string]interface{})
	for _, data := range inputStrSlice {
		key, value, err := splitVarParam(data)
		if err != nil {
			return nil, err
		}

		byteData := []byte{}
		if sr.IsFile(value) {
			// data from file
			byteData, err = ioutil.ReadFile(value)
			if err != nil {
				return nil, err
			}
		} else if value == "-" {
			// data from stdin
			stdin, err := readStdinData()
			if err != nil {
				return nil, err
			}
			byteData = stdin
		} else {
			// inline data
			byteData = []byte(value)
		}

		data, err := sr.Unmarshal(byteData)
		if err != nil {
			return nil, err
		}

		retData[key] = data
	}
	return retData, nil
}
