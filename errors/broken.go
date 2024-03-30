package errors

// BrokenError is an interface for set errors like Broken errors.
// If the function returns an error, which implements this interface, the allure step will has a broken status
type BrokenError interface {
	IsBroken() bool
	SetBroken(bool)
	Error() string
}

// NewBrokenError returns error with a Broken tag for Allure
func NewBrokenError(err string) error {
	return &CuteError{
		Broken:  true,
		Message: err,
	}
}

// WrapBrokenError returns error with a Broken tag for Allure
func WrapBrokenError(err error) error {
	return &CuteError{
		Broken: true,
		Err:    err,
	}
}
