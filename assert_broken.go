package cute

import (
	"net/http"

	"github.com/ozontech/cute/errors"
)

func brokenAssertHeaders(assert AssertHeaders) AssertHeaders {
	return func(headers http.Header) error {
		err := assert(headers)

		return wrapBrokenError(err)
	}
}

func brokenAssertBody(assert AssertBody) AssertBody {
	return func(body []byte) error {
		err := assert(body)

		return wrapBrokenError(err)
	}
}

func brokenAssertResponse(assert AssertResponse) AssertResponse {
	return func(resp *http.Response) error {
		err := assert(resp)

		return wrapBrokenError(err)
	}
}

func brokenAssertHeadersT(assert AssertHeadersT) AssertHeadersT {
	return func(t T, headers http.Header) error {
		err := assert(t, headers)

		return wrapBrokenError(err)
	}
}

func brokenAssertBodyT(assert AssertBodyT) AssertBodyT {
	return func(t T, body []byte) error {
		err := assert(t, body)

		return wrapBrokenError(err)
	}
}

func brokenAssertResponseT(assert AssertResponseT) AssertResponseT {
	return func(t T, resp *http.Response) error {
		err := assert(t, resp)

		return wrapBrokenError(err)
	}
}

func wrapBrokenError(err error) error {
	if err == nil {
		return nil
	}

	if tErr, ok := err.(errors.BrokenError); ok {
		tErr.SetBroken(true)

		return tErr.(error)
	}

	return errors.NewBrokenError(err.Error())
}
