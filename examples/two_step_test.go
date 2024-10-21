//go:build example
// +build example

package examples

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	"github.com/ozontech/cute"
)

func Test_TwoSteps_1(t *testing.T) {
	cute.NewTestBuilder().
		Title("Test with two requests.").
		Tags("two_steps").
		Parallel().
		CreateStep("Create entry /posts/1").

		// CreateWithStep first step

		RequestBuilder(
			cute.WithURI("https://jsonplaceholder.typicode.com/posts/1/comments"),
			cute.WithMethod(http.MethodGet),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusCreated).
		NextTest().
		CreateStep("Delete entry").

		// CreateWithStep second step for delete
		RequestBuilder(
			cute.WithURI("https://jsonplaceholder.typicode.com/posts/1/comments"),
			cute.WithMethod(http.MethodDelete),
			cute.WithHeaders(map[string][]string{
				"some_auth_token": []string{fmt.Sprint(11111)},
			}),
		).
		ExecuteTest(context.Background(), t)
}

func Test_TwoSteps_2_AllureRunner(t *testing.T) {
	runner.Run(t, "Test with two steps", func(t provider.T) {
		testBuilder := cute.NewHTTPTestMaker().NewTestBuilder()

		testBuilder.
			Title("Test with two requests executed by allure-go").
			Tag("two_steps").
			Description("some_description").
			CreateStep("Request 1").
			RequestBuilder(
				cute.WithURI("https://jsonplaceholder.typicode.com/posts/1/comments"),
				cute.WithMethod(http.MethodGet),
			).
			ExpectStatus(http.StatusOK).
			ExecuteTest(context.Background(), t)

		testBuilder.
			CreateStep("Request 2").
			RequestBuilder(
				cute.WithURI("https://jsonplaceholder.typicode.com/posts/2/comments"),
				cute.WithMethod(http.MethodGet),
			).
			ExpectExecuteTimeout(10*time.Second).
			ExpectStatus(http.StatusOK).
			ExecuteTest(context.Background(), t)
	})
}

func Test_TwoSteps_3(t *testing.T) {
	responseCode := 0

	// First step.
	cute.NewTestBuilder().
		Title("Test with two requests and parse body.").
		Tag("two_steps").
		Create().
		RequestBuilder(
			cute.WithURI("https://jsonplaceholder.typicode.com/posts/1/comments"),
			cute.WithMethod(http.MethodGet),
		).
		ExpectStatus(http.StatusOK).
		RequireBody(func(body []byte) error {
			return errors.New("example")
		}).
		NextTest().
		AfterTestExecute(func(response *http.Response, errors []error) error { // Execute after first step
			responseCode = response.StatusCode

			fmt.Println("Hello from after test execute")
			fmt.Println("Response code", responseCode)

			return nil
		}).
		// Second step. This test isn't run, because previous test has failed require validation
		Create().
		RequestBuilder(
			cute.WithURI("https://jsonplaceholder.typicode.com/posts/2/comments"),
			cute.WithMethod(http.MethodDelete),
		).
		ExecuteTest(context.Background(), t)

	fmt.Println("Response code from first request", responseCode)
}
