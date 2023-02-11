package cute

import (
	"net/http"
)

type testResults struct {
	isFailed bool
	name     string
	resp     *http.Response
	errors   []error
}

func newTestResult(name string, resp *http.Response, isFailed bool, errs []error) ResultsHTTPBuilder {
	return &testResults{
		name:     name,
		resp:     resp,
		isFailed: isFailed,
		errors:   errs,
	}
}

func (r *testResults) GetHTTPResponse() *http.Response {
	return r.resp
}

func (r *testResults) GetErrors() []error {
	return r.errors
}

func (r *testResults) GetName() string {
	return r.name
}

func (r *testResults) IsFailed() bool {
	return r.isFailed
}
