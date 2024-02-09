package errors

// BrokenError is interface for set error like Broken error.
// If function returns error, which implement this interface, allure step will have failed status
type BrokenError interface {
	IsBroken() bool
	SetBroken(bool)
	Error() string
}

type brokenError struct {
	err    string
	broken bool
}

// NewBrokenError ...
func NewBrokenError(err string) error {
	return &brokenError{
		broken: true,
		err:    err,
	}
}

// Error ..
func (o *brokenError) Error() string {
	return o.err
}

// IsBroken ...
func (o *brokenError) IsBroken() bool {
	return o.broken
}

// SetBroken ...
func (o *brokenError) SetBroken(Broken bool) {
	o.broken = Broken
}
