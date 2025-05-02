package cute

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/core/common"
	"github.com/ozontech/cute/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateRequest(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://go.com", nil)
	require.NoError(t, err)

	ht := &Test{
		Request: &Request{
			Base: req,
		},
	}

	resReq, err := ht.createRequest(context.Background())
	require.NoError(t, err)
	require.Equal(t, req, resReq)
}

func TestCreateRequestBuilder(t *testing.T) {
	var (
		body    = []byte("HELLO")
		headers = map[string][]string{
			"Good": []string{"Day"},
			"Mad":  []string{"Max"},
		}
		url = "http://go.com"
	)

	req, err := http.NewRequest(http.MethodGet, url, io.NopCloser(bytes.NewReader(body)))
	require.NoError(t, err)

	req.Header = headers

	ht := &Test{
		Request: &Request{
			Builders: []RequestBuilder{
				WithURI(url),
				WithMethod("GET"),
				WithHeaders(headers),
				WithBody(body),
			},
		},
	}

	resReq, err := ht.createRequest(context.Background())
	require.NoError(t, err)
	require.Equal(t, req, resReq)
}

func TestCreateRequestBuilder_MarshalBody(t *testing.T) {
	var (
		str = struct {
			name string
		}{
			"hello",
		}
	)

	ht := &Test{
		jsonMarshaler: jsonMarshaler{},
		Request: &Request{
			Builders: []RequestBuilder{
				WithMarshalBody(str),
			},
		},
	}

	resReq, err := ht.createRequest(context.Background())
	require.NoError(t, err)

	getBody, err := utils.GetBody(resReq.Body)
	require.NoError(t, err)

	require.NotEmpty(t, getBody)
}

func TestValidateRequestEmptyUrl(t *testing.T) {
	ht := &Test{}

	require.Error(t, ht.validateRequest(&http.Request{}))
}

func TestValidateRequestEmptyMethod(t *testing.T) {
	ht := &Test{}
	u, _ := url.Parse("https://go.com")

	require.Error(t, ht.validateRequest(&http.Request{
		URL: u,
	}))
}

func TestValidateResponseEmpty(t *testing.T) {
	ht := &Test{
		Expect: new(Expect),
	}

	temp := common.NewT(t)

	errs := ht.validateResponse(temp, &http.Response{})
	require.Empty(t, errs)
}

func TestValidateResponseCode(t *testing.T) {
	ht := &Test{
		Expect: &Expect{Code: 200},
	}
	temp := common.NewT(t)

	errs := ht.validateResponse(temp, &http.Response{StatusCode: http.StatusOK})
	require.Empty(t, errs)
}

func TestValidateResponseWithErrors(t *testing.T) {
	var (
		ht = &Test{
			Expect: &Expect{
				Code:       200,
				JSONSchema: new(ExpectJSONSchema),
				AssertHeaders: []AssertHeaders{
					func(headers http.Header) error {
						return errors.New("two error")
					},
				},
				AssertResponse: []AssertResponse{
					func(response *http.Response) error {
						if response.StatusCode != http.StatusOK || len(response.Header["auth"]) == 0 { //nolint
							return errors.New("bad response")
						}
						return nil
					},
				},
			},
		}

		reader = bytes.NewReader([]byte(`{"a":"ab","b":"bc"}`))
		temp   = createAllureT(t)
		resp   = &http.Response{
			StatusCode: http.StatusBadRequest,
			Header: map[string][]string{
				"key":  []string{"value"},
				"auth": []string{"sometoken"},
			},
			Body: io.NopCloser(reader),
		}
	)

	ht.initEmptyFields()

	errs := ht.validateResponse(temp, resp)

	require.Len(t, errs, 2)
}

func TestSanitizeURLHook(t *testing.T) {
	var sanitized string

	test := &Test{
		Name: "Test with sanitization",
		Request: &Request{
			Builders: []RequestBuilder{
				WithMethod(http.MethodGet),
				WithURI("http://localhost/api?key=123"),
			},
		},
		SanitizeURL: func(req *http.Request) {
			q := req.URL.Query()
			q.Set("key", "****")
			req.URL.RawQuery = q.Encode()
			decoded, err := url.QueryUnescape(req.URL.RawQuery)
			require.NoError(t, err)
			sanitized = decoded
		},
	}

	_, err := test.createRequest(context.Background())
	require.NoError(t, err)
	require.Equal(t, "key=****", sanitized)
}
