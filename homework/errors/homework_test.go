package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	Errors []error
}

func (e *MultiError) Error() string {
	result := fmt.Sprintf("%d errors occured:\n", len(e.Errors))

	for _, e := range e.Errors {
		result += fmt.Sprintf("\t* %s", e.Error())
	}

	result += "\n"

	return result
}

func Append(err error, errs ...error) *MultiError {
	multi, ok := err.(*MultiError)
	if !ok {
		return &MultiError{
			Errors: errs,
		}
	}

	multi.Errors = append(multi.Errors, errs...)

	return multi
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
