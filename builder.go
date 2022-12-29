package cute

import (
	"net/http"
	"time"

	"github.com/ozontech/allure-go/pkg/allure"
)

const defaultHTTPTimeout = 30

var (
	errorAssertIsNil = "Assert must be not nil"
)

// HTTPTestMaker is a creator tests
type HTTPTestMaker struct {
	httpClient *http.Client
	middleware *Middleware

	hardValidation bool

	// todo add marshaler
}

type options struct {
	httpClient       *http.Client
	httpTimeout      time.Duration
	httpRoundTripper http.RoundTripper

	middleware *Middleware

	hardValidation bool
}

type Option func(*options)

// WithHTTPClient is a function for set custom http client
func WithHTTPClient(client *http.Client) Option {
	return func(o *options) {
		o.httpClient = client
	}
}

// WithCustomHTTPTimeout is a function for set custom http client timeout
func WithCustomHTTPTimeout(t time.Duration) Option {
	return func(o *options) {
		o.httpTimeout = t
	}
}

// WithCustomHTTPRoundTripper is a function for set custom http round tripper
func WithCustomHTTPRoundTripper(r http.RoundTripper) Option {
	return func(o *options) {
		o.httpRoundTripper = r
	}
}

// WithMiddlewareAfter ...
func WithMiddlewareAfter(after ...AfterExecute) Option {
	return func(o *options) {
		o.middleware.After = append(o.middleware.After, after...)
	}
}

// WithMiddlewareAfterT ...
func WithMiddlewareAfterT(after ...AfterExecuteT) Option {
	return func(o *options) {
		o.middleware.AfterT = append(o.middleware.AfterT, after...)
	}
}

// WithMiddlewareBefore ...
func WithMiddlewareBefore(before ...BeforeExecute) Option {
	return func(o *options) {
		o.middleware.Before = append(o.middleware.Before, before...)
	}
}

// WithMiddlewareBeforeT ...
func WithMiddlewareBeforeT(beforeT ...BeforeExecuteT) Option {
	return func(o *options) {
		o.middleware.BeforeT = append(o.middleware.BeforeT, beforeT...)
	}
}

// WithHardValidation ...
func WithHardValidation() Option {
	return func(o *options) {
		o.hardValidation = true
	}
}

// NewHTTPTestMaker is function for set options for all cute.
func NewHTTPTestMaker(opts ...Option) *HTTPTestMaker {
	var (
		o = &options{
			middleware: new(Middleware),
		}

		timeout      = defaultHTTPTimeout * time.Second
		roundTripper = http.DefaultTransport
	)

	for _, opt := range opts {
		opt(o)
	}

	if o.httpTimeout != 0 {
		timeout = o.httpTimeout
	}

	if o.httpRoundTripper != nil { //nolint
		roundTripper = o.httpRoundTripper
	}

	httpClient := &http.Client{
		Transport: roundTripper,
		Timeout:   timeout,
	}

	if o.httpClient != nil {
		httpClient = o.httpClient
	}

	m := &HTTPTestMaker{
		hardValidation: o.hardValidation,
		httpClient:     httpClient,
		middleware:     o.middleware,
	}

	return m
}

// NewTestBuilder is a function for initialization foundation for cute
func (m *HTTPTestMaker) NewTestBuilder() AllureBuilder {
	tests := createDefaultTests(m)

	return &cute{
		baseProps:    m,
		countTests:   0,
		tests:        tests,
		allureInfo:   new(allureInformation),
		allureLinks:  new(allureLinks),
		allureLabels: new(allureLabels),
		parallel:     false,
	}
}

func createDefaultTests(m *HTTPTestMaker) []*Test {
	tests := make([]*Test, 1)
	tests[0] = createDefaultTest(m)

	return tests
}

func createDefaultTest(m *HTTPTestMaker) *Test {
	after := make([]AfterExecute, 0, len(m.middleware.After))
	after = append(after, m.middleware.After...)

	afterT := make([]AfterExecuteT, 0, len(m.middleware.AfterT))
	afterT = append(afterT, m.middleware.AfterT...)

	before := make([]BeforeExecute, 0, len(m.middleware.Before))
	before = append(before, m.middleware.Before...)

	beforeT := make([]BeforeExecuteT, 0, len(m.middleware.BeforeT))
	beforeT = append(beforeT, m.middleware.BeforeT...)

	middleware := &Middleware{
		After:   after,
		AfterT:  afterT,
		Before:  before,
		BeforeT: beforeT,
	}

	return &Test{
		HardValidation: m.hardValidation,
		httpClient:     m.httpClient,
		Middleware:     middleware,
		AllureStep:     new(AllureStep),
		Request: &Request{
			Repeat: new(RequestRepeatPolitic),
		},
		Expect: &Expect{JSONSchema: new(ExpectJSONSchema)},
	}
}

func (it *cute) Title(title string) AllureBuilder {
	it.allureInfo.title = title

	return it
}

func (it *cute) Epic(epic string) AllureBuilder {
	it.allureLabels.epic = epic

	return it
}

func (it *cute) SetIssue(issue string) AllureBuilder {
	it.allureLinks.issue = issue

	return it
}

func (it *cute) SetTestCase(testCase string) AllureBuilder {
	it.allureLinks.testCase = testCase

	return it
}

func (it *cute) Link(link *allure.Link) AllureBuilder {
	it.allureLinks.link = link

	return it
}

func (it *cute) ID(value string) AllureBuilder {
	it.allureLabels.id = value

	return it
}

func (it *cute) AllureID(value string) AllureBuilder {
	it.allureLabels.allureID = value

	return it
}

func (it *cute) AddSuiteLabel(value string) AllureBuilder {
	it.allureLabels.suiteLabel = value

	return it
}

func (it *cute) AddSubSuite(value string) AllureBuilder {
	it.allureLabels.subSuite = value

	return it
}

func (it *cute) AddParentSuite(value string) AllureBuilder {
	it.allureLabels.parentSuite = value

	return it
}

func (it *cute) Story(value string) AllureBuilder {
	it.allureLabels.story = value

	return it
}

func (it *cute) Tag(value string) AllureBuilder {
	it.allureLabels.tag = value

	return it
}

func (it *cute) Severity(value allure.SeverityType) AllureBuilder {
	it.allureLabels.severity = value

	return it
}

func (it *cute) Owner(value string) AllureBuilder {
	it.allureLabels.owner = value

	return it
}

func (it *cute) Lead(value string) AllureBuilder {
	it.allureLabels.lead = value

	return it
}

func (it *cute) Label(label *allure.Label) AllureBuilder {
	it.allureLabels.label = label

	return it
}

func (it *cute) Labels(labels ...*allure.Label) AllureBuilder {
	it.allureLabels.labels = labels

	return it
}

func (it *cute) Description(description string) AllureBuilder {
	it.allureInfo.description = description

	return it
}

func (it *cute) Tags(tags ...string) AllureBuilder {
	it.allureLabels.tags = tags

	return it
}

func (it *cute) Feature(feature string) AllureBuilder {
	it.allureLabels.feature = feature

	return it
}

func (it *cute) Create() MiddlewareRequest {
	return it
}

func (it *cute) CreateStep(name string) MiddlewareRequest {
	it.tests[it.countTests].AllureStep.Name = name

	return it
}

func (it *cute) Parallel() AllureBuilder {
	it.parallel = true

	return it
}

func (it *cute) CreateRequest() RequestHTTPBuilder {
	return it
}

func (it *cute) StepName(name string) MiddlewareRequest {
	it.tests[it.countTests].AllureStep.Name = name

	return it
}

func (it *cute) BeforeExecute(fs ...BeforeExecute) MiddlewareRequest {
	it.tests[it.countTests].Middleware.Before = append(it.tests[it.countTests].Middleware.Before, fs...)

	return it
}

func (it *cute) BeforeExecuteT(fs ...BeforeExecuteT) MiddlewareRequest {
	it.tests[it.countTests].Middleware.BeforeT = append(it.tests[it.countTests].Middleware.BeforeT, fs...)

	return it
}

func (it *cute) After(fs ...AfterExecute) ExpectHTTPBuilder {
	it.tests[it.countTests].Middleware.After = append(it.tests[it.countTests].Middleware.After, fs...)

	return it
}

func (it *cute) AfterT(fs ...AfterExecuteT) ExpectHTTPBuilder {
	it.tests[it.countTests].Middleware.AfterT = append(it.tests[it.countTests].Middleware.AfterT, fs...)

	return it
}

func (it *cute) AfterExecute(fs ...AfterExecute) MiddlewareRequest {
	it.tests[it.countTests].Middleware.After = append(it.tests[it.countTests].Middleware.After, fs...)

	return it
}

func (it *cute) AfterExecuteT(fs ...AfterExecuteT) MiddlewareRequest {
	it.tests[it.countTests].Middleware.AfterT = append(it.tests[it.countTests].Middleware.AfterT, fs...)

	return it
}

func (it *cute) AfterTestExecute(fs ...AfterExecute) NextTestBuilder {
	previousTest := 0
	if it.countTests != 0 {
		previousTest = it.countTests - 1
	}

	it.tests[previousTest].Middleware.After = append(it.tests[previousTest].Middleware.After, fs...)

	return it
}

func (it *cute) AfterTestExecuteT(fs ...AfterExecuteT) NextTestBuilder {
	previousTest := 0
	if it.countTests != 0 {
		previousTest = it.countTests - 1
	}

	it.tests[previousTest].Middleware.AfterT = append(it.tests[previousTest].Middleware.AfterT, fs...)

	return it
}

func (it *cute) RequestRepeat(count int) RequestHTTPBuilder {
	it.tests[it.countTests].Request.Repeat.Count = count

	return it
}

func (it *cute) RequestRepeatDelay(delay time.Duration) RequestHTTPBuilder {
	it.tests[it.countTests].Request.Repeat.Delay = delay

	return it
}

func (it *cute) Request(r *http.Request) ExpectHTTPBuilder {
	it.tests[it.countTests].Request.Base = r

	return it
}

func (it *cute) RequestBuilder(r ...RequestBuilder) ExpectHTTPBuilder {
	it.tests[it.countTests].Request.Builders = append(it.tests[it.countTests].Request.Builders, r...)

	return it
}

func (it *cute) ExpectExecuteTimeout(t time.Duration) ExpectHTTPBuilder {
	it.tests[it.countTests].Expect.ExecuteTime = t

	return it
}

func (it *cute) ExpectStatus(code int) ExpectHTTPBuilder {
	it.tests[it.countTests].Expect.Code = code

	return it
}

func (it *cute) ExpectJSONSchemaString(schema string) ExpectHTTPBuilder {
	it.tests[it.countTests].Expect.JSONSchema.String = schema

	return it
}

func (it *cute) ExpectJSONSchemaByte(schema []byte) ExpectHTTPBuilder {
	it.tests[it.countTests].Expect.JSONSchema.Byte = schema

	return it
}

func (it *cute) ExpectJSONSchemaFile(filePath string) ExpectHTTPBuilder {
	it.tests[it.countTests].Expect.JSONSchema.File = filePath

	return it
}

func (it *cute) AssertBody(asserts ...AssertBody) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}
	}

	it.tests[it.countTests].Expect.AssertBody = append(it.tests[it.countTests].Expect.AssertBody, asserts...)

	return it
}

func (it *cute) OptionalAssertBody(asserts ...AssertBody) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		it.tests[it.countTests].Expect.AssertBody = append(it.tests[it.countTests].Expect.AssertBody, optionalAssertBody(assert))
	}

	return it
}

func (it *cute) AssertHeaders(asserts ...AssertHeaders) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}
	}

	it.tests[it.countTests].Expect.AssertHeaders = append(it.tests[it.countTests].Expect.AssertHeaders, asserts...)

	return it
}

func (it *cute) OptionalAssertHeaders(asserts ...AssertHeaders) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		it.tests[it.countTests].Expect.AssertHeaders = append(it.tests[it.countTests].Expect.AssertHeaders, optionalAssertHeaders(assert))
	}

	return it
}

func (it *cute) AssertResponse(asserts ...AssertResponse) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}
	}

	it.tests[it.countTests].Expect.AssertResponse = append(it.tests[it.countTests].Expect.AssertResponse, asserts...)

	return it
}

func (it *cute) OptionalAssertResponse(asserts ...AssertResponse) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		it.tests[it.countTests].Expect.AssertResponse = append(it.tests[it.countTests].Expect.AssertResponse, optionalAssertResponse(assert))
	}

	return it
}

func (it *cute) AssertBodyT(asserts ...AssertBodyT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}
	}

	it.tests[it.countTests].Expect.AssertBodyT = append(it.tests[it.countTests].Expect.AssertBodyT, asserts...)

	return it
}

func (it *cute) OptionalAssertBodyT(asserts ...AssertBodyT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		it.tests[it.countTests].Expect.AssertBodyT = append(it.tests[it.countTests].Expect.AssertBodyT, optionalAssertBodyT(assert))
	}

	return it
}

func (it *cute) AssertHeadersT(asserts ...AssertHeadersT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}
	}

	it.tests[it.countTests].Expect.AssertHeadersT = append(it.tests[it.countTests].Expect.AssertHeadersT, asserts...)

	return it
}

func (it *cute) OptionalAssertHeadersT(asserts ...AssertHeadersT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		it.tests[it.countTests].Expect.AssertHeadersT = append(it.tests[it.countTests].Expect.AssertHeadersT, optionalAssertHeadersT(assert))
	}

	return it
}

func (it *cute) AssertResponseT(asserts ...AssertResponseT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}
	}

	it.tests[it.countTests].Expect.AssertResponseT = append(it.tests[it.countTests].Expect.AssertResponseT, asserts...)

	return it
}

func (it *cute) OptionalAssertResponseT(asserts ...AssertResponseT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		it.tests[it.countTests].Expect.AssertResponseT = append(it.tests[it.countTests].Expect.AssertResponseT, optionalAssertResponseT(assert))
	}

	return it
}

func (it *cute) EnableHardValidation() ExpectHTTPBuilder {
	it.tests[it.countTests].HardValidation = true

	return it
}

func (it *cute) CreateTableTest() MiddlewareTable {
	it.isTableTest = true

	return it
}

func (it *cute) PutNewTest(name string, r *http.Request, expect *Expect) TableTest {
	// Validate, that first step is empty
	if it.countTests == 0 {
		if it.tests[0].Request.Base == nil &&
			len(it.tests[0].Request.Builders) == 0 {
			it.tests[0].Expect = expect
			it.tests[0].Name = name
			it.tests[0].Request.Base = r

			return it
		}
	}

	newTest := createDefaultTest(it.baseProps)
	newTest.Expect = expect
	newTest.Name = name
	newTest.Request.Base = r
	it.tests = append(it.tests, newTest)
	it.countTests++ // async?

	return it
}

func (it *cute) PutTests(params ...*Test) TableTest {
	for _, param := range params {
		// Validate, that first step is empty
		if it.countTests == 0 {
			if it.tests[0].Request.Base == nil &&
				len(it.tests[0].Request.Builders) == 0 {
				it.tests[0] = param

				continue
			}
		}

		it.tests = append(it.tests, param)
		it.countTests++
	}

	return it
}

func (it *cute) NextTest() NextTestBuilder {
	it.countTests++ // async?

	it.tests = append(it.tests, createDefaultTest(it.baseProps))

	return it
}
