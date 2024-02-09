//go:build example
// +build example

package examples

import (
	"context"
	"errors"
	"fmt"
	"io"
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
		RequestRepeat(3).
		RequestBuilder(
			cute.WithURI("https://jsonplaceholder.typicode.com/posts/1/comments"),
			cute.WithMarshalBody(struct {
				Name string `json:"name"`
			}{
				Name: "Vasya Pupkin",
			}),
			cute.WithQueryKV("socks", "42"),
			cute.WithMethod(http.MethodGet),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusOK).
		AssertBody(json.Diff("{\"aaa\":\"bb\"}")).
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
		After(
			func(response *http.Response, errors []error) error {
				b, err := io.ReadAll(response.Body)
				if err != nil {
					return err
				}

				email, err := json.GetValueFromJSON(b, "$[0].email")
				if err != nil {
					return err
				}

				fmt.Println("Email from test", email)

				return nil
			},
		).
		ExecuteTest(context.Background(), t)
}

func Test_Single_Broken111(t *testing.T) {
	cute.NewTestBuilder().
		Title("Test_Single_Broken").
		Create().
		RequestBuilder(
			cute.WithURI("https://jsonplaceholder.typicode.com/posts/1/comments"),
		).
		BrokenAssertBodyT(func(t cute.T, body []byte) error {
			return errors.New("example broken error")
		}).
		ExpectStatus(http.StatusOK).
		NextTest().
		Create().
		RequestBuilder(
			cute.WithURI("https://jsonplaceholder.typicode.com/posts/1/comments"),
		).
		AssertBody(func(body []byte) error {
			return errors.New("ssss")
		},
		).
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
