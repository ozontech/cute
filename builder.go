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
		httpClient:    m.httpClient,
		jsonMarshaler: m.jsonMarshaler,
		Middleware:    middleware,
		AllureStep:    new(AllureStep),
		Request: &Request{
			Repeat: new(RequestRepeatPolitic),
		},
		Expect: &Expect{JSONSchema: new(ExpectJSONSchema)},
	}
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
