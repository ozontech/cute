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

// NewAssertErrorWithMessage ...
func NewAssertErrorWithMessage(name string, message string) error {
	return &CuteError{
		Name:    name,
		Message: message,
	}
}

// NewEmptyAssertError ...
func NewEmptyAssertError(name string, message string) AssertError {
	return &CuteError{
		Name:    name,
		Message: message,
		Fields:  map[string]interface{}{},
	}
}

// Unwrap ...
func (a *CuteError) Unwrap() error {
	return a.Err
}

// Error ...
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

// GetName ...
func (a *CuteError) GetName() string {
	return a.Name
}

// SetName ...
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
