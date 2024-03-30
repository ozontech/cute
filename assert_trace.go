package cute

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/ozontech/cute/errors"
)

// assertHeadersWithTrace is a function to add trace inside assert headers error
func assertHeadersWithTrace(assert AssertHeaders, trace string) AssertHeaders {
	return func(headers http.Header) error {
		err := assert(headers)

		return wrapWithTrace(err, trace)
	}
}

// assertBodyWithTrace is a function to add trace inside assert body error
func assertBodyWithTrace(assert AssertBody, trace string) AssertBody {
	return func(body []byte) error {
		err := assert(body)

		return wrapWithTrace(err, trace)
	}
}

// assertResponseWithTrace is a function to add trace inside assert response error
func assertResponseWithTrace(assert AssertResponse, trace string) AssertResponse {
	return func(resp *http.Response) error {
		err := assert(resp)

		return wrapWithTrace(err, trace)
	}
}

// assertHeadersTWithTrace is a function to add trace inside assert headers error
func assertHeadersTWithTrace(assert AssertHeadersT, trace string) AssertHeadersT {
	return func(t T, headers http.Header) error {
		err := assert(t, headers)

		return wrapWithTrace(err, trace)
	}
}

// assertBodyTWithTrace is a function to add trace inside assert body error
func assertBodyTWithTrace(assert AssertBodyT, trace string) AssertBodyT {
	return func(t T, body []byte) error {
		err := assert(t, body)

		return wrapWithTrace(err, trace)
	}
}

// assertResponseTWithTrace is a function to add trace inside assert response error
func assertResponseTWithTrace(assert AssertResponseT, trace string) AssertResponseT {
	return func(t T, resp *http.Response) error {
		err := assert(t, resp)

		return wrapWithTrace(err, trace)
	}
}

// wrapWithTrace is a function to add trace inside error
func wrapWithTrace(err error, trace string) error {
	if err == nil {
		return nil
	}

	if tErr, ok := err.(errors.WithTrace); ok {
		tErr.SetTrace(trace)

		return tErr.(error)
	}

	return errors.WrapErrorWithTrace(err, trace)
}

func getTrace() string {
	pcs := make([]uintptr, 10)
	depth := runtime.Callers(3, pcs)

	if depth == 0 {
		fmt.Println("Couldn't get the stack information")
		return ""
	}

	callers := runtime.CallersFrames(pcs[:depth])
	caller, _ := callers.Next()

	return fmt.Sprintf("%s:%d", caller.File, caller.Line)
}
