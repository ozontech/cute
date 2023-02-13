package cute

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequest(t *testing.T) {
	var (
		req     = newRequestOptions()
		headers = map[string][]string{
			"key": []string{
				"value",
			},
		}
		query = map[string][]string{
			"query_key": []string{
				"query_value",
			},
		}
		method = http.MethodGet
		mBody  = map[string]interface{}{
			"key": map[string]interface{}{
				"some": "value",
			},
			"key_twi": "more",
		}
		uri  = "https://goo.com"
		u, _ = url.Parse("https://ho.com")
		body = []byte("body")
	)

	WithHeaders(headers)(req)
	WithMarshalBody(mBody)(req)
	WithURI(uri)(req)
	WithMethod(method)(req)
	WithURL(u)(req)
	WithBody(body)(req)
	WithQuery(query)(req)

	require.Equal(t, req.headers, headers)
	require.Equal(t, req.uri, uri)
	require.Equal(t, req.bodyMarshal, mBody)
	require.Equal(t, req.method, method)
	require.Equal(t, req.query, query)

	require.Equal(t, req.body, body)
	require.Equal(t, req.url, u)
}
