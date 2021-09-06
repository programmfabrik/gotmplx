package main

import (
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
