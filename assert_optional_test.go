package cute

import (
	"errors"
	"net/http"
	"testing"

	cuteErrors "github.com/ozontech/cute/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptionalAssertResponse(t *testing.T) {
	v := &http.Response{}
	f := func(resp *http.Response) error {
		return errors.New("test error")
	}

	err := optionalAssertResponse(f)(v)

	if optionalError, ok := err.(cuteErrors.OptionalError); assert.True(t, ok) {
		require.True(t, optionalError.IsOptional())
	}
}

func TestOptionalAssertResponseT(t *testing.T) {
	v := &http.Response{}
	f := func(t T, resp *http.Response) error {
		return errors.New("test error")
	}

	err := optionalAssertResponseT(f)(nil, v)

	if optionalError, ok := err.(cuteErrors.OptionalError); assert.True(t, ok) {
		require.True(t, optionalError.IsOptional())
	}
}

func TestOptionalAssertHeaders(t *testing.T) {
	h := http.Header{}
	f := func(headers http.Header) error {
		return errors.New("test error")
	}

	err := optionalAssertHeaders(f)(h)

	if optionalError, ok := err.(cuteErrors.OptionalError); assert.True(t, ok) {
		require.True(t, optionalError.IsOptional())
	}
}

func TestOptionalAssertHeadersT(t *testing.T) {
	h := http.Header{}
	f := func(t T, headers http.Header) error {
		return errors.New("test error")
	}

	err := optionalAssertHeadersT(f)(nil, h)

	if optionalError, ok := err.(cuteErrors.OptionalError); assert.True(t, ok) {
		require.True(t, optionalError.IsOptional())
	}
}

func TestOptionalAssertBody(t *testing.T) {
	v := []byte{}
	f := func(body []byte) error {
		return errors.New("test error")
	}

	err := optionalAssertBody(f)(v)

	if optionalError, ok := err.(cuteErrors.OptionalError); assert.True(t, ok) {
		require.True(t, optionalError.IsOptional())
	}
}

func TestOptionalAssertBodyT(t *testing.T) {
	v := []byte{}
	f := func(t T, body []byte) error {
		return errors.New("test error")
	}

	err := optionalAssertBodyT(f)(nil, v)

	if optionalError, ok := err.(cuteErrors.OptionalError); assert.True(t, ok) {
		require.True(t, optionalError.IsOptional())
	}
}

func TestWrapOptionalError(t *testing.T) {
	err := errors.New("test error")

	optError := wrapOptionalError(err)
	if optionalError, ok := optError.(cuteErrors.OptionalError); assert.True(t, ok) {
		require.True(t, optionalError.IsOptional())
	}
}
