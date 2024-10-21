package cute

import (
	"net/http"
	"time"
)

const defaultHTTPTimeout = 30

var (
	errorAssertIsNil = "assert must be not nil"
)

// HTTPTestMaker is a creator tests
type HTTPTestMaker struct {
	httpClient    *http.Client
	middleware    *Middleware
	jsonMarshaler JSONMarshaler
}

// NewHTTPTestMaker is function for set options for all cute.
// For example, you can set timeout for all requests  or set custom http client
// Options:
// - WithCustomHTTPTimeout - set timeout for all requests
// - WithHTTPClient - set custom http client
// - WithCustomHTTPRoundTripper - set custom http round tripper
// - WithJSONMarshaler - set custom json marshaler
// - WithMiddlewareAfter - set function which will run AFTER test execution
// - WithMiddlewareAfterT - set function which will run AFTER test execution with TB
// - WithMiddlewareBefore - set function which will run BEFORE test execution
// - WithMiddlewareBeforeT - set function which will run BEFORE test execution with TB
func NewHTTPTestMaker(opts ...Option) *HTTPTestMaker {
	var (
		o = &options{
			middleware: new(Middleware),
		}

		timeout                    = defaultHTTPTimeout * time.Second
		roundTripper               = http.DefaultTransport
		jsMarshaler  JSONMarshaler = &jsonMarshaler{}
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

	if o.jsonMarshaler != nil {
		jsMarshaler = o.jsonMarshaler
	}

	m := &HTTPTestMaker{
		httpClient:    httpClient,
		jsonMarshaler: jsMarshaler,
		middleware:    o.middleware,
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
	return &Test{
		httpClient:    m.httpClient,
		jsonMarshaler: m.jsonMarshaler,
		Middleware:    createMiddlewareFromTemplate(m.middleware),
		AllureStep:    new(AllureStep),
		Request: &Request{
			Retry: new(RequestRetryPolitic),
		},
		Expect: &Expect{JSONSchema: new(ExpectJSONSchema)},
	}
}

func createMiddlewareFromTemplate(m *Middleware) *Middleware {
	after := make([]AfterExecute, 0, len(m.After))
	after = append(after, m.After...)

	afterT := make([]AfterExecuteT, 0, len(m.AfterT))
	afterT = append(afterT, m.AfterT...)

	before := make([]BeforeExecute, 0, len(m.Before))
	before = append(before, m.Before...)

	beforeT := make([]BeforeExecuteT, 0, len(m.BeforeT))
	beforeT = append(beforeT, m.BeforeT...)

	middleware := &Middleware{
		After:   after,
		AfterT:  afterT,
		Before:  before,
		BeforeT: beforeT,
	}

	return middleware
}

func (qt *cute) Create() MiddlewareRequest {
	return qt
}

func (qt *cute) CreateStep(name string) MiddlewareRequest {
	qt.tests[qt.countTests].AllureStep.Name = name

	return qt
}

func (qt *cute) CreateRequest() RequestHTTPBuilder {
	return qt
}
