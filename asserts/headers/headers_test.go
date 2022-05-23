package headers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPresent(t *testing.T) {
	headers := http.Header{
		"Content-Type": []string{"application/json"},
	}

	err := Present("Content-Type")(headers)
	require.NoError(t, err)
}

func TestPresentError(t *testing.T) {
	headers := http.Header{
		"Content-Type": []string{},
	}

	err := Present("not-present")(headers)
	require.Error(t, err)
}

func TestNotPresent(t *testing.T) {
	headers := http.Header{
		"Content-Type": []string{"", "application/json"},
	}

	err := NotPresent("Content-Type")(headers)
	require.Error(t, err)
}

func TestNotPresentError(t *testing.T) {
	headers := http.Header{
		"Content-Type": []string{},
	}

	err := NotPresent("not-present")(headers)
	require.NoError(t, err)
}
