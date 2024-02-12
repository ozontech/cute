package errors

import "errors"

// RequireError is an interface for set errors like require error.
// If the function returns an error, which implements this interface, the allure step will has failed status
type RequireError interface {
	IsRequire() bool
	SetRequire(bool)
}

type requireError struct {
	err     error
	require bool
}

// NewRequireError ...
func NewRequireError(err string) error {
	return &requireError{
		require: true,
		err:     errors.New(err),
	}
}

// WrapRequireError returns error with flag for execute t.FailNow() and finish test after this error
func WrapRequireError(err error) error {
	return &requireError{
		require: true,
		err:     err,
	}
}

// Error ..
func (o *requireError) Error() string {
	return o.err.Error()
}

// IsRequire ...
func (o *requireError) IsRequire() bool {
	return o.require
}

// SetRequire ...
func (o *requireError) SetRequire(require bool) {
	o.require = require
}
