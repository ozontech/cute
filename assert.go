package cute

import (
	"net/http"
)

// This is type of asserts, for create some assert with using custom logic.

// AssertBody is type for create custom assertions for body
// Example asserts:
// - json.LengthGreaterThan
// - json.LengthGreaterOrEqualThan
// - json.LengthLessThan
// - json.LengthLessOrEqualThan
// - json.Present
// - json.NotEmpty
// - json.NotPresent
type AssertBody func(body []byte) error

// AssertHeaders is type for create custom assertions for headers
// Example asserts:
// - headers.Present
// - headers.NotPresent
type AssertHeaders func(headers http.Header) error

// AssertResponse is type for create custom assertions for response
type AssertResponse func(response *http.Response) error

// This is type for create custom assertions with using allure and testing.allureProvider

// AssertBodyT is type for create custom assertions for body with TB
// Check example in AssertBody
// TB is testing.T and it can be used for require ore assert from testify or another packages
type AssertBodyT func(t T, body []byte) error

// AssertHeadersT is type for create custom assertions for headers with TB
// Check example in AssertHeaders
// TB is testing.T and it can be used for require ore assert from testify or another packages
type AssertHeadersT func(t T, headers http.Header) error

// AssertResponseT is type for create custom assertions for response with TB
// Check example in AssertResponse
// TB is testing.T and it can be used for require ore assert from testify or another packages
type AssertResponseT func(t T, response *http.Response) error

func (it *Test) assertHeaders(t internalT, headers http.Header) []error {
	var (
		asserts = it.Expect.AssertHeaders
		assertT = it.Expect.AssertHeadersT
	)

	if len(asserts) == 0 && len(assertT) == 0 {
		return nil
	}

	return executeWithStep(t, "Assert headers", func(t T) []error {
		errs := make([]error, 0)
		// Execute assert only response
		for _, f := range asserts {
			err := f(headers)
			if err != nil {
				errs = append(errs, err)
			}
		}

		// Execute assert for response with TB
		for _, f := range assertT {
			err := f(t, headers)
			if err != nil {
				errs = append(errs, err)
			}
		}

		return errs
	})
}

func (it *Test) assertResponse(t internalT, resp *http.Response) []error {
	var (
		asserts = it.Expect.AssertResponse
		assertT = it.Expect.AssertResponseT
	)

	if len(asserts) == 0 && len(assertT) == 0 {
		return nil
	}

	return executeWithStep(t, "Assert response", func(t T) []error {
		errs := make([]error, 0)
		// Execute assert only response
		for _, f := range asserts {
			err := f(resp)
			if err != nil {
				errs = append(errs, err)
			}
		}

		// Execute assert for response with TB
		for _, f := range assertT {
			err := f(t, resp)
			if err != nil {
				errs = append(errs, err)
			}
		}

		return errs
	})
}

func (it *Test) assertBody(t internalT, body []byte) []error {
	var (
		asserts = it.Expect.AssertBody
		assertT = it.Expect.AssertBodyT
	)

	if len(asserts) == 0 && len(assertT) == 0 {
		return nil
	}

	return executeWithStep(t, "Assert body", func(t T) []error {
		errs := make([]error, 0)
		// Execute assert only response
		for _, f := range asserts {
			err := f(body)
			if err != nil {
				errs = append(errs, err)
			}
		}

		// Execute assert for response with TB
		for _, f := range assertT {
			err := f(t, body)
			if err != nil {
				errs = append(errs, err)
			}
		}

		return errs
	})
}
