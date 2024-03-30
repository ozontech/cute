package errors

// OptionalError is an interface for set errors like Optional errors.
// If the function returns an error, which implements this interface, the allure step will has to skip status
type OptionalError interface {
	IsOptional() bool
	SetOptional(bool)
}

// NewOptionalError returns error with an Optional tag for Allure
func NewOptionalError(err string) error {
	return &CuteError{
		Optional: true,
		Message:  err,
	}
}

// WrapOptionalError returns error with an Optional tag for Allure
func WrapOptionalError(err error) error {
	return &CuteError{
		Optional: true,
		Err:      err,
	}
}
