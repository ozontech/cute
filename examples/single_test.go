//go:build example
// +build example

package examples

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"testing"
	"time"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"

	"github.com/ozontech/cute"
	"github.com/ozontech/cute/asserts/json"
)

func Test_Single_1(t *testing.T) {
	cute.NewTestBuilder().
		Title("Single test with default T").
		Tag("single_test").
		Description("some_description").
		Parallel().
		Create().
		RequestBuilder(
			cute.WithURI("https://jsonplaceholder.typicode.com/posts/1/comments"),
			cute.WithMethod(http.MethodGet),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusOK).
		AssertBody(
			json.Present("$[1].name"),
			json.Present("$[0].passport"), // Example fail
			json.Equal("$[0].email", "Eliseo@gardner.biz"),
			CustomAssertBody(),
		).
		AssertBodyT(func(t cute.T, body []byte) error {
			t.Step(allure.NewSimpleStep("inside Assert body. 1 ", allure.NewParameters("key", "value")...))

			return nil
		}).
		ExecuteTest(context.Background(), t)
}

func Test_Single_2_AllureRunner(t *testing.T) {
	runner.Run(t, "Single test with allure-go Runner", func(t provider.T) {
		var (
			testMaker   = cute.NewHTTPTestMaker()
			testBuilder = testMaker.NewTestBuilder()
		)

		u, _ := url.Parse("https://jsonplaceholder.typicode.com/")
		u.Path = path.Join(u.Path, "/posts/1/comments")

		testBuilder.
			Title("Single test with allure.T and repeat errors").
			Tag("single_test").
			Description("some_description").
			Create().
			RequestRepeatDelay(3*time.Second). // delay before new try
			RequestRepeat(3).                  // count attempts
			RequestBuilder(
				cute.WithURL(u),
				cute.WithMethod(http.MethodGet),
			).
			ExpectExecuteTimeout(10*time.Second).
			ExpectStatus(http.StatusBadGateway).
			AssertBody(
				json.Equal("$[0].email", "Eliseo@gardner.biz"),
				json.Present("$[1].name"),
			).
			OptionalAssertBody(
				json.Present("$[0].photo"), // Example optional fail
			).
			ExecuteTest(context.Background(), t)
	})
}
