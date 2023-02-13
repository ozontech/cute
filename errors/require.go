package errors

// RequireError is interface for set error like require error.
// If function returns error, which implement this interface, allure step will have failed status
type RequireError interface {
	IsRequire() bool
	SetRequire(bool)
}

type requireError struct {
	err     string
	require bool
}

// NewRequireError ...
func NewRequireError(err string) error {
	return &requireError{
		require: true,
		err:     err,
	}
}

// Error ..
func (o *requireError) Error() string {
	return o.err
}

// IsRequire ...
func (o *requireError) IsRequire() bool {
	return o.require
}

// SetRequire ...
func (o *requireError) SetRequire(require bool) {
	o.require = require
}
