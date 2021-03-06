package errors

const (
	actualField   = "Actual"
	expectedField = "Expected"
)

// WithNameError is interface for creates allure step.
// If function returns error, which implement this interface, allure step will create automatically
type WithNameError interface {
	GetName() string
	SetName(string)
}

// WithFields is interface for put parameters in allure step.
// If function returns error, which implement this interface, parameters will add to allure step
type WithFields interface {
	GetFields() map[string]interface{}
	PutFields(map[string]interface{})
}

// OptionalError is interface for put parameters in allure step.
// If function returns error, which implement this interface, allure step will have skip status
type OptionalError interface {
	IsOptional() bool
	SetOptional(bool)
}

type assertError struct {
	optional bool
	name     string
	message  string
	fields   map[string]interface{}
}

func NewAssertError(name string, message string, actual interface{}, expected interface{}) error {
	return &assertError{
		name:    name,
		message: message,
		fields: map[string]interface{}{
			actualField:   actual,
			expectedField: expected,
		},
	}
}

func (a *assertError) Error() string {
	return a.message
}

func (a *assertError) GetName() string {
	return a.name
}

func (a *assertError) SetName(name string) {
	a.name = name
}

func (a *assertError) GetFields() map[string]interface{} {
	return a.fields
}

func (a *assertError) PutFields(fields map[string]interface{}) {
	for k, v := range fields {
		a.fields[k] = v
	}
}

func (a *assertError) IsOptional() bool {
	return a.optional
}

func (a *assertError) SetOptional(opt bool) {
	a.optional = opt
}

type optionalError struct {
	err      string
	optional bool
}

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
