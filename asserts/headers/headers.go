package headers

import (
	"fmt"
	"net/http"

	"github.com/ozontech/cute"
	"github.com/ozontech/cute/errors"
)

// Present is a function to asserts that header is present
func Present(key string) cute.AssertHeaders {
	return func(headers http.Header) error {
		if v := headers.Get(key); v == "" {
			return errors.NewAssertError("Present", fmt.Sprintf("header %s is not present", key), nil, nil)
		}

		return nil
	}
}

// NotPresent is a function to asserts that header is not present
func NotPresent(key string) cute.AssertHeaders {
	return func(headers http.Header) error {
		if v := headers.Values(key); len(v) > 0 {
			return errors.NewAssertError("NotPresent", fmt.Sprintf("header %s is present", key), nil, nil)
		}

		return nil
	}
}
