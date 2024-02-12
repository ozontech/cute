package cute

import (
	"net/http"
)

// ResultState is state of test
type ResultState int

// ResultState is state of test
const (
	ResultStateSuccess ResultState = iota
	ResultStateBroken
	ResultStateFail
	// ResultStateFailNow is state for require validations
	ResultStateFailNow
)

type testResults struct {
	name   string
	state  ResultState
	resp   *http.Response
	errors []error
}

func newTestResult(name string, resp *http.Response, state ResultState, errs []error) ResultsHTTPBuilder {
	return &testResults{
		name:   name,
		resp:   resp,
		state:  state,
		errors: errs,
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

// IsFailed ...
// Deprecated please use GetResultState
func (r *testResults) IsFailed() bool {
	return r.state == ResultStateFail
}

func (r *testResults) GetResultState() ResultState {
	return r.state
}
