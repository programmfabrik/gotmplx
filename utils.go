package main

import (
	"fmt"
	"io/ioutil"
	"os"
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

var (
	// stdinData serves as a cache for stdin data
	stdinData []byte
)

// readStdinData reads data from stdin and stores the values in the package variable "stdinData" if "stdinData" is nil, otherwise it returns "stdinData".
func readStdinData() ([]byte, error) {
	if stdinData != nil {
		return stdinData, nil
	}
	if os.Stdin != nil {
		stdinBts, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("unable read stdin data: %w", err)
		}
		stdinData = stdinBts
		return stdinData, nil
	}
	return nil, nil
}
