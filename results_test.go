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
		title   = "title"
		feature = "feature"
		epic    = "epic"
		desc    = "desc"
		tags    = []string{"tag_1", "tag_2"}

		isParallel = true

		ht = &test{
			httpClient: &http.Client{},
			allureInfo: &allureInformation{
				title:       title,
				description: desc,
			},
			parallel:    isParallel,
			allureLinks: &allureLinks{},
			allureLabels: &allureLabels{
				feature: feature,
				epic:    epic,
				tags:    tags,
			},
		}

		testResults ResultsHTTPBuilder = &testResults{
			resp: resp,
			errors: []error{
				firstErr,
				secondErr,
			},
			httpTest: ht,
		}
	)

	require.Equal(t, resp, testResults.GetHTTPResponse())
	require.Equal(t, []error{firstErr, secondErr}, testResults.GetErrors())

	hRes, ok := testResults.NextTest().(*test)
	require.True(t, ok)

	require.Equal(t, title, hRes.allureInfo.title)
	require.Equal(t, tags, hRes.allureLabels.tags)
	require.Equal(t, desc, hRes.allureInfo.description)
	require.Equal(t, epic, hRes.allureLabels.epic)
	require.Equal(t, feature, hRes.allureLabels.feature)
	require.Equal(t, isParallel, hRes.parallel)

	require.NotNil(t, hRes.expect)
	require.NotNil(t, hRes.request)
	require.NotNil(t, hRes.middleware)
	require.NotNil(t, hRes.allureStep)
	require.NotNil(t, hRes.httpClient)
	require.NotNil(t, hRes.request)
	require.NotNil(t, hRes.request.repeat)
}
