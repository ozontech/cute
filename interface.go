package cute

import (
	"context"
	"net/http"
	"time"

	"github.com/ozontech/allure-go/pkg/allure"
)

// AllureBuilder is a scope of methods for create allure information (Title, Tags, etc.)
type AllureBuilder interface {
	AllureInfoBuilder
	AllureLabelsBuilder
	AllureLinksBuilder

	CreateBuilder

	// Parallel signals that this Test is to be run in parallel with (and only with) other parallel tests.
	// This function is not thread save. If you use multiply parallel with one T Test will panic.
	Parallel() AllureBuilder
}

// AllureInfoBuilder is a scope of methods for create allure information (Title, Tags, etc.)
type AllureInfoBuilder interface {
	// Title is a function for set title in allure information
	Title(title string) AllureBuilder
	Titlef(format string, args ...interface{}) AllureBuilder
	// Description is a function for set description in allure information
	Description(description string) AllureBuilder
	Descriptionf(format string, args ...interface{}) AllureBuilder
	Stage(stage string) AllureBuilder
	Stagef(format string, args ...interface{}) AllureBuilder
}

// AllureLinksBuilder is a scope of methods to set allure links
type AllureLinksBuilder interface {
	SetIssue(issue string) AllureBuilder
	SetTestCase(testCase string) AllureBuilder
	Link(link *allure.Link) AllureBuilder
	TmsLink(tmsLink string) AllureBuilder
	TmsLinks(tmsLinks ...string) AllureBuilder
}

// AllureLabelsBuilder is a scope of methods to set allure labels
type AllureLabelsBuilder interface {
	Feature(feature string) AllureBuilder
	Epic(epic string) AllureBuilder
	AllureID(value string) AllureBuilder
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
	Label(label *allure.Label) AllureBuilder
	Labels(labels ...*allure.Label) AllureBuilder
	Layer(value string) AllureBuilder
	Stagef(format string, args ...interface{}) AllureBuilder
	Stage(stage string) AllureBuilder
}

// CreateBuilder is functions for create Test or table tests
type CreateBuilder interface {
	// Create is a function for save main information about allure and start write tests
	Create() MiddlewareRequest

	// CreateStep is a function for create step inside suite for Test
	CreateStep(string) MiddlewareRequest

	// CreateTableTest is function for create table Test
	CreateTableTest() MiddlewareTable
}

// MiddlewareTable is functions for create table Test
type MiddlewareTable interface {
	TableTest

	BeforeTest
	AfterTest
}

// MiddlewareRequest is function for create requests or add After/Before functions
type MiddlewareRequest interface {
	RequestHTTPBuilder

	BeforeTest
	AfterTest
}

// BeforeTest are functions for processing request before test execution
// Same functions:
// Before
type BeforeTest interface {
	// BeforeExecute is function for processing request before test execution
	BeforeExecute(...BeforeExecute) MiddlewareRequest
	// BeforeExecuteT is function for processing request before test execution
	BeforeExecuteT(...BeforeExecuteT) MiddlewareRequest
}

// After are functions for processing response after test execution
// Same functions:
// AfterText
// AfterTestExecute
type After interface {
	// After is function for processing response after test execution
	After(...AfterExecute) ExpectHTTPBuilder
	// AfterT is function for processing response after test execution
	AfterT(...AfterExecuteT) ExpectHTTPBuilder
}

// AfterTest are functions for processing response after test execution
// Same functions:
// After
// AfterTestExecute
type AfterTest interface {
	// AfterExecute is function for processing response after test execution
	AfterExecute(...AfterExecute) MiddlewareRequest
	// AfterExecuteT is function for processing response after test execution
	AfterExecuteT(...AfterExecuteT) MiddlewareRequest
}

// AfterTestExecute are functions for processing response after test execution
// Same functions:
// After
// AfterText
type AfterTestExecute interface {
	// AfterTestExecute is function for processing response after test execution
	AfterTestExecute(...AfterExecute) NextTestBuilder
	// AfterTestExecuteT is function for processing response after test execution
	AfterTestExecuteT(...AfterExecuteT) NextTestBuilder
}

// TableTest is function for put request and assert for table tests
type TableTest interface {
	// PutNewTest is function for put request and assert for table Test
	PutNewTest(name string, r *http.Request, expect *Expect) TableTest
	// PutTests is function for put requests and asserts for table Test
	PutTests(params ...*Test) TableTest
	ControlTest
}

// RequestHTTPBuilder is a scope of methods to create HTTP requests
type RequestHTTPBuilder interface {
	// Request is function for set http.Request
	Request(r *http.Request) ExpectHTTPBuilder
	// RequestBuilder is function for set http.Request with builders
	// Available builders:
	// WithMethod
	// WithURL
	// WithHeaders
	// WithHeadersKV
	// WithBody
	// WithMarshalBody
	// WithBody
	// WithURI
	// WithQuery
	// WithQueryKV
	// WithFileForm
	// WithFileFormKV
	// WithForm
	// WithFormKV
	RequestBuilder(r ...RequestBuilder) ExpectHTTPBuilder

	RequestParams
}

// RequestParams is a scope of methods to configure request
type RequestParams interface {
	// RequestRepeat is a function for set options in request
	// if response.Code != Expect.Code, than request will repeat counts with delay.
	// Default delay is 1 second.
	RequestRepeat(count int) RequestHTTPBuilder

	// RequestRepeatDelay set delay for request repeat.
	// if response.Code != Expect.Code, than request will repeat counts with delay.
	// Default delay is 1 second.
	RequestRepeatDelay(delay time.Duration) RequestHTTPBuilder

	// RequestRepeatPolitic is a politic for repeat request.
	// if response.Code != Expect.Code, than request will repeat counts with delay.
	// if Optional is true and request is failed, than test step allure will be skipped, and t.Fail() will not execute.
	// If Broken is true and request is failed, than test step allure will be broken, and t.Fail() will execute.
	RequestRepeatPolitic(politic *RequestRepeatPolitic) RequestHTTPBuilder

	// RequestRepeatOptional is a option politic for repeat request.
	// if Optional is true and request is failed, than test step allure will be skipped, and t.Fail() will not execute.
	RequestRepeatOptional(optional bool) RequestHTTPBuilder

	// RequestRepeatBroken is a broken politic for repeat request.
	// If Broken is true and request is failed, than test step allure will be broken, and t.Fail() will execute.
	RequestRepeatBroken(broken bool) RequestHTTPBuilder
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
	// RequireBody implements the same assertions as the `AssertBody`, but stops test execution when a test fails.
	RequireBody(asserts ...AssertBody) ExpectHTTPBuilder
	// OptionalAssertBody is not a mandatory assert.
	// Mark in allure as Skipped
	OptionalAssertBody(asserts ...AssertBody) ExpectHTTPBuilder
	// BrokenAssertBody  is function for validate response, if it's failed, then test will be Broken.
	// Mark in allure as Broken
	BrokenAssertBody(asserts ...AssertBody) ExpectHTTPBuilder
	// AssertBodyT is function for validate response body with help testing.TB and allure allureProvider.
	// You may create allure step inside assert, add attachment, log information, etc.
	AssertBodyT(asserts ...AssertBodyT) ExpectHTTPBuilder
	// RequireBodyT implements the same assertions as the `AssertBodyT`, but stops test execution when a test fails.
	RequireBodyT(asserts ...AssertBodyT) ExpectHTTPBuilder
	// OptionalAssertBodyT is not a mandatory assert.
	// Mark in allure as Skipped
	OptionalAssertBodyT(asserts ...AssertBodyT) ExpectHTTPBuilder
	// BrokenAssertBodyT  is function for validate response, if it's failed, then test will be Broken.
	// Mark in allure as Broken
	BrokenAssertBodyT(asserts ...AssertBodyT) ExpectHTTPBuilder

	// AssertHeaders is function for validate response headers
	// Available asserts from asserts/headers/headers.go:
	// Present is a function to asserts header is present
	// NotPresent is a function to asserts header is present
	// Also you can write you assert.
	AssertHeaders(asserts ...AssertHeaders) ExpectHTTPBuilder
	// RequireHeaders implements the same assertions as the `AssertHeaders`, but stops test execution when a test fails.
	RequireHeaders(asserts ...AssertHeaders) ExpectHTTPBuilder
	// OptionalAssertHeaders is not a mandatory assert.
	// Mark in allure as Skipped
	OptionalAssertHeaders(asserts ...AssertHeaders) ExpectHTTPBuilder
	// BrokenAssertHeaders  is function for validate response, if it's failed, then test will be Broken.
	// Mark in allure as Broken
	BrokenAssertHeaders(asserts ...AssertHeaders) ExpectHTTPBuilder
	// AssertHeadersT is function for validate headers body with help testing.TB and allure allureProvider.
	// You may create allure step inside assert, add attachment, log information, etc.
	AssertHeadersT(asserts ...AssertHeadersT) ExpectHTTPBuilder
	// RequireHeadersT implements the same assertions as the `AssertHeadersT`, but stops test execution when a test fails.
	RequireHeadersT(asserts ...AssertHeadersT) ExpectHTTPBuilder
	// OptionalAssertHeadersT is not a mandatory assert.
	// Mark in allure as Skipped
	OptionalAssertHeadersT(asserts ...AssertHeadersT) ExpectHTTPBuilder
	// BrokenAssertHeadersT is function for validate response, if it's failed, then test will be Broken.
	// Mark in allure as Broken
	BrokenAssertHeadersT(asserts ...AssertHeadersT) ExpectHTTPBuilder

	// AssertResponse is function for validate response.
	AssertResponse(asserts ...AssertResponse) ExpectHTTPBuilder
	// RequireResponse implements the same assertions as the `AssertResponse`, but stops test execution when a test fails.
	RequireResponse(asserts ...AssertResponse) ExpectHTTPBuilder
	// OptionalAssertResponse is not a mandatory assert.
	// Mark in allure as Skipped
	OptionalAssertResponse(asserts ...AssertResponse) ExpectHTTPBuilder
	// BrokenAssertResponse  is function for validate response, if it's failed, then test will be Broken.
	// Mark in allure as Broken
	BrokenAssertResponse(asserts ...AssertResponse) ExpectHTTPBuilder
	// AssertResponseT is function for validate response with help testing.TB.
	// You may create allure step inside assert, add attachment, log information, etc.
	AssertResponseT(asserts ...AssertResponseT) ExpectHTTPBuilder
	// RequireResponseT implements the same assertions as the `AssertResponseT`, but stops test execution when a test fails.
	RequireResponseT(asserts ...AssertResponseT) ExpectHTTPBuilder
	// OptionalAssertResponseT is not a mandatory assert.
	// Mark in allure as Skipped
	OptionalAssertResponseT(asserts ...AssertResponseT) ExpectHTTPBuilder
	// BrokenAssertResponseT is function for validate response, if it's failed, then test will be Broken.
	// Mark in allure as Broken
	BrokenAssertResponseT(asserts ...AssertResponseT) ExpectHTTPBuilder

	After
	ControlTest
}

// ControlTest is function for manipulating tests
type ControlTest interface {
	NextTest() NextTestBuilder

	// ExecuteTest is a function for execute Test
	ExecuteTest(ctx context.Context, t tProvider) []ResultsHTTPBuilder
}

// NextTestBuilder is a scope of methods for processing response, after Test.
type NextTestBuilder interface {
	AfterTestExecute

	CreateBuilder
}

// ResultsHTTPBuilder is a scope of methods for processing results
type ResultsHTTPBuilder interface {
	// GetHTTPResponse is a function, which returns http response
	GetHTTPResponse() *http.Response
	// GetErrors is a function, which returns all errors from test
	GetErrors() []error
	// GetName is a function, which returns name of Test
	GetName() string
	// GetResultState is a function, which returns state of test
	// State could be ResultStateSuccess, ResultStateBroken, ResultStateFail
	GetResultState() ResultState
}

// BeforeExecute is a function for processing request before test execution
type BeforeExecute func(*http.Request) error

// BeforeExecuteT is a function for processing request before test execution
type BeforeExecuteT func(T, *http.Request) error

// AfterExecute is a function for processing response after test execution
type AfterExecute func(*http.Response, []error) error

// AfterExecuteT is a function for processing response after test execution
type AfterExecuteT func(T, *http.Response, []error) error
