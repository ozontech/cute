package cute

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ozontech/allure-go/pkg/allure"
)

const defaultHTTPTimeout = 30

var (
	errorAssertIsNil = "assert must be not nil"
)

// HTTPTestMaker is a creator tests
type HTTPTestMaker struct {
	httpClient *http.Client
	middleware *Middleware

	// todo add marshaler
}

type options struct {
	httpClient       *http.Client
	httpTimeout      time.Duration
	httpRoundTripper http.RoundTripper

	middleware *Middleware
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
		httpClient: httpClient,
		middleware: o.middleware,
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
		httpClient: m.httpClient,
		Middleware: middleware,
		AllureStep: new(AllureStep),
		Request: &Request{
			Repeat: new(RequestRepeatPolitic),
		},
		Expect: &Expect{JSONSchema: new(ExpectJSONSchema)},
	}
}

func (qt *cute) Title(title string) AllureBuilder {
	qt.allureInfo.title = title

	return qt
}

func (qt *cute) Epic(epic string) AllureBuilder {
	qt.allureLabels.epic = epic

	return qt
}

func (qt *cute) Titlef(format string, args ...interface{}) AllureBuilder {
	qt.allureInfo.title = fmt.Sprintf(format, args...)

	return qt
}

func (qt *cute) Descriptionf(format string, args ...interface{}) AllureBuilder {
	qt.allureInfo.description = fmt.Sprintf(format, args...)

	return qt
}

func (qt *cute) Stage(stage string) AllureBuilder {
	qt.allureInfo.stage = stage

	return qt
}

func (qt *cute) Stagef(format string, args ...interface{}) AllureBuilder {
	qt.allureInfo.stage = fmt.Sprintf(format, args...)

	return qt
}

func (qt *cute) Layer(value string) AllureBuilder {
	qt.allureLabels.layer = value

	return qt
}

func (qt *cute) TmsLink(tmsLink string) AllureBuilder {
	qt.allureLinks.tmsLink = tmsLink

	return qt
}

func (qt *cute) TmsLinks(tmsLinks ...string) AllureBuilder {
	qt.allureLinks.tmsLinks = append(qt.allureLinks.tmsLinks, tmsLinks...)

	return qt
}

func (qt *cute) SetIssue(issue string) AllureBuilder {
	qt.allureLinks.issue = issue

	return qt
}

func (qt *cute) SetTestCase(testCase string) AllureBuilder {
	qt.allureLinks.testCase = testCase

	return qt
}

func (qt *cute) Link(link *allure.Link) AllureBuilder {
	qt.allureLinks.link = link

	return qt
}

func (qt *cute) ID(value string) AllureBuilder {
	qt.allureLabels.id = value

	return qt
}

func (qt *cute) AllureID(value string) AllureBuilder {
	qt.allureLabels.allureID = value

	return qt
}

func (qt *cute) AddSuiteLabel(value string) AllureBuilder {
	qt.allureLabels.suiteLabel = value

	return qt
}

func (qt *cute) AddSubSuite(value string) AllureBuilder {
	qt.allureLabels.subSuite = value

	return qt
}

func (qt *cute) AddParentSuite(value string) AllureBuilder {
	qt.allureLabels.parentSuite = value

	return qt
}

func (qt *cute) Story(value string) AllureBuilder {
	qt.allureLabels.story = value

	return qt
}

func (qt *cute) Tag(value string) AllureBuilder {
	qt.allureLabels.tag = value

	return qt
}

func (qt *cute) Severity(value allure.SeverityType) AllureBuilder {
	qt.allureLabels.severity = value

	return qt
}

func (qt *cute) Owner(value string) AllureBuilder {
	qt.allureLabels.owner = value

	return qt
}

func (qt *cute) Lead(value string) AllureBuilder {
	qt.allureLabels.lead = value

	return qt
}

func (qt *cute) Label(label *allure.Label) AllureBuilder {
	qt.allureLabels.label = label

	return qt
}

func (qt *cute) Labels(labels ...*allure.Label) AllureBuilder {
	qt.allureLabels.labels = labels

	return qt
}

func (qt *cute) Description(description string) AllureBuilder {
	qt.allureInfo.description = description

	return qt
}

func (qt *cute) Tags(tags ...string) AllureBuilder {
	qt.allureLabels.tags = tags

	return qt
}

func (qt *cute) Feature(feature string) AllureBuilder {
	qt.allureLabels.feature = feature

	return qt
}

func (qt *cute) Create() MiddlewareRequest {
	return qt
}

func (qt *cute) CreateStep(name string) MiddlewareRequest {
	qt.tests[qt.countTests].AllureStep.Name = name

	return qt
}

func (qt *cute) Parallel() AllureBuilder {
	qt.parallel = true

	return qt
}

func (qt *cute) CreateRequest() RequestHTTPBuilder {
	return qt
}

func (qt *cute) StepName(name string) MiddlewareRequest {
	qt.tests[qt.countTests].AllureStep.Name = name

	return qt
}

func (qt *cute) BeforeExecute(fs ...BeforeExecute) MiddlewareRequest {
	qt.tests[qt.countTests].Middleware.Before = append(qt.tests[qt.countTests].Middleware.Before, fs...)

	return qt
}

func (qt *cute) BeforeExecuteT(fs ...BeforeExecuteT) MiddlewareRequest {
	qt.tests[qt.countTests].Middleware.BeforeT = append(qt.tests[qt.countTests].Middleware.BeforeT, fs...)

	return qt
}

func (qt *cute) After(fs ...AfterExecute) ExpectHTTPBuilder {
	qt.tests[qt.countTests].Middleware.After = append(qt.tests[qt.countTests].Middleware.After, fs...)

	return qt
}

func (qt *cute) AfterT(fs ...AfterExecuteT) ExpectHTTPBuilder {
	qt.tests[qt.countTests].Middleware.AfterT = append(qt.tests[qt.countTests].Middleware.AfterT, fs...)

	return qt
}

func (qt *cute) AfterExecute(fs ...AfterExecute) MiddlewareRequest {
	qt.tests[qt.countTests].Middleware.After = append(qt.tests[qt.countTests].Middleware.After, fs...)

	return qt
}

func (qt *cute) AfterExecuteT(fs ...AfterExecuteT) MiddlewareRequest {
	qt.tests[qt.countTests].Middleware.AfterT = append(qt.tests[qt.countTests].Middleware.AfterT, fs...)

	return qt
}

func (qt *cute) AfterTestExecute(fs ...AfterExecute) NextTestBuilder {
	previousTest := 0
	if qt.countTests != 0 {
		previousTest = qt.countTests - 1
	}

	qt.tests[previousTest].Middleware.After = append(qt.tests[previousTest].Middleware.After, fs...)

	return qt
}

func (qt *cute) AfterTestExecuteT(fs ...AfterExecuteT) NextTestBuilder {
	previousTest := 0
	if qt.countTests != 0 {
		previousTest = qt.countTests - 1
	}

	qt.tests[previousTest].Middleware.AfterT = append(qt.tests[previousTest].Middleware.AfterT, fs...)

	return qt
}

func (qt *cute) RequestRepeat(count int) RequestHTTPBuilder {
	qt.tests[qt.countTests].Request.Repeat.Count = count

	return qt
}

func (qt *cute) RequestRepeatDelay(delay time.Duration) RequestHTTPBuilder {
	qt.tests[qt.countTests].Request.Repeat.Delay = delay

	return qt
}

func (qt *cute) Request(r *http.Request) ExpectHTTPBuilder {
	qt.tests[qt.countTests].Request.Base = r

	return qt
}

func (qt *cute) RequestBuilder(r ...RequestBuilder) ExpectHTTPBuilder {
	qt.tests[qt.countTests].Request.Builders = append(qt.tests[qt.countTests].Request.Builders, r...)

	return qt
}

func (qt *cute) ExpectExecuteTimeout(t time.Duration) ExpectHTTPBuilder {
	qt.tests[qt.countTests].Expect.ExecuteTime = t

	return qt
}

func (qt *cute) ExpectStatus(code int) ExpectHTTPBuilder {
	qt.tests[qt.countTests].Expect.Code = code

	return qt
}

func (qt *cute) ExpectJSONSchemaString(schema string) ExpectHTTPBuilder {
	qt.tests[qt.countTests].Expect.JSONSchema.String = schema

	return qt
}

func (qt *cute) ExpectJSONSchemaByte(schema []byte) ExpectHTTPBuilder {
	qt.tests[qt.countTests].Expect.JSONSchema.Byte = schema

	return qt
}

func (qt *cute) ExpectJSONSchemaFile(filePath string) ExpectHTTPBuilder {
	qt.tests[qt.countTests].Expect.JSONSchema.File = filePath

	return qt
}

func (qt *cute) AssertBody(asserts ...AssertBody) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}
	}

	qt.tests[qt.countTests].Expect.AssertBody = append(qt.tests[qt.countTests].Expect.AssertBody, asserts...)

	return qt
}

func (qt *cute) OptionalAssertBody(asserts ...AssertBody) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertBody = append(qt.tests[qt.countTests].Expect.AssertBody, optionalAssertBody(assert))
	}

	return qt
}

func (qt *cute) RequireBody(asserts ...AssertBody) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertBody = append(qt.tests[qt.countTests].Expect.AssertBody, requireAssertBody(assert))
	}

	return qt
}

func (qt *cute) AssertHeaders(asserts ...AssertHeaders) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}
	}

	qt.tests[qt.countTests].Expect.AssertHeaders = append(qt.tests[qt.countTests].Expect.AssertHeaders, asserts...)

	return qt
}

func (qt *cute) OptionalAssertHeaders(asserts ...AssertHeaders) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertHeaders = append(qt.tests[qt.countTests].Expect.AssertHeaders, optionalAssertHeaders(assert))
	}

	return qt
}

func (qt *cute) RequireHeaders(asserts ...AssertHeaders) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertHeaders = append(qt.tests[qt.countTests].Expect.AssertHeaders, requireAssertHeaders(assert))
	}

	return qt
}

func (qt *cute) AssertResponse(asserts ...AssertResponse) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}
	}

	qt.tests[qt.countTests].Expect.AssertResponse = append(qt.tests[qt.countTests].Expect.AssertResponse, asserts...)

	return qt
}

func (qt *cute) OptionalAssertResponse(asserts ...AssertResponse) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertResponse = append(qt.tests[qt.countTests].Expect.AssertResponse, optionalAssertResponse(assert))
	}

	return qt
}

func (qt *cute) RequireResponse(asserts ...AssertResponse) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertResponse = append(qt.tests[qt.countTests].Expect.AssertResponse, requireAssertResponse(assert))
	}

	return qt
}

func (qt *cute) AssertBodyT(asserts ...AssertBodyT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}
	}

	qt.tests[qt.countTests].Expect.AssertBodyT = append(qt.tests[qt.countTests].Expect.AssertBodyT, asserts...)

	return qt
}

func (qt *cute) OptionalAssertBodyT(asserts ...AssertBodyT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertBodyT = append(qt.tests[qt.countTests].Expect.AssertBodyT, optionalAssertBodyT(assert))
	}

	return qt
}

func (qt *cute) RequireBodyT(asserts ...AssertBodyT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertBodyT = append(qt.tests[qt.countTests].Expect.AssertBodyT, requireAssertBodyT(assert))
	}

	return qt
}

func (qt *cute) AssertHeadersT(asserts ...AssertHeadersT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}
	}

	qt.tests[qt.countTests].Expect.AssertHeadersT = append(qt.tests[qt.countTests].Expect.AssertHeadersT, asserts...)

	return qt
}

func (qt *cute) OptionalAssertHeadersT(asserts ...AssertHeadersT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertHeadersT = append(qt.tests[qt.countTests].Expect.AssertHeadersT, optionalAssertHeadersT(assert))
	}

	return qt
}

func (qt *cute) RequireHeadersT(asserts ...AssertHeadersT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertHeadersT = append(qt.tests[qt.countTests].Expect.AssertHeadersT, requireAssertHeadersT(assert))
	}

	return qt
}

func (qt *cute) AssertResponseT(asserts ...AssertResponseT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}
	}

	qt.tests[qt.countTests].Expect.AssertResponseT = append(qt.tests[qt.countTests].Expect.AssertResponseT, asserts...)

	return qt
}

func (qt *cute) OptionalAssertResponseT(asserts ...AssertResponseT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertResponseT = append(qt.tests[qt.countTests].Expect.AssertResponseT, optionalAssertResponseT(assert))
	}

	return qt
}

func (qt *cute) RequireResponseT(asserts ...AssertResponseT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		if assert == nil {
			panic(errorAssertIsNil)
		}

		qt.tests[qt.countTests].Expect.AssertResponseT = append(qt.tests[qt.countTests].Expect.AssertResponseT, requireAssertResponseT(assert))
	}

	return qt
}

func (qt *cute) CreateTableTest() MiddlewareTable {
	qt.isTableTest = true

	return qt
}

func (qt *cute) PutNewTest(name string, r *http.Request, expect *Expect) TableTest {
	// Validate, that first step is empty
	if qt.countTests == 0 {
		if qt.tests[0].Request.Base == nil &&
			len(qt.tests[0].Request.Builders) == 0 {
			qt.tests[0].Expect = expect
			qt.tests[0].Name = name
			qt.tests[0].Request.Base = r

			return qt
		}
	}

	newTest := createDefaultTest(qt.baseProps)
	newTest.Expect = expect
	newTest.Name = name
	newTest.Request.Base = r
	qt.tests = append(qt.tests, newTest)
	qt.countTests++ // async?

	return qt
}

func (qt *cute) PutTests(params ...*Test) TableTest {
	for _, param := range params {
		// Validate, that first step is empty
		if qt.countTests == 0 {
			if qt.tests[0].Request.Base == nil &&
				len(qt.tests[0].Request.Builders) == 0 {
				qt.tests[0] = param

				continue
			}
		}

		qt.tests = append(qt.tests, param)
		qt.countTests++
	}

	return qt
}

func (qt *cute) NextTest() NextTestBuilder {
	qt.countTests++ // async?

	qt.tests = append(qt.tests, createDefaultTest(qt.baseProps))

	return qt
}
