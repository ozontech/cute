package errors

import "errors"

// OptionalError is an interface for set errors like optional errors.
// If the function returns an error, which implements this interface, the allure step will has to skip status
type OptionalError interface {
	IsOptional() bool
	SetOptional(bool)
}

type optionalError struct {
	err      error
	optional bool
}

// NewOptionalError ...
func NewOptionalError(err string) error {
	return &optionalError{
		optional: true,
		err:      errors.New(err),
	}
}

// WrapOptionalError returns error with an Optional tag for Allure
func WrapOptionalError(err error) error {
	return &optionalError{
		optional: true,
		err:      err,
	}
}

func (o *optionalError) Error() string {
	return o.err.Error()
}

func (o *optionalError) IsOptional() bool {
	return o.optional
}

func (o *optionalError) SetOptional(opt bool) {
	o.optional = opt
}
