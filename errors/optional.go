package errors

// OptionalError is interface for set error like optional error.
// If function returns error, which implement this interface, allure step will have skip status
type OptionalError interface {
	IsOptional() bool
	SetOptional(bool)
}

type optionalError struct {
	err      string
	optional bool
}

// NewOptionalError ...
func NewOptionalError(err string) error {
	return &optionalError{
		optional: true,
		err:      err,
	}
}

func (o *optionalError) Error() string {
	return o.err
}

func (o *optionalError) IsOptional() bool {
	return o.optional
}

func (o *optionalError) SetOptional(opt bool) {
	o.optional = opt
}
