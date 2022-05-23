package cute

import (
	"net/http"
)

type testResults struct {
	resp     *http.Response
	errors   []error
	httpTest *test
}

func (r *testResults) GetHTTPResponse() *http.Response {
	return r.resp
}

func (r *testResults) GetErrors() []error {
	return r.errors
}

func (r *testResults) NextTest() Middleware {
	return &test{
		httpClient:   r.httpTest.httpClient,
		allureInfo:   r.httpTest.allureInfo,
		allureLinks:  r.httpTest.allureLinks,
		allureLabels: r.httpTest.allureLabels,
		parallel:     r.httpTest.parallel,
		allureStep:   new(allureStep),
		middleware:   new(middleware),
		request:      new(request),
		expect:       new(expect),
	}
}

func (r *testResults) NextTestWithStep() StepBuilder {
	return r.NextTest().(StepBuilder)
}
