package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/pkg/errors"
)

var (
	ErrUnableToIncrementStdinRefCounter = errors.New("unable to increment stdin reference counter. Counter already increased")
)

type cliInt struct {
	stdinRefCounter uint8
	inputCh         io.Reader
}

// newCliIntWithStdinInputCh initializes the inputCh with the os.Stdin channel
func newCliIntWithStdinInputCh() *cliInt {
	return &cliInt{
		inputCh: os.Stdin,
	}
}

// incrementStdinRef checks if the stdinRefCounter is set to 0 and increments the counter by one.
// If stdinRefCounter is at this point > 1 an error is returned
func (ci *cliInt) incrementStdinRef() error {
	if ci.stdinRefCounter > 0 {
		return ErrUnableToIncrementStdinRefCounter
	}
	ci.stdinRefCounter++
	return nil
}

// inputIndicatesStdinData checks if input matches -
func (ci *cliInt) inputIndicatesStdinData(input string) bool {
	if input != "-" {
		return false
	}
	return true
}

// extractData checks if inputStrSlice provides key=value pairs and validates them using the following input techniques:
//   * file
//   * inline
//   * stdin
//
// The "tformat" parameter defines the parser we should use to extract the data.
func (ci *cliInt) extractData(inputStrSlice []string, tformat format) (map[string]interface{}, error) {
	retData := make(map[string]interface{})
	for _, data := range inputStrSlice {
		key, value, err := splitVarParam(data)
		if err != nil {
			return nil, err
		}

		// check if format is of type var
		if tformat == FormatVar {
			retData[key] = value
			continue
		}

		bts := []byte{}

		matched, err := regexp.MatchString(".(json|csv)$", value)
		if err != nil {
			return nil, fmt.Errorf("regex file extension failed with error: %w", err)
		}

		if ci.inputIndicatesStdinData(value) {
			// data from stdin
			err = ci.incrementStdinRef()
			if err != nil {
				return nil, fmt.Errorf("%w: for new key %q", err, key)
			}

			bts, err = ioutil.ReadAll(ci.inputCh)
			if err != nil {
				return nil, fmt.Errorf("unable read stdin data: %w", err)
			}
		} else if matched {
			// data from file
			bts, err = ioutil.ReadFile(value)
			if err != nil {
				return nil, fmt.Errorf("unable to read file %q: %w", value, err)
			}
		} else {
			// inline data
			bts = []byte(value)
		}

		parser := ValueParser{
			FormatType: tformat,
		}

		pData, err := parser.Unmarshal(bts)
		if err != nil {
			return nil, fmt.Errorf("%w failed for\n%+v\n", err, string(bts))
		}

		retData[key] = pData
	}
	return retData, nil
}
