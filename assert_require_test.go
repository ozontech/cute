package cute

import (
	"errors"
	"net/http"
	"testing"

	cuteErrors "github.com/ozontech/cute/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequireAssertResponse(t *testing.T) {
	v := &http.Response{}
	f := func(_ *http.Response) error {
		return errors.New("test error")
	}

	err := requireAssertResponse(f)(v)

	if RequireError, ok := err.(cuteErrors.RequireError); assert.True(t, ok) {
		require.True(t, RequireError.IsRequire())
	}
}

func TestRequireAssertResponseT(t *testing.T) {
	v := &http.Response{}
	f := func(T, *http.Response) error {
		return errors.New("test error")
	}

	err := requireAssertResponseT(f)(nil, v)

	if RequireError, ok := err.(cuteErrors.RequireError); assert.True(t, ok) {
		require.True(t, RequireError.IsRequire())
	}
}

func TestRequireAssertHeaders(t *testing.T) {
	h := http.Header{}
	f := func(http.Header) error {
		return errors.New("test error")
	}

	err := requireAssertHeaders(f)(h)

	if RequireError, ok := err.(cuteErrors.RequireError); assert.True(t, ok) {
		require.True(t, RequireError.IsRequire())
	}
}

func TestRequireAssertHeadersT(t *testing.T) {
	h := http.Header{}
	f := func(T, http.Header) error {
		return errors.New("test error")
	}

	err := requireAssertHeadersT(f)(nil, h)

	if RequireError, ok := err.(cuteErrors.RequireError); assert.True(t, ok) {
		require.True(t, RequireError.IsRequire())
	}
}

func TestRequireAssertBody(t *testing.T) {
	v := []byte{}
	f := func([]byte) error {
		return errors.New("test error")
	}

	err := requireAssertBody(f)(v)

	if RequireError, ok := err.(cuteErrors.RequireError); assert.True(t, ok) {
		require.True(t, RequireError.IsRequire())
	}
}

func TestRequireAssertBodyT(t *testing.T) {
	v := []byte{}
	f := func(T, []byte) error {
		return errors.New("test error")
	}

	err := requireAssertBodyT(f)(nil, v)

	if RequireError, ok := err.(cuteErrors.RequireError); assert.True(t, ok) {
		require.True(t, RequireError.IsRequire())
	}
}

func TestWrapRequireError(t *testing.T) {
	err := errors.New("test error")

	optError := wrapRequireError(err)
	if RequireError, ok := optError.(cuteErrors.RequireError); assert.True(t, ok) {
		require.True(t, RequireError.IsRequire())
	}
}
