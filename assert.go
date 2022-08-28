package cute

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
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
		assertHeaders  = it.Expect.AssertHeaders
		assertHeadersT = it.Expect.AssertHeadersT

		errs = make([]error, 0)
	)

	if len(assertHeaders) == 0 && len(assertHeadersT) == 0 {
		return nil
	}

	t.WithNewStep("Assert headers", func(stepCtx provider.StepCtx) {
		isOption := false
		isOptionT := false

		// Execute assert only body
		for _, f := range assertHeaders {
			executeWithStep(stepCtx, getFunctionName(f), func(t T) []error {
				err := f(headers)
				if err != nil {
					errs = append(errs, err)

					isOption = isOptionError(err)

					return []error{err}
				}

				return nil
			}, true)
		}

		// Execute assert for body with TB
		for _, f := range assertHeadersT {
			executeWithStep(stepCtx, getFunctionName(f), func(t T) []error {
				err := f(t, headers)
				if err != nil {
					errs = append(errs, err)

					isOptionT = isOptionError(err)

					return []error{err}
				}

				return nil
			}, true)
		}

		if isOption && isOptionT {
			stepCtx.CurrentStep().Status = allure.Skipped
		} else {
			stepCtx.CurrentStep().Status = allure.Failed
		}
	})

	return errs
}

func (it *Test) assertResponse(t internalT, response *http.Response) []error {
	var (
		assertResponse  = it.Expect.AssertResponse
		assertResponseT = it.Expect.AssertResponseT

		errs = make([]error, 0)
	)

	if len(assertResponse) == 0 && len(assertResponseT) == 0 {
		return nil
	}

	t.WithNewStep("Assert response", func(stepCtx provider.StepCtx) {
		isOption := false
		isOptionT := false

		// Execute assert only body
		for _, f := range assertResponse {
			executeWithStep(stepCtx, getFunctionName(f), func(t T) []error {
				err := f(response)
				if err != nil {
					errs = append(errs, err)

					isOption = isOptionError(err)

					return []error{err}
				}

				return nil
			}, true)
		}

		// Execute assert for body with TB
		for _, f := range assertResponseT {
			executeWithStep(stepCtx, getFunctionName(f), func(t T) []error {
				err := f(t, response)
				if err != nil {
					errs = append(errs, err)

					isOptionT = isOptionError(err)

					return []error{err}
				}

				return nil
			}, true)
		}

		if isOption && isOptionT {
			stepCtx.CurrentStep().Status = allure.Skipped
		} else {
			stepCtx.CurrentStep().Status = allure.Failed
		}
	})

	return errs
}

func (it *Test) assertBody(t internalT, body []byte) []error {
	var (
		assertBody  = it.Expect.AssertBody
		assertBodyT = it.Expect.AssertBodyT

		errs = make([]error, 0)
	)

	if len(assertBody) == 0 && len(assertBodyT) == 0 {
		return nil
	}

	t.WithNewStep("Assert body", func(stepCtx provider.StepCtx) {
		isOption := false
		isOptionT := false

		// Execute assert only body
		for _, f := range assertBody {
			executeWithStep(stepCtx, getFunctionName(f), func(t T) []error {
				err := f(body)
				if err != nil {
					errs = append(errs, err)

					isOption = isOptionError(err)

					return []error{err}
				}

				return nil
			}, true)
		}

		// Execute assert for body with TB
		for _, f := range assertBodyT {
			executeWithStep(stepCtx, getFunctionName(f), func(t T) []error {
				err := f(t, body)
				if err != nil {
					errs = append(errs, err)

					isOptionT = isOptionError(err)

					return []error{err}
				}

				return nil
			}, true)
		}

		if isOption && isOptionT {
			stepCtx.CurrentStep().Status = allure.Skipped
		} else {
			stepCtx.CurrentStep().Status = allure.Failed
		}
	})

	return errs
}

func isOptionError(err error) bool {
	if tErr, ok := err.(errors.OptionalError); ok {
		return tErr.IsOptional()
	}

	return false
}

func getFunctionName(temp interface{}) string {
	strs := strings.Split(runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name(), ".")
	return strs[len(strs)-2]
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
