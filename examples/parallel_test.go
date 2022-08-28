//go:build example
// +build example

package examples

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/cute"
)

func Test_Async_1(t *testing.T) {
	cute.NewTestBuilder().
		Title("Title async test 1").
		Tags("parallel_test").
		Parallel().
		Create().
		BeforeExecuteT(
			func(t cute.T, r *http.Request) error {
				t.WithNewStep("insideBefore", func(stepCtx provider.StepCtx) {
					time.Sleep(time.Second)
					now := time.Now()
					stepCtx.Logf("Test 1. Start time %v", now)
					stepCtx.WithNewParameters("Test 1. Time", now)
				})

				return nil
			},
		).
		AfterExecuteT(
			func(t cute.T, resp *http.Response, errs []error) error {
				t.WithNewStep("insideAfter", func(stepCtx provider.StepCtx) {
					now := time.Now()
					stepCtx.Logf("Test 1. Stop time %v", now)
					stepCtx.WithNewParameters("Test 1. Stop time", now)
				})

				return nil
			}).
		RequestBuilder(
			cute.WithURI("https://jsonplaceholder.typicode.com/posts/1/comments"),
			cute.WithMethod(http.MethodGet),
		).
		ExecuteTest(context.Background(), t)
}

func Test_Async_2(t *testing.T) {
	cute.NewTestBuilder().
		Title("Title async test 2").
		Tags("parallel_test").
		Parallel().
		Create().
		BeforeExecuteT(
			func(t cute.T, r *http.Request) error {
				t.WithNewStep("insideBefore", func(stepCtx provider.StepCtx) {
					now := time.Now()
					stepCtx.Logf("Test 2. Start time %v", now)
					stepCtx.WithNewParameters("Test 2. Start time", now)
					time.Sleep(2 * time.Second)
				})

				return nil
			},
		).
		AfterExecuteT(
			func(t cute.T, resp *http.Response, errs []error) error {
				t.WithNewStep("insideAfter", func(stepCtx provider.StepCtx) {
					now := time.Now()
					stepCtx.Logf("test 2. Stop time %v", now)
					stepCtx.WithNewParameters("Test 2. Stop time", now)
				})

				return nil
			}).
		RequestBuilder(
			cute.WithURI("https://jsonplaceholder.typicode.com/posts/1/comments"),
			cute.WithMethod(http.MethodGet),
		).
		ExpectStatus(200).
		ExecuteTest(context.Background(), t)
}
