package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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
	inputCh := os.Stdin
	retData := make(map[string]interface{})
	for _, data := range inputStrSlice {
		key, value, err := splitVarParam(data)
		if err != nil {
			return nil, err
		}

		byteData := []byte{}
		if sr.IsFile(value) {
			byteData, err = ioutil.ReadFile(value)
			if err != nil {
				return nil, err
			}
		} else if strings.ContainsRune(value, '-') {
			// data from stdin
			byteData, err = ioutil.ReadAll(inputCh)
			if err != nil {
				return nil, fmt.Errorf("unable read stdin data: %w", err)
			}
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

// stringSliceToMap extracts key=value pairs from the string slice and writes them as key=value pair to the map
func stringSliceToMap(strs []string) (map[string]string, error) {
	if len(strs) < 1 {
		return nil, errors.New("need at least one key=value pair")
	}
	rslt := map[string]string{}
	for _, idxValue := range strs {
		key, value, err := splitVarParam(idxValue)
		if err != nil {
			return nil, err
		}
		rslt[key] = value
	}
	return rslt, nil
}
