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

type CuteError struct {
	Optional bool
	Require  bool
	Broken   bool

	Name        string
	Message     string
	Err         error
	Trace       string
	Fields      map[string]interface{}
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

func (a *CuteError) Unwrap() error {
	return a.Err
}

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

func (a *CuteError) GetName() string {
	return a.Name
}

func (a *CuteError) SetName(name string) {
	a.Name = name
}

func (a *CuteError) GetFields() map[string]interface{} {
	return a.Fields
}

func (a *CuteError) PutFields(fields map[string]interface{}) {
	for k, v := range fields {
		a.Fields[k] = v
	}
}

func (a *CuteError) GetAttachments() []*Attachment {
	return a.Attachments
}

func (a *CuteError) PutAttachment(attachment *Attachment) {
	a.Attachments = append(a.Attachments, attachment)
}

func (a *CuteError) IsOptional() bool {
	return a.Optional
}

func (a *CuteError) SetOptional(opt bool) {
	a.Optional = opt
}

func (a *CuteError) IsRequire() bool {
	return a.Require
}

func (a *CuteError) SetRequire(b bool) {
	a.Require = b
}

func (a *CuteError) IsBroken() bool {
	return a.Broken
}

func (a *CuteError) SetBroken(b bool) {
	a.Broken = b
}

func (a *CuteError) GetTrace() string {
	return a.Trace
}

func (a *CuteError) SetTrace(trace string) {
	a.Trace = trace
}
