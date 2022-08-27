package cute

import (
	"net/http"
	"time"

	"github.com/ozontech/allure-go/pkg/allure"
)

const defaultHTTPTimeout = 30

// HTTPTestMaker is a creator tests
type HTTPTestMaker struct {
	httpClient *http.Client
	// todo add marshaler
}

type options struct {
	httpClient       *http.Client
	httpTimeout      time.Duration
	httpRoundTripper http.RoundTripper
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

// NewHTTPTestMaker is function for set options for all cute.
func NewHTTPTestMaker(opts ...Option) *HTTPTestMaker {
	var (
		o            = new(options)
		timeout      = defaultHTTPTimeout * time.Second
		roundTripper = http.DefaultTransport
	)

	for _, opt := range opts {
		opt(o)
	}

	if o.httpTimeout != 0 {
		timeout = o.httpTimeout
	}

	if o.httpRoundTripper != nil {
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
	}

	return m
}

// NewTestBuilder is a function for initialization foundation for cute
func (m *HTTPTestMaker) NewTestBuilder() AllureBuilder {
	tests := make([]*test, 1, 1)
	tests[0] = createDefaultTest()

	return &cute{
		httpClient:   m.httpClient,
		countTests:   0,
		tests:        tests,
		allureInfo:   new(allureInformation),
		allureLinks:  new(allureLinks),
		allureLabels: new(allureLabels),
		parallel:     false,
	}
}

func createDefaultTest() *test {
	return &test{
		allureStep: new(allureStep),
		middleware: new(middleware),
		request: &request{
			repeat: new(requestRepeatPolitic),
		},
		expect: new(Expect),
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

func (it *cute) Link(link allure.Link) AllureBuilder {
	it.allureLinks.link = link

	return it
}

func (it *cute) ID(value string) AllureBuilder {
	it.allureLabels.id = value

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

func (it *cute) Label(label allure.Label) AllureBuilder {
	it.allureLabels.label = label

	return it
}

func (it *cute) Labels(labels ...allure.Label) AllureBuilder {
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

func (it *cute) CreateWithStep() StepBuilder {
	return it
}

func (it *cute) Create() Middleware {
	return it
}

func (it *cute) CreateStep(name string) Middleware {
	it.tests[it.countTests].allureStep.name = name

	return it
}

func (it *cute) Parallel() AllureBuilder {
	it.parallel = true

	return it
}

func (it *cute) CreateRequest() RequestHTTPBuilder {
	return it
}

func (it *cute) StepName(name string) Middleware {
	it.tests[it.countTests].allureStep.name = name

	return it
}

func (it *cute) BeforeExecute(fs ...BeforeExecute) Middleware {
	it.tests[it.countTests].middleware.before = append(it.tests[it.countTests].middleware.before, fs...)

	return it
}

func (it *cute) BeforeExecuteT(fs ...BeforeExecuteT) Middleware {
	it.tests[it.countTests].middleware.beforeT = append(it.tests[it.countTests].middleware.beforeT, fs...)

	return it
}

func (it *cute) AfterExecute(fs ...AfterExecute) Middleware {
	it.tests[it.countTests].middleware.after = append(it.tests[it.countTests].middleware.after, fs...)

	return it
}

func (it *cute) AfterExecuteT(fs ...AfterExecuteT) Middleware {
	it.tests[it.countTests].middleware.afterT = append(it.tests[it.countTests].middleware.afterT, fs...)

	return it
}

func (it *cute) RequestRepeat(count int) RequestHTTPBuilder {
	it.tests[it.countTests].request.repeat.count = count

	return it
}

func (it *cute) RequestRepeatDelay(delay time.Duration) RequestHTTPBuilder {
	it.tests[it.countTests].request.repeat.delay = delay

	return it
}

func (it *cute) Request(r *http.Request) ExpectHTTPBuilder {
	it.tests[it.countTests].request.base = r

	return it
}

func (it *cute) RequestBuilder(r ...requestBuilder) ExpectHTTPBuilder {
	it.tests[it.countTests].request.builders = append(it.tests[it.countTests].request.builders, r...)

	return it
}

func (it *cute) ExpectExecuteTimeout(t time.Duration) ExpectHTTPBuilder {
	it.tests[it.countTests].expect.ExecuteTime = t

	return it
}

func (it *cute) ExpectStatus(code int) ExpectHTTPBuilder {
	it.tests[it.countTests].expect.Code = code

	return it
}

func (it *cute) ExpectJSONSchemaString(schema string) ExpectHTTPBuilder {
	it.tests[it.countTests].expect.JSONSchemaString = schema

	return it
}

func (it *cute) ExpectJSONSchemaByte(schema []byte) ExpectHTTPBuilder {
	it.tests[it.countTests].expect.JSONSchemaByte = schema

	return it
}

func (it *cute) ExpectJSONSchemaFile(filePath string) ExpectHTTPBuilder {
	it.tests[it.countTests].expect.JSONSchemaFile = filePath

	return it
}

func (it *cute) AssertBody(asserts ...AssertBody) ExpectHTTPBuilder {
	it.tests[it.countTests].expect.AssertBody = append(it.tests[it.countTests].expect.AssertBody, asserts...)

	return it
}

func (it *cute) OptionalAssertBody(asserts ...AssertBody) ExpectHTTPBuilder {
	for _, assert := range asserts {
		it.tests[it.countTests].expect.AssertBody = append(it.tests[it.countTests].expect.AssertBody, optionalAssertBody(assert))
	}

	return it
}

func (it *cute) AssertHeaders(asserts ...AssertHeaders) ExpectHTTPBuilder {
	it.tests[it.countTests].expect.AssertHeaders = append(it.tests[it.countTests].expect.AssertHeaders, asserts...)

	return it
}

func (it *cute) OptionalAssertHeaders(asserts ...AssertHeaders) ExpectHTTPBuilder {
	for _, assert := range asserts {
		it.tests[it.countTests].expect.AssertHeaders = append(it.tests[it.countTests].expect.AssertHeaders, optionalAssertHeaders(assert))
	}

	return it
}

func (it *cute) AssertResponse(asserts ...AssertResponse) ExpectHTTPBuilder {
	it.tests[it.countTests].expect.AssertResponse = append(it.tests[it.countTests].expect.AssertResponse, asserts...)

	return it
}

func (it *cute) OptionalAssertResponse(asserts ...AssertResponse) ExpectHTTPBuilder {
	for _, assert := range asserts {
		it.tests[it.countTests].expect.AssertResponse = append(it.tests[it.countTests].expect.AssertResponse, optionalAssertResponse(assert))
	}

	return it
}

func (it *cute) AssertBodyT(asserts ...AssertBodyT) ExpectHTTPBuilder {
	it.tests[it.countTests].expect.AssertBodyT = append(it.tests[it.countTests].expect.AssertBodyT, asserts...)

	return it
}

func (it *cute) OptionalAssertBodyT(asserts ...AssertBodyT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		it.tests[it.countTests].expect.AssertBodyT = append(it.tests[it.countTests].expect.AssertBodyT, optionalAssertBodyT(assert))
	}

	return it
}

func (it *cute) AssertHeadersT(asserts ...AssertHeadersT) ExpectHTTPBuilder {
	it.tests[it.countTests].expect.AssertHeadersT = append(it.tests[it.countTests].expect.AssertHeadersT, asserts...)

	return it
}

func (it *cute) OptionalAssertHeadersT(asserts ...AssertHeadersT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		it.tests[it.countTests].expect.AssertHeadersT = append(it.tests[it.countTests].expect.AssertHeadersT, optionalAssertHeadersT(assert))
	}

	return it
}

func (it *cute) AssertResponseT(asserts ...AssertResponseT) ExpectHTTPBuilder {
	it.tests[it.countTests].expect.AssertResponseT = append(it.tests[it.countTests].expect.AssertResponseT, asserts...)

	return it
}

func (it *cute) OptionalAssertResponseT(asserts ...AssertResponseT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		it.tests[it.countTests].expect.AssertResponseT = append(it.tests[it.countTests].expect.AssertResponseT, optionalAssertResponseT(assert))
	}

	return it
}

func (it *cute) CreateTableTest() MiddlewareTable {
	it.isTableTest = true

	return it
}

func (it *cute) PutTest(name string, r *http.Request, expect *Expect) TableTest {
	// Validate, that first step is empty
	if it.countTests == 0 {
		if it.tests[0].request.base == nil &&
			len(it.tests[0].request.builders) == 0 {
			it.tests[0].expect = expect
			it.tests[0].name = name
			it.tests[0].request.base = r

			return it
		}
	}

	newTest := createDefaultTest()
	newTest.expect = expect
	newTest.name = name
	newTest.request.base = r
	it.tests = append(it.tests, newTest)
	it.countTests++ // async?

	return it
}

func (it *cute) PutTests(params ...*TableTestParam) TableTest {
	for _, param := range params {
		it.PutTest(param.Name, param.Request, param.Expect)
	}

	return it
}

func (it *cute) NextTest() NextTestBuilder {
	it.countTests++ // async?

	it.tests = append(it.tests, createDefaultTest())

	return it
}
