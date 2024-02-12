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
		name = "Name"

		testResults ResultsHTTPBuilder = &testResults{
			name:  name,
			state: ResultStateBroken,
			resp:  resp,
			errors: []error{
				firstErr,
				secondErr,
			},
		}
	)

	require.Equal(t, name, testResults.GetName())
	require.Equal(t, false, testResults.IsFailed())
	require.Equal(t, false, testResults.IsFailed())
	require.Equal(t, resp, testResults.GetHTTPResponse())
	require.Equal(t, []error{firstErr, secondErr}, testResults.GetErrors())
}
