package util

import (
	"errors"
)

func Check(errors ...error) {
	for _, err := range errors {
		if err != nil {
			panic(err)
		}
	}
}

func AppendLineNumToErr(err error) error {
	// skip = 2, count myself in.
	newMsg := GetLineNumSkip(2) + " " + err.Error()
	return errors.New(newMsg)
}
