package cute

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/core/common"
	"github.com/stretchr/testify/require"
)

func TestValidateRequestEmptyUrl(t *testing.T) {
	ht := test{}

	require.Error(t, ht.validateRequest(&http.Request{}))
}

func TestValidateRequestEmptyMethod(t *testing.T) {
	ht := test{}
	u, _ := url.Parse("https://go.com")

	require.Error(t, ht.validateRequest(&http.Request{
		URL: u,
	}))
}

func TestValidateResponseEmpty(t *testing.T) {
	ht := test{
		expect: new(expect),
	}
	temp := common.NewT(t, "package", t.Name())

	errs := ht.validateResponse(temp, &http.Response{})
	require.Empty(t, errs)
}

func TestValidateResponseCode(t *testing.T) {
	ht := test{
		expect: &expect{code: 200},
	}
	temp := common.NewT(t, "package", t.Name())

	errs := ht.validateResponse(temp, &http.Response{StatusCode: http.StatusOK})
	require.Empty(t, errs)
}

func TestValidateResponseWithErrors(t *testing.T) {
	ht := test{
		expect: &expect{
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
	}
	reader := bytes.NewReader([]byte(`{"a":"ab","b":"bc"}`))
	temp := common.NewT(t, "package", t.Name())
	temp.NewTest(t.Name(), "package")
	temp.TestContext()

	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Header: map[string][]string{
			"key":  []string{"value"},
			"auth": []string{"sometoken"},
		},
		Body: ioutil.NopCloser(reader),
	}

	errs := ht.validateResponse(temp, resp)

	require.Len(t, errs, 3)
}
