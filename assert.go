package cute

import (
	"net/http"

	"github.com/ozontech/cute/errors"
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

func isOptionError(err error) bool {
	if tErr, ok := err.(errors.OptionalError); ok {
		return tErr.IsOptional()
	}

	return false
}

func optionalAssertHeaders(assert AssertHeaders) AssertHeaders {
	return func(headers http.Header) error {
		err := assert(headers)

		return wrapOptionalError(err)
	}
}

func optionalAssertBody(assert AssertBody) AssertBody {
	return func(body []byte) error {
		err := assert(body)

		return wrapOptionalError(err)
	}
}

func optionalAssertResponse(assert AssertResponse) AssertResponse {
	return func(resp *http.Response) error {
		err := assert(resp)

		return wrapOptionalError(err)
	}
}

func optionalAssertHeadersT(assert AssertHeadersT) AssertHeadersT {
	return func(t T, headers http.Header) error {
		err := assert(t, headers)

		return wrapOptionalError(err)
	}
}

func optionalAssertBodyT(assert AssertBodyT) AssertBodyT {
	return func(t T, body []byte) error {
		err := assert(t, body)

		return wrapOptionalError(err)
	}
}

func optionalAssertResponseT(assert AssertResponseT) AssertResponseT {
	return func(t T, resp *http.Response) error {
		err := assert(t, resp)

		return wrapOptionalError(err)
	}
}

func wrapOptionalError(err error) error {
	if err == nil {
		return nil
	}

	if tErr, ok := err.(errors.OptionalError); ok {
		tErr.SetOptional(true)

		return tErr.(error)
	}

	return errors.NewOptionalError(err.Error())
}
