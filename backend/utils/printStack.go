package utils

import (
	"errors"
	"log/slog"
	"strings"
)

func PrintStack(errPtr *error) {
	err := *errPtr
	if err == nil {
		return
	}

	errorStack := []string{}
	for err != nil {
		errorStack = append(errorStack, err.Error())
		err = errors.Unwrap(err)
	}
	if len(errorStack) != 0 {
		slog.Error(strings.Join(errorStack, "\n"))
	}
}
