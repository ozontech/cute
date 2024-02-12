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
	f := func(*http.Response) error {
		return errors.New("test error")
	}

	err := optionalAssertResponse(f)(v)

	if optionalError, ok := err.(cuteErrors.OptionalError); assert.True(t, ok) {
		require.True(t, optionalError.IsOptional())
	}
}

func TestOptionalAssertResponseT(t *testing.T) {
	v := &http.Response{}
	f := func(T, *http.Response) error {
		return errors.New("test error")
	}

	err := optionalAssertResponseT(f)(nil, v)

	if optionalError, ok := err.(cuteErrors.OptionalError); assert.True(t, ok) {
		require.True(t, optionalError.IsOptional())
	}
}

func TestOptionalAssertHeaders(t *testing.T) {
	h := http.Header{}
	f := func(http.Header) error {
		return errors.New("test error")
	}

	err := optionalAssertHeaders(f)(h)

	if optionalError, ok := err.(cuteErrors.OptionalError); assert.True(t, ok) {
		require.True(t, optionalError.IsOptional())
	}
}

func TestOptionalAssertHeadersT(t *testing.T) {
	h := http.Header{}
	f := func(T, http.Header) error {
		return errors.New("test error")
	}

	err := optionalAssertHeadersT(f)(nil, h)

	if optionalError, ok := err.(cuteErrors.OptionalError); assert.True(t, ok) {
		require.True(t, optionalError.IsOptional())
	}
}

func TestOptionalAssertBody(t *testing.T) {
	v := []byte{}
	f := func([]byte) error {
		return errors.New("test error")
	}

	err := optionalAssertBody(f)(v)

	if optionalError, ok := err.(cuteErrors.OptionalError); assert.True(t, ok) {
		require.True(t, optionalError.IsOptional())
	}
}

func TestOptionalAssertBodyT(t *testing.T) {
	v := []byte{}
	f := func(T, []byte) error {
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
