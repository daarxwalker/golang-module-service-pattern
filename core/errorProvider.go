package core

import (
	"errors"
	"fmt"
)

type ErrorProvider interface {
	getError() error
	getMessage() string
	getFormatted() string
}

type errorProvider struct {
	value   error
	message string
}

func createError(value error, message string) ErrorProvider {
	return errorProvider{
		value,
		message,
	}
}

func createNewError(value string) ErrorProvider {
	return errorProvider{
		value: errors.New(value),
	}
}

func (e errorProvider) getError() error {
	return e.value
}

func (e errorProvider) getMessage() string {
	return e.message
}

func (e errorProvider) getFormatted() string {
	if len(e.message) == 0 {
		return fmt.Sprintf("error: %s", e.value)
	}
	return fmt.Sprintf("error: %s - %s", e.message, e.value)
}
