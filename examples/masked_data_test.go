//go:build example
// +build example

package examples

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/ozontech/testo"
	allure "github.com/ozontech/testo-allure"

	"github.com/ozontech/cute"
)

func TestSanitizer(t *testing.T) {
	t.Run("Single test with request and response sanitizer", testo.Test(func(t T) {

		allure.Step(t, "First step", func(t T) {
			allure.Step(t, "Inside first step", func(T) {})
		})

		allure.Step(t, "Step name", func(t T) {
			u, _ := url.Parse("https://jsonplaceholder.typicode.com/posts/1/comments?example=11")
			query := u.Query()
			query.Set("name", "Vasya")
			u.RawQuery = query.Encode()

			cute.NewTestBuilder().
				Title("Super simple test").
				Tags("simple", "suite", "some_local_tag", "json").
				Parallel().
				Create().
				RequestSanitizerHook(func(req *http.Request) {
					req.URL.Path = "/path/masked"

					values := req.URL.Query()
					values.Set("example", "masked")

					req.URL.RawQuery = values.Encode()

					req.Header["some_header"] = []string{"masked"}
				}).
				ResponseSanitizerHook(func(resp *http.Response) {
					resp.Header["some_header"] = []string{"masked"}
					resp.Header["Content-Type"] = []string{"masked"}
				}).
				RequestBuilder(
					cute.WithHeaders(map[string][]string{
						"some_header": []string{"something"},
					}),
					cute.WithURL(u),
					cute.WithMethod(http.MethodPost),
				).
				ExpectExecuteTimeout(10*time.Second).
				ExpectStatus(http.StatusCreated).
				ExecuteTest(context.Background(), t)
		})
	}))

}
