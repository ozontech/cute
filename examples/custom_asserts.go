package examples

import (
	"errors"
	"net/http"

	"github.com/ozontech/cute"
	cuteErrors "github.com/ozontech/cute/errors"
	allure "github.com/ozontech/testo-allure"
	"github.com/stretchr/testify/require"
)

func CustomAssertBodyWithCustomError() cute.AssertBody {
	return func(bytes []byte) error {
		if len(bytes) == 0 {
			return cuteErrors.NewAssertError("customAssertBodyWithCustomError", "body must be not empty", "len is 0", "len more 0")
		}

		return nil
	}
}

func CustomAssertBody() cute.AssertBody {
	return func(bytes []byte) error {
		if len(bytes) == 0 {
			return errors.New("response body is empty")
		}

		return nil
	}
}

func CustomAssertBodyT() cute.AssertBodyT {
	return func(t cute.T, bytes []byte) error {
		t.Parameters(allure.NewParameter("example_parameter", "example"))
		require.GreaterOrEqual(t, len(bytes), 100)
		return nil
	}
}

func CustomAssertBodyWithAllureStep() cute.AssertBodyT {
	return func(t cute.T, bytes []byte) error {
		allure.Step(t, "Custom assert step", func(t cute.T) {
			if len(bytes) == 0 {
				t.Attach("Error", allure.Bytes("response body is empty").As(allure.TextPlain))
				t.Status(allure.StatusFailed)
			}
		})

		return nil
	}
}

func CustomAssertHeaders() cute.AssertHeaders {
	return func(headers http.Header) error {
		if len(headers) == 0 {
			return errors.New("response without headers")
		}

		return nil
	}
}

func CustomAssertResponse() cute.AssertResponse {
	return func(resp *http.Response) error {
		if resp.ContentLength == 0 {
			return errors.New("content length is zero")
		}

		return nil
	}
}
