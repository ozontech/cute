package errors

import "errors"

// BrokenError is interface for set error like Broken error.
// If function returns error, which implement this interface, allure step will has broken status
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

// WrapBrokenError ...
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
