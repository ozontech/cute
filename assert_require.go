package cute

import (
	"net/http"

	"github.com/ozontech/cute/errors"
)

func requireAssertHeaders(assert AssertHeaders) AssertHeaders {
	return func(headers http.Header) error {
		err := assert(headers)

		return wrapRequireError(err)
	}
}

func requireAssertBody(assert AssertBody) AssertBody {
	return func(body []byte) error {
		err := assert(body)

		return wrapRequireError(err)
	}
}

func requireAssertResponse(assert AssertResponse) AssertResponse {
	return func(resp *http.Response) error {
		err := assert(resp)

		return wrapRequireError(err)
	}
}

func requireAssertHeadersT(assert AssertHeadersT) AssertHeadersT {
	return func(t T, headers http.Header) error {
		err := assert(t, headers)

		return wrapRequireError(err)
	}
}

func requireAssertBodyT(assert AssertBodyT) AssertBodyT {
	return func(t T, body []byte) error {
		err := assert(t, body)

		return wrapRequireError(err)
	}
}

func requireAssertResponseT(assert AssertResponseT) AssertResponseT {
	return func(t T, resp *http.Response) error {
		err := assert(t, resp)

		return wrapRequireError(err)
	}
}

func wrapRequireError(err error) error {
	if err == nil {
		return nil
	}

	if tErr, ok := err.(errors.RequireError); ok {
		tErr.SetRequire(true)

		return tErr.(error)
	}

	return errors.NewRequireError(err.Error())
}
