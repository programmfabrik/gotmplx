package main

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// splitVarParam splits "param" into key, value pairs
func splitVarParam(param string) (string, string, error) {
	parts := strings.Split(param, "=")
	if len(parts) < 2 {
		return "", "", errors.Errorf("flag arguments should be `name=value`, given %s", param)
	}
	return parts[0], strings.Join(parts[1:], "="), nil
}

// sliceKeyValueToMap tries to generate key-value pairs from "strSlc". Each string must have a "key=value" syntax.
func sliceKeyValueToMap(strSlc []string) (map[string]interface{}, error) {
	retData := map[string]interface{}{}
	for _, data := range strSlc {
		key, value, err := splitVarParam(data)
		if err != nil {
			return nil, fmt.Errorf("unable to parse env var with key %q: %w", key, err)
		}
		retData[key] = value
	}
	return retData, nil
}
