package errors

type WithNameError interface {
	GetName() string
}

type OptionalError interface {
	IsOptional() bool
	SetOptional()
}

type ExpectedError interface {
	GetActual() interface{}
	GetExpected() interface{}
}

type assertError struct {
	optional bool
	name     string
	message  string
	actual   interface{}
	expected interface{}
}

func NewAssertError(name string, message string, actual interface{}, expected interface{}) error {
	return &assertError{
		name:     name,
		message:  message,
		actual:   actual,
		expected: expected,
	}
}

func (a *assertError) Error() string {
	return a.message
}

func (a *assertError) GetName() string {
	return a.name
}

func (a *assertError) GetActual() interface{} {
	return a.actual
}

func (a *assertError) GetExpected() interface{} {
	return a.expected
}

func (a *assertError) IsOptional() bool {
	return a.optional
}

func (a *assertError) SetOptional() {
	a.optional = true
}

type optionalError struct {
	error
	optional bool
}

func NewOptionalError(err error) error {
	return &optionalError{
		error:    err,
		optional: true,
	}
}

func (a *optionalError) IsOptional() bool {
	return a.optional
}

func (a *optionalError) SetOptional() {
	a.optional = true
}
