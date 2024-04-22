package errors

import "fmt"

const (
	// ActualField is a key for actual value in error fields
	ActualField = "Actual"
	// ExpectedField is a key for expected value in error fields
	ExpectedField = "Expected"
)

// AssertError is a common interface for all errors in the package
type AssertError interface {
	error
	WithNameError
	WithFields
	WithAttachments
	WithTrace
}

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

// WithTrace is interface for put trace in logs
type WithTrace interface {
	GetTrace() string
	SetTrace(string)
}

// Attachment represents an attachment to Allure with properties like name, MIME type, and content.
type Attachment struct {
	Name     string // Name of the attachment.
	MimeType string // MIME type of the attachment.
	Content  []byte // Content of the attachment.
}

// WithAttachments is an interface that defines methods for managing attachments.
type WithAttachments interface {
	GetAttachments() []*Attachment
	PutAttachment(a *Attachment)
}

// CuteError is a struct for error with additional fields for allure and logs
type CuteError struct {
	// Optional is a flag to determine if the error is optional
	// If the error is optional, it will not fail the test
	Optional bool
	// Require is a flag to determine if the error is required
	// If the error is required, it will fail the test
	Require bool
	// Broken is a flag to determine if the error is broken
	// If the error is broken, it will fail the test and mark the test as broken in allure
	Broken bool

	// Name is a name of the error
	Name string
	// Message is a message of the error
	Message string
	// Err is a wrapped error
	Err error

	// Trace is a trace of the error
	// It could be a file path, function name, or any other information
	Trace string

	// Fields is a map of additional fields for the error
	// It could be actual and expected values, parameters, or any other information
	// ActualField and ExpectedField fields will be logged
	Fields map[string]interface{}
	// Attachments is a slice of attachments for the error
	Attachments []*Attachment
}

// NewCuteError is the function, which creates cute error with "Name" and "Message" for allure
func NewCuteError(name string, err error) *CuteError {
	return &CuteError{
		Name: name,
		Err:  err,
	}
}

// NewAssertError is the function, which creates error with "Actual" and "Expected" for allure
func NewAssertError(name string, message string, actual interface{}, expected interface{}) error {
	return &CuteError{
		Name:    name,
		Message: message,
		Fields: map[string]interface{}{
			ActualField:   actual,
			ExpectedField: expected,
		},
	}
}

// NewAssertErrorWithMessage is the function, which creates error with "Name" and "Message" for allure
// Deprecated: use NewEmptyAssertError instead
func NewAssertErrorWithMessage(name string, message string) error {
	return NewEmptyAssertError(name, message)
}

// NewEmptyAssertError is the function, which creates error with "Name" and "Message" for allure
// Returns AssertError with empty fields
// You can use PutFields and PutAttachment to add additional information
// You can use SetOptional, SetRequire, SetBroken to change error behavior
func NewEmptyAssertError(name string, message string) AssertError {
	return &CuteError{
		Name:    name,
		Message: message,
		Fields:  map[string]interface{}{},
	}
}

// Unwrap is a method to get wrapped error
// It is used for errors.Is and errors.As functions
func (a *CuteError) Unwrap() error {
	return a.Err
}

// Error is a method to get error message
// It is used for fmt.Errorf and fmt.Println functions
func (a *CuteError) Error() string {
	if a.Trace == "" {
		return a.Message
	}

	errText := a.Message

	if a.Err != nil {
		errText = a.Err.Error()
	}

	return fmt.Sprintf("%s\nCalled from: %s", errText, a.Trace)
}

// GetName is a method to get error name
// It is used for allure step name
func (a *CuteError) GetName() string {
	return a.Name
}

// SetName is a method to set error name
// It is used for allure step name
func (a *CuteError) SetName(name string) {
	a.Name = name
}

// GetFields ...
func (a *CuteError) GetFields() map[string]interface{} {
	return a.Fields
}

// PutFields ...
func (a *CuteError) PutFields(fields map[string]interface{}) {
	for k, v := range fields {
		a.Fields[k] = v
	}
}

// GetAttachments ...
func (a *CuteError) GetAttachments() []*Attachment {
	return a.Attachments
}

// PutAttachment ...
func (a *CuteError) PutAttachment(attachment *Attachment) {
	a.Attachments = append(a.Attachments, attachment)
}

// IsOptional ...
func (a *CuteError) IsOptional() bool {
	return a.Optional
}

// SetOptional ...
func (a *CuteError) SetOptional(opt bool) {
	a.Optional = opt
}

// IsRequire ...
func (a *CuteError) IsRequire() bool {
	return a.Require
}

// SetRequire ...
func (a *CuteError) SetRequire(b bool) {
	a.Require = b
}

// IsBroken ...
func (a *CuteError) IsBroken() bool {
	return a.Broken
}

// SetBroken ...
func (a *CuteError) SetBroken(b bool) {
	a.Broken = b
}

// GetTrace ...
func (a *CuteError) GetTrace() string {
	return a.Trace
}

// SetTrace ...
func (a *CuteError) SetTrace(trace string) {
	a.Trace = trace
}
