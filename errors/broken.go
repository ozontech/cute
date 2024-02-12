package errors

import "errors"

// BrokenError is an interface for set errors like Broken errors.
// If the function returns an error, which implements this interface, the allure step will has a broken status
type BrokenError interface {
	IsBroken() bool
	SetBroken(bool)
	Error() string
}

type brokenError struct {
	err    error
	broken bool
}

// NewBrokenError ...
func NewBrokenError(err string) error {
	return &brokenError{
		broken: true,
		err:    errors.New(err),
	}
}

// WrapBrokenError returns error with a Broken tag for Allure
func WrapBrokenError(err error) error {
	return &brokenError{
		broken: true,
		err:    err,
	}
}

// Error ..
func (o *brokenError) Error() string {
	return o.err.Error()
}

// IsBroken ...
func (o *brokenError) IsBroken() bool {
	return o.broken
}

// SetBroken ...
func (o *brokenError) SetBroken(broken bool) {
	o.broken = broken
}
