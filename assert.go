package cute

import (
	"fmt"
	"net/http"
	"runtime"
)

// This is type of asserts, for create some assert with using custom logic.

// AssertBody ...
type AssertBody func(body []byte) error

// AssertHeaders ...
type AssertHeaders func(headers http.Header) error

// AssertResponse ...
type AssertResponse func(response *http.Response) error

// This is type for create custom assertions with using allure and testing.allureProvider

// AssertBodyT ...
type AssertBodyT func(t T, body []byte) error

// AssertHeadersT ...
type AssertHeadersT func(t T, headers http.Header) error

// AssertResponseT ...
type AssertResponseT func(t T, response *http.Response) error

func callAssertWrapper(assert AssertBody, withCall string) func(body []byte) error {
	return func(body []byte) error {
		err := assert(body)

		if err != nil {
			return fmt.Errorf("%w\n%s", err, withCall)
		}

		return nil
	}
}

func getCaller() string {
	pcs := make([]uintptr, 10)
	depth := runtime.Callers(3, pcs)
	if depth == 0 {
		fmt.Println("Couldn't get the stack information")
		return ""
	}
	callers := runtime.CallersFrames(pcs[:depth])
	caller, _ := callers.Next()

	return fmt.Sprintf("Called from: %s:%d\n", caller.File, caller.Line)
}

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
