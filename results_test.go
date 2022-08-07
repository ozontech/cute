package cute

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResult(t *testing.T) {
	var (
		firstErr  = errors.New("first_error")
		secondErr = errors.New("second_error")
		resp      = &http.Response{
			Status:     "OK",
			StatusCode: 200,
		}
		name = "name"

		testResults ResultsHTTPBuilder = &testResults{
			name: name,
			resp: resp,
			errors: []error{
				firstErr,
				secondErr,
			},
		}
	)

	require.Equal(t, resp, testResults.GetHTTPResponse())
	require.Equal(t, []error{firstErr, secondErr}, testResults.GetErrors())
	require.Equal(t, name, testResults.GetName())
}
