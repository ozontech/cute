package cute

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/ozontech/allure-go/pkg/allure"
)

// AllureBuilder is a scope of methods for create allure information (Title, Tags, etc.)
type AllureBuilder interface {
	AllureInfoBuilder
	AllureLabelsBuilder
	AllureLinksBuilder

	CreateBuilder

	// Parallel signals that this test is to be run in parallel with (and only with) other parallel tests.
	// This function is not thread save. If you use multiply parallel with one T test will panic.
	Parallel() AllureBuilder
}

type CreateBuilder interface {
	// Create is a function for save main information about allure and start write tests
	Create() Middleware

	// CreateWithStep is a function for create step and log some information inside
	// Deprecated use CreateStep(string)
	CreateWithStep() StepBuilder
	// CreateStep is a function for create step inside suite for test
	CreateStep(string) Middleware
}

type AllureInfoBuilder interface {
	// Title is a function for set title in allure information
	Title(title string) AllureBuilder
	// Description is a function for set description in allure information
	Description(description string) AllureBuilder
}

type AllureLinksBuilder interface {
	SetIssue(issue string) AllureBuilder
	SetTestCase(testCase string) AllureBuilder
	Link(link allure.Link) AllureBuilder
}

type AllureLabelsBuilder interface {
	Feature(feature string) AllureBuilder
	Epic(epic string) AllureBuilder
	Tags(tags ...string) AllureBuilder
	ID(value string) AllureBuilder
	AddSuiteLabel(value string) AllureBuilder
	AddSubSuite(value string) AllureBuilder
	AddParentSuite(value string) AllureBuilder
	Story(value string) AllureBuilder
	Tag(value string) AllureBuilder
	Severity(value allure.SeverityType) AllureBuilder
	Owner(value string) AllureBuilder
	Lead(value string) AllureBuilder
	Label(label allure.Label) AllureBuilder
	Labels(labels ...allure.Label) AllureBuilder
}

// StepBuilder is a scope of methods for set step information
// Deprecated.
type StepBuilder interface {
	// StepName is a function to wrap a test in new steps with name
	StepName(name string) Middleware
}

type Middleware interface {
	RequestHTTPBuilder

	// BeforeExecute is function for processing request before createRequest request
	BeforeExecute(...BeforeExecute) Middleware
	// BeforeExecuteT is function for processing request before createRequest request
	BeforeExecuteT(...BeforeExecuteT) Middleware

	// AfterExecute is function will run after allureProvider asserts
	AfterExecute(...AfterExecute) Middleware
	// AfterExecuteT is function will run after allureProvider asserts
	AfterExecuteT(...AfterExecuteT) Middleware
}

// RequestHTTPBuilder is a scope of methods for create HTTP requests
type RequestHTTPBuilder interface {
	// Request is function for set http.Request
	Request(r *http.Request) ExpectHTTPBuilder
	// RequestBuilder is function for create http.Request with help builder.
	// Available builders:
	// WithMethod
	// WithURL
	// WithHeaders
	// WithBody
	// WithMarshalBody
	// WithBody
	// WithURI
	RequestBuilder(r ...requestBuilder) ExpectHTTPBuilder
	// RequestRepeat is a count of repeat request, if request was failed.
	RequestRepeat(count int) RequestHTTPBuilder
	// RequestRepeatDelay is a time between repeat request, if request was failed.
	// Default 1 second
	RequestRepeatDelay(delay time.Duration) RequestHTTPBuilder
}

// ExpectHTTPBuilder is a scope of methods for validate http response
type ExpectHTTPBuilder interface {
	// ExpectExecuteTimeout is function for validate time of execution
	// Default value - 10 seconds
	ExpectExecuteTimeout(t time.Duration) ExpectHTTPBuilder

	// ExpectStatus is function for validate response status code
	ExpectStatus(code int) ExpectHTTPBuilder

	// ExpectJSONSchemaString is function for validate response by json schema from string
	ExpectJSONSchemaString(schema string) ExpectHTTPBuilder
	// ExpectJSONSchemaByte is function for validate response by json schema from byte
	ExpectJSONSchemaByte(schema []byte) ExpectHTTPBuilder
	// ExpectJSONSchemaFile is function for validate response by json schema from file
	// For get file from network use:
	// "http://www.some_host.com/schema.json"
	// For get local file use:
	// "file://./project/me/schema.json"
	ExpectJSONSchemaFile(path string) ExpectHTTPBuilder

	// AssertBody is function for validate response body.
	// Available asserts from asserts/json/json.go:
	// Contains is a function to assert that a jsonpath expression extracts a value in an array
	// Equal is a function to assert that a jsonpath expression matches the given value
	// NotEqual is a function to check jsonpath expression value is not equal to given value
	// Length is a function to asserts that jsonpath expression value is the expected length
	// GreaterThan is a function to asserts that jsonpath expression value is greater than the given length
	// LessThan is a function to asserts that jsonpath expression value is less than the given length
	// Present is a function to asserts that jsonpath expression value is present
	// NotPresent is a function to asserts that jsonpath expression value is not present
	// Also you can write you assert.
	AssertBody(asserts ...AssertBody) ExpectHTTPBuilder
	// OptionalAssertBody is not a mandatory assert.
	// Mark in allure as Broken
	OptionalAssertBody(asserts ...AssertBody) ExpectHTTPBuilder
	// AssertBodyT is function for validate response body with help testing.TB and allure allureProvider.
	// You may create allure step inside assert, add attachment, log information, etc.
	AssertBodyT(asserts ...AssertBodyT) ExpectHTTPBuilder
	// OptionalAssertBodyT is not a mandatory assert.
	// Mark in allure as Broken
	OptionalAssertBodyT(asserts ...AssertBodyT) ExpectHTTPBuilder

	// AssertHeaders is function for validate response headers
	// Available asserts from asserts/headers/headers.go:
	// Present is a function to asserts header is present
	// NotPresent is a function to asserts header is present
	// Also you can write you assert.
	AssertHeaders(asserts ...AssertHeaders) ExpectHTTPBuilder
	// OptionalAssertHeaders is not a mandatory assert.
	// Mark in allure as Broken
	OptionalAssertHeaders(asserts ...AssertHeaders) ExpectHTTPBuilder
	// AssertHeadersT is function for validate headers body with help testing.TB and allure allureProvider.
	// You may create allure step inside assert, add attachment, log information, etc.
	AssertHeadersT(asserts ...AssertHeadersT) ExpectHTTPBuilder
	// OptionalAssertHeadersT is not a mandatory assert.
	// Mark in allure as Broken
	OptionalAssertHeadersT(asserts ...AssertHeadersT) ExpectHTTPBuilder

	// AssertResponse is function for validate response
	AssertResponse(asserts ...AssertResponse) ExpectHTTPBuilder
	// OptionalAssertResponse is not a mandatory assert.
	// Mark in allure as Broken
	OptionalAssertResponse(asserts ...AssertResponse) ExpectHTTPBuilder
	// AssertResponseT is function for validate response with help testing.TB.
	// You may create allure step inside assert, add attachment, log information, etc.
	AssertResponseT(asserts ...AssertResponseT) ExpectHTTPBuilder
	// OptionalAssertResponseT is not a mandatory assert.
	// Mark in allure as Broken
	OptionalAssertResponseT(asserts ...AssertResponseT) ExpectHTTPBuilder

	NextTest() NextTestBuilder

	// ExecuteTest is a function for execute test
	ExecuteTest(ctx context.Context, t testing.TB) []ResultsHTTPBuilder
}

type NextTestBuilder interface {
	ProcessBody([]byte) NextTestBuilder
	ProcessBodyT(T, []byte) NextTestBuilder

	ProcessHeaders(headers http.Header) NextTestBuilder
	ProcessHeadersT(T, headers http.Header) NextTestBuilder

	ProcessBodyResponse(*http.Response) NextTestBuilder
	ProcessBodyResponseT(T, *http.Response) NextTestBuilder

	CreateBuilder
}

// ResultsHTTPBuilder is a scope of methods for processing results
type ResultsHTTPBuilder interface {
	// GetHTTPResponse is a function, which returns http response
	GetHTTPResponse() *http.Response
	// GetErrors is a function, which returns errors
	GetErrors() []error

	// TODO
	GetName() string
}

// BeforeExecute ...
type BeforeExecute func(*http.Request) error

// BeforeExecuteT ...
type BeforeExecuteT func(T, *http.Request) error

// AfterExecute ...
type AfterExecute func(*http.Response, []error) error

// AfterExecuteT ...
type AfterExecuteT func(T, *http.Response, []error) error
