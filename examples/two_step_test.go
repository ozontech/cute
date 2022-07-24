package examples

import (
	"context"
	"io"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	"github.com/ozontech/cute"
)

func TestExample_TwoSteps(t *testing.T) {
	runner.Run(t, "Test with two steps", func(t provider.T) {
		test := cute.NewTestBuilder().
			Title("Two steps").
			Description("some_description").
			CreateWithStep().
			StepName("Request 1").
			RequestBuilder(
				cute.WithURI("https://jsonplaceholder.typicode.com/posts/1/comments"),
				cute.WithMethod(http.MethodGet),
			).
			ExpectStatus(http.StatusOK).
			ExecuteTest(context.Background(), t)

		bodyBytes, err := io.ReadAll(test.GetHTTPResponse().Body)
		if err != nil {
			log.Fatal(err)
		}
		// process body
		_ = string(bodyBytes)

		cute.NewTestBuilder().
			CreateWithStep().
			StepName("Request 2").
			RequestBuilder(
				cute.WithURI("https://jsonplaceholder.typicode.com/posts/1/comments"),
				cute.WithMethod(http.MethodGet),
			).
			ExpectExecuteTimeout(10*time.Second).
			ExpectStatus(http.StatusOK).
			ExecuteTest(context.Background(), t)
	})
}
