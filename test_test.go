package cute

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
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

	ht := cute{
		countTests:  1,
		correctTest: 0,
		tests: []*test{
			{
				request: &request{
					base: req,
				},
			},
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

	req, err := http.NewRequest(http.MethodGet, url, ioutil.NopCloser(bytes.NewReader(body)))
	require.NoError(t, err)

	req.Header = headers

	ht := cute{
		countTests:  1,
		correctTest: 0,
		tests: []*test{
			{
				request: &request{
					builders: []requestBuilder{
						WithURI(url),
						WithMethod("GET"),
						WithHeaders(headers),
						WithBody(body),
					},
				},
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

	ht := cute{
		countTests:  1,
		correctTest: 0,
		tests: []*test{
			{
				request: &request{
					builders: []requestBuilder{
						WithMarshalBody(str),
					},
				},
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
	ht := cute{}

	require.Error(t, ht.validateRequest(&http.Request{}))
}

func TestValidateRequestEmptyMethod(t *testing.T) {
	ht := cute{}
	u, _ := url.Parse("https://go.com")

	require.Error(t, ht.validateRequest(&http.Request{
		URL: u,
	}))
}

func TestValidateResponseEmpty(t *testing.T) {
	ht := cute{
		tests: []*test{
			{
				expect: new(Expect),
			},
		},
	}
	temp := common.NewT(t)

	errs := ht.validateResponse(temp, &http.Response{})
	require.Empty(t, errs)
}

func TestValidateResponseCode(t *testing.T) {
	ht := cute{
		tests: []*test{
			{
				expect: &Expect{code: 200},
			},
		},
	}
	temp := common.NewT(t)

	errs := ht.validateResponse(temp, &http.Response{StatusCode: http.StatusOK})
	require.Empty(t, errs)
}

func TestValidateResponseWithErrors(t *testing.T) {
	var (
		ht = cute{
			tests: []*test{
				{
					expect: &Expect{
						code: 200,
						assertHeaders: []AssertHeaders{
							func(headers http.Header) error {
								return errors.New("two error")
							},
						},
						assertResponse: []AssertResponse{
							func(response *http.Response) error {
								if response.StatusCode != http.StatusOK || len(response.Header["auth"]) == 0 {
									return errors.New("bad response")
								}
								return nil
							},
						},
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
			Body: ioutil.NopCloser(reader),
		}
	)

	errs := ht.validateResponse(temp, resp)

	require.Len(t, errs, 2)
}
