package cute

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ozontech/testo"
	allure "github.com/ozontech/testo-allure"

	"github.com/ozontech/cute/internal/utils"
)

// TestMain redirects allure output of the whole package to a temp dir so
// unit tests leave no artifacts (the flag is also picked up by the child
// tests that Test.Execute spawns, which do not inherit plugin options).
func TestMain(m *testing.M) {
	flag.Parse()

	dir, err := os.MkdirTemp("", "cute-allure-results")
	if err == nil {
		_ = flag.Set("allure.dir", dir)
	}

	code := m.Run()

	if dir != "" {
		_ = os.RemoveAll(dir)
	}

	os.Exit(code)
}

// runCuteTest runs f with a live cute.T (testo.T + allure plugin).
func runCuteTest(t *testing.T, f func(ct defaultT)) {
	testo.RunTest(t, f, allure.WithOutputDir(t.TempDir()))
}

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

	runCuteTest(t, func(ct defaultT) {
		errs := ht.validateResponse(ct, &http.Response{})
		require.Empty(ct, errs)
	})
}

func TestValidateResponseCode(t *testing.T) {
	ht := &Test{
		Expect: &Expect{Code: 200},
	}
	runCuteTest(t, func(ct defaultT) {
		errs := ht.validateResponse(ct, &http.Response{StatusCode: http.StatusOK})
		require.Empty(ct, errs)
	})
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

	runCuteTest(t, func(ct defaultT) {
		errs := ht.validateResponse(ct, resp)

		require.Len(ct, errs, 2)
	})
}

type mockRoundTripper struct{}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Request:    req,
		Body:       io.NopCloser(strings.NewReader("mock response")),
	}, nil
}

func TestSanitizeURLHook(t *testing.T) {
	client := &http.Client{
		Transport: &mockRoundTripper{},
	}

	test := &Test{
		httpClient: client,
		Retry: &Retry{
			currentCount: 0,
			MaxAttempts:  0,
			Delay:        0,
		},
		Request: &Request{
			Builders: []RequestBuilder{
				WithMethod(http.MethodGet),
				WithURI("http://localhost/api?key=123"),
			},
			Retry: &RequestRetryPolitic{
				Count: 1,
				Delay: 2,
			},
		},
		RequestSanitizer: sanitizeKeyParam("****"),
	}

	req, err := test.createRequest(context.Background())
	require.NoError(t, err)
	require.NotNil(t, req)

	runCuteTest(t, func(ct defaultT) {
		err = test.addInformationRequest(ct, req)
		require.NoError(ct, err)

		decodedQuery, decodeErr := url.QueryUnescape(req.URL.RawQuery)
		require.NoError(ct, decodeErr)
		require.Equal(ct, "key=****", decodedQuery)
	})
}

func TestSanitizeURL_LastRequestURL(t *testing.T) {
	client := &http.Client{
		Transport: &mockRoundTripper{},
	}

	test := &Test{
		httpClient: client,
		Request: &Request{
			Builders: []RequestBuilder{
				WithMethod(http.MethodGet),
				WithURI("http://localhost/api?key=123"),
			},
		},
		RequestSanitizer: sanitizeKeyParam("****"),
	}

	runCuteTest(t, func(ct defaultT) {
		test.Execute(context.Background(), ct)
	})

	decodedURL, err := url.QueryUnescape(test.lastRequestURL)
	require.NoError(t, err)
	require.Contains(t, decodedURL, "key=****", "Expected masked key in lastRequestURL")
}

func sanitizeKeyParam(mask string) RequestSanitizerHook {
	return func(req *http.Request) {
		q := req.URL.Query()
		q.Set("key", mask)
		req.URL.RawQuery = q.Encode()
	}
}

func TestSanitizeURL_RealRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		t.Logf("Server received URL: %s, Body: %s", r.URL.String(), string(body))
		require.Contains(t, r.URL.String(), "key=123", "Sanitizer must not change real request")
		w.WriteHeader(200)
	}))
	defer ts.Close()

	client := &http.Client{}
	test := &Test{
		httpClient: client,
		Request: &Request{
			Builders: []RequestBuilder{
				WithMethod(http.MethodGet),
				WithURI(ts.URL + "/api?key=123"),
			},
		},
		RequestSanitizer: sanitizeKeyParam("****"),
	}

	runCuteTest(t, func(ct defaultT) {
		test.Execute(context.Background(), ct)
	})
}
