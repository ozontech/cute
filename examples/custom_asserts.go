package examples

import (
	"errors"
	"net/http"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/stretchr/testify/require"

	"github.com/ozontech/cute"
	cErrors "github.com/ozontech/cute/errors"
)

func CustomAssertBodyWithCustomError() cute.AssertBody {
	return func(bytes []byte) error {
		if len(bytes) == 0 {
			return cErrors.NewAssertError("customAssertBodyWithCustomError", "body must be not empty", "len is 0", "len more 0")
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
		t.WithNewParameters("example_parameter", "example")
		require.GreaterOrEqual(t, len(bytes), 100)
		return nil
	}
}

func CustomAssertBodyWithAllureStep() cute.AssertBodyT {
	return func(t cute.T, bytes []byte) error {

		step := allure.NewSimpleStep("Custom assert step")
		defer func() {
			t.Step(step)
		}()

		if len(bytes) == 0 {
			step.Status = allure.Failed
			step.WithAttachments(allure.NewAttachment("Error", allure.Text, []byte("response body is empty")))

			return nil
		}

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
