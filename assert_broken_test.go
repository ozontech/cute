package cute

import (
	"errors"
	"net/http"
	"testing"

	cuteErrors "github.com/ozontech/cute/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBrokenAssertResponse(t *testing.T) {
	v := &http.Response{}
	f := func(_ *http.Response) error {
		return errors.New("test error")
	}

	err := brokenAssertResponse(f)(v)

	if BrokenError, ok := err.(cuteErrors.BrokenError); assert.True(t, ok) {
		require.True(t, BrokenError.IsBroken())
	}
}

func TestBrokenAssertResponseT(t *testing.T) {
	v := &http.Response{}
	f := func(T, *http.Response) error {
		return errors.New("test error")
	}

	err := brokenAssertResponseT(f)(nil, v)

	if BrokenError, ok := err.(cuteErrors.BrokenError); assert.True(t, ok) {
		require.True(t, BrokenError.IsBroken())
	}
}

func TestBrokenAssertHeaders(t *testing.T) {
	h := http.Header{}
	f := func(_ http.Header) error {
		return errors.New("test error")
	}

	err := brokenAssertHeaders(f)(h)

	if BrokenError, ok := err.(cuteErrors.BrokenError); assert.True(t, ok) {
		require.True(t, BrokenError.IsBroken())
	}
}

func TestBrokenAssertHeadersT(t *testing.T) {
	h := http.Header{}
	f := func(T, http.Header) error {
		return errors.New("test error")
	}

	err := brokenAssertHeadersT(f)(nil, h)

	if BrokenError, ok := err.(cuteErrors.BrokenError); assert.True(t, ok) {
		require.True(t, BrokenError.IsBroken())
	}
}

func TestBrokenAssertBody(t *testing.T) {
	v := []byte{}
	f := func(_ []byte) error {
		return errors.New("test error")
	}

	err := brokenAssertBody(f)(v)

	if BrokenError, ok := err.(cuteErrors.BrokenError); assert.True(t, ok) {
		require.True(t, BrokenError.IsBroken())
	}
}

func TestBrokenAssertBodyT(t *testing.T) {
	v := []byte{}
	f := func(T, []byte) error {
		return errors.New("test error")
	}

	err := brokenAssertBodyT(f)(nil, v)

	if BrokenError, ok := err.(cuteErrors.BrokenError); assert.True(t, ok) {
		require.True(t, BrokenError.IsBroken())
	}
}

func TestWrapBrokenError(t *testing.T) {
	err := errors.New("test error")

	optError := wrapBrokenError(err)
	if BrokenError, ok := optError.(cuteErrors.BrokenError); assert.True(t, ok) {
		require.True(t, BrokenError.IsBroken())
	}
}
