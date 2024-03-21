package cute

import (
	"testing"

	cute_errors "github.com/ozontech/cute/errors"
)

func TestCallAssertWrapper(t *testing.T) {
	f := callAssertWrapper(func(body []byte) error {
		return cute_errors.NewAssertErrorWithMessage("kek", "lol")
	}, "it should be in new line and link")

	err := f(nil)

	t.Logf(err.Error())

	cuteError := err.(cute_errors.CuteError)

	//cerr := errors.Unwrap(err)

	t.Logf("%+v", cuteError.GetName())
	//cerr, ok := err.(cute_errors.WithNameError)
	//require.True(t, ok)
	//require.Equal(t, "kek", cerr.GetName())

	t.Logf(err.Error())
}
