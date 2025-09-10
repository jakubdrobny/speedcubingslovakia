package utils

import (
	"errors"
	"log/slog"
	"strings"
)

func PrintStack(err error) {
	errorStack := []string{}
	for err != nil {
		errorStack = append(errorStack, err.Error())
		err = errors.Unwrap(err)
	}
	if len(errorStack) != 0 {
		slog.Error(strings.Join(errorStack, "\n"))
	}
}
