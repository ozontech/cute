package cute

import (
	"net/http"

	"github.com/ozontech/cute/errors"
)

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

	return errors.WrapOptionalError(err)
}
