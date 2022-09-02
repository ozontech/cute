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

// NewHTTPTestMaker is function for set options for all test.
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

// NewTestBuilder is a function for initialization foundation for test
func (m *HTTPTestMaker) NewTestBuilder() AllureBuilder {
	return &test{
		httpClient:   m.httpClient,
		allureStep:   new(allureStep),
		allureInfo:   new(allureInformation),
		allureLinks:  new(allureLinks),
		allureLabels: new(allureLabels),
		middleware:   new(middleware),
		request: &request{
			repeat: new(requestRepeatPolitic),
		},
		expect:   new(expect),
		parallel: false,
	}
}

func (it *test) Title(title string) AllureBuilder {
	it.allureInfo.title = title

	return it
}

func (it *test) Epic(epic string) AllureBuilder {
	it.allureLabels.epic = epic

	return it
}

func (it *test) SetIssue(issue string) AllureBuilder {
	it.allureLinks.issue = issue

	return it
}

func (it *test) SetTestCase(testCase string) AllureBuilder {
	it.allureLinks.testCase = testCase

	return it
}

func (it *test) Link(link *allure.Link) AllureBuilder {
	it.allureLinks.link = link

	return it
}

func (it *test) ID(value string) AllureBuilder {
	it.allureLabels.id = value

	return it
}

func (it *test) AllureID(value string) AllureBuilder {
	it.allureLabels.allureID = value

	return it
}

func (it *test) AddSuiteLabel(value string) AllureBuilder {
	it.allureLabels.suiteLabel = value

	return it
}

func (it *test) AddSubSuite(value string) AllureBuilder {
	it.allureLabels.subSuite = value

	return it
}

func (it *test) AddParentSuite(value string) AllureBuilder {
	it.allureLabels.parentSuite = value

	return it
}

func (it *test) Story(value string) AllureBuilder {
	it.allureLabels.story = value

	return it
}

func (it *test) Tag(value string) AllureBuilder {
	it.allureLabels.tag = value

	return it
}

func (it *test) Severity(value allure.SeverityType) AllureBuilder {
	it.allureLabels.severity = value

	return it
}

func (it *test) Owner(value string) AllureBuilder {
	it.allureLabels.owner = value

	return it
}

func (it *test) Lead(value string) AllureBuilder {
	it.allureLabels.lead = value

	return it
}

func (it *test) Label(label *allure.Label) AllureBuilder {
	it.allureLabels.label = label

	return it
}

func (it *test) Labels(labels ...*allure.Label) AllureBuilder {
	it.allureLabels.labels = labels

	return it
}

func (it *test) Description(description string) AllureBuilder {
	it.allureInfo.description = description

	return it
}

func (it *test) Tags(tags ...string) AllureBuilder {
	it.allureLabels.tags = tags

	return it
}

func (it *test) Feature(feature string) AllureBuilder {
	it.allureLabels.feature = feature

	return it
}

func (it *test) CreateWithStep() StepBuilder {
	return it
}

func (it *test) Create() Middleware {
	return it
}

func (it *test) Parallel() AllureBuilder {
	it.parallel = true

	return it
}

func (it *test) CreateRequest() RequestHTTPBuilder {
	return it
}

func (it *test) StepName(name string) Middleware {
	it.allureStep.name = name

	return it
}

func (it *test) BeforeExecute(fs ...BeforeExecute) Middleware {
	it.middleware.before = append(it.middleware.before, fs...)

	return it
}

func (it *test) BeforeExecuteT(fs ...BeforeExecuteT) Middleware {
	it.middleware.beforeT = append(it.middleware.beforeT, fs...)

	return it
}

func (it *test) AfterExecute(fs ...AfterExecute) Middleware {
	it.middleware.after = append(it.middleware.after, fs...)

	return it
}

func (it *test) AfterExecuteT(fs ...AfterExecuteT) Middleware {
	it.middleware.afterT = append(it.middleware.afterT, fs...)

	return it
}

func (it *test) RequestRepeat(count int) RequestHTTPBuilder {
	it.request.repeat.count = count

	return it
}

func (it *test) RequestRepeatDelay(delay time.Duration) RequestHTTPBuilder {
	it.request.repeat.delay = delay

	return it
}

func (it *test) Request(r *http.Request) ExpectHTTPBuilder {
	it.request.base = r

	return it
}

func (it *test) RequestBuilder(r ...requestBuilder) ExpectHTTPBuilder {
	it.request.builders = append(it.request.builders, r...)

	return it
}

func (it *test) ExpectExecuteTimeout(t time.Duration) ExpectHTTPBuilder {
	it.expect.executeTime = t

	return it
}

func (it *test) ExpectStatus(code int) ExpectHTTPBuilder {
	it.expect.code = code

	return it
}

func (it *test) ExpectJSONSchemaString(schema string) ExpectHTTPBuilder {
	it.expect.jsSchemaString = schema

	return it
}

func (it *test) ExpectJSONSchemaByte(schema []byte) ExpectHTTPBuilder {
	it.expect.jsSchemaByte = schema

	return it
}

func (it *test) ExpectJSONSchemaFile(filePath string) ExpectHTTPBuilder {
	it.expect.jsSchemaFile = filePath

	return it
}

func (it *test) AssertBody(asserts ...AssertBody) ExpectHTTPBuilder {
	it.expect.assertBody = append(it.expect.assertBody, asserts...)

	return it
}

func (it *test) OptionalAssertBody(asserts ...AssertBody) ExpectHTTPBuilder {
	for _, assert := range asserts {
		it.expect.assertBody = append(it.expect.assertBody, optionalAssertBody(assert))
	}

	return it
}

func (it *test) AssertHeaders(asserts ...AssertHeaders) ExpectHTTPBuilder {
	it.expect.assertHeaders = append(it.expect.assertHeaders, asserts...)

	return it
}

func (it *test) OptionalAssertHeaders(asserts ...AssertHeaders) ExpectHTTPBuilder {
	for _, assert := range asserts {
		it.expect.assertHeaders = append(it.expect.assertHeaders, optionalAssertHeaders(assert))
	}

	return it
}

func (it *test) AssertResponse(asserts ...AssertResponse) ExpectHTTPBuilder {
	it.expect.assertResponse = append(it.expect.assertResponse, asserts...)

	return it
}

func (it *test) OptionalAssertResponse(asserts ...AssertResponse) ExpectHTTPBuilder {
	for _, assert := range asserts {
		it.expect.assertResponse = append(it.expect.assertResponse, optionalAssertResponse(assert))
	}

	return it
}

func (it *test) AssertBodyT(asserts ...AssertBodyT) ExpectHTTPBuilder {
	it.expect.assertBodyT = append(it.expect.assertBodyT, asserts...)

	return it
}

func (it *test) OptionalAssertBodyT(asserts ...AssertBodyT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		it.expect.assertBodyT = append(it.expect.assertBodyT, optionalAssertBodyT(assert))
	}

	return it
}

func (it *test) AssertHeadersT(asserts ...AssertHeadersT) ExpectHTTPBuilder {
	it.expect.assertHeadersT = append(it.expect.assertHeadersT, asserts...)

	return it
}

func (it *test) OptionalAssertHeadersT(asserts ...AssertHeadersT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		it.expect.assertHeadersT = append(it.expect.assertHeadersT, optionalAssertHeadersT(assert))
	}

	return it
}

func (it *test) AssertResponseT(asserts ...AssertResponseT) ExpectHTTPBuilder {
	it.expect.assertResponseT = append(it.expect.assertResponseT, asserts...)

	return it
}

func (it *test) OptionalAssertResponseT(asserts ...AssertResponseT) ExpectHTTPBuilder {
	for _, assert := range asserts {
		it.expect.assertResponseT = append(it.expect.assertResponseT, optionalAssertResponseT(assert))
	}

	return it
}
