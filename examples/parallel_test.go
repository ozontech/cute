//go:build example
// +build example

package examples

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/ozontech/cute"
	allure "github.com/ozontech/testo-allure"
)

func Test_Async_1(t *testing.T) {
	cute.NewTestBuilder().
		Title("Title async test 1").
		Tags("parallel_test").
		Parallel().
		Create().
		BeforeExecuteT(
			func(t cute.T, r *http.Request) error {
				allure.Step(t, "insideBefore", func(stepT cute.T) {
					time.Sleep(time.Second)
					now := time.Now()
					stepT.Logf("Test 1. Start time %v", now)
					stepT.Parameters(allure.NewParameter("Test 1. Time", now))
				})

				return nil
			},
		).
		AfterExecuteT(
			func(t cute.T, resp *http.Response, errs []error) error {
				allure.Step(t, "insideAfter", func(stepT cute.T) {
					now := time.Now()
					stepT.Logf("Test 1. Stop time %v", now)
					stepT.Parameters(allure.NewParameter("Test 1. Stop time", now))
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
				allure.Step(t, "insideBefore", func(stepT cute.T) {
					now := time.Now()
					stepT.Logf("Test 2. Start time %v", now)
					stepT.Parameters(allure.NewParameter("Test 2. Start time", now))
					time.Sleep(2 * time.Second)
				})

				return nil
			},
		).
		AfterExecuteT(
			func(t cute.T, resp *http.Response, errs []error) error {
				allure.Step(t, "insideAfter", func(stepT cute.T) {
					now := time.Now()
					stepT.Logf("test 2. Stop time %v", now)
					stepT.Parameters(allure.NewParameter("Test 2. Stop time", now))
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
