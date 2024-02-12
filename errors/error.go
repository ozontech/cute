package errors

const (
	actualField   = "Actual"
	expectedField = "Expected"
)

// CuteError ...
type CuteError interface {
	error
	WithNameError
	WithFields
	WithAttachments
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

type assertError struct {
	optional bool
	require  bool
	broken   bool

	name        string
	message     string
	fields      map[string]interface{}
	attachments []*Attachment
}

// NewAssertError is the function, which creates error with "Actual" and "Expected" for allure
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

// NewAssertErrorWithMessage ...
func NewAssertErrorWithMessage(name string, message string) error {
	return &assertError{
		name:    name,
		message: message,
	}
}

// NewEmptyAssertError ...
func NewEmptyAssertError(name string, message string) CuteError {
	return &assertError{
		name:    name,
		message: message,
		fields:  map[string]interface{}{},
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

func (a *assertError) GetAttachments() []*Attachment {
	return a.attachments
}

func (a *assertError) PutAttachment(attachment *Attachment) {
	a.attachments = append(a.attachments, attachment)
}

func (a *assertError) IsOptional() bool {
	return a.optional
}

func (a *assertError) SetOptional(opt bool) {
	a.optional = opt
}

func (a *assertError) IsRequire() bool {
	return a.require
}

func (a *assertError) SetRequire(b bool) {
	a.require = b
}

func (a *assertError) IsBroken() bool {
	return a.broken
}

func (a *assertError) SetBroken(b bool) {
	a.broken = b
}
