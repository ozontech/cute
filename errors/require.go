package errors

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

// NewRequireError returns error with flag for execute t.FailNow() and finish test after this error
func NewRequireError(err string) error {
	return &CuteError{
		Require: true,
		Message: err,
	}
}

// WrapRequireError returns error with flag for execute t.FailNow() and finish test after this error
func WrapRequireError(err error) error {
	return &CuteError{
		Require: true,
		Err:     err,
	}
}
