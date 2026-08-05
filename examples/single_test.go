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

	cuteErrors "github.com/ozontech/cute/errors"
	"github.com/ozontech/testo"
	allure "github.com/ozontech/testo-allure"

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
		RequestRetry(3).
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
			allure.Step(t, "inside Assert body. 1 ", func(stepT cute.T) {
				stepT.Parameters(allure.NewParameter("key", "value"))
			})

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

func Test_Single_Broken(t *testing.T) {
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
			return errors.New("it's NOT must be run")
		},
		).
		ExecuteTest(context.Background(), t)

	t.Skip()
}

func Test_Single_RepeatPolitic_Optional_Success_Test(t *testing.T) {
	cute.NewTestBuilder().
		Title("Test_Single_RepeatPolitic_Optional_Success_Test").
		Create().
		RequestRetry(2).
		RequestRetryOptional(true).
		RequestBuilder(
			cute.WithURI("https://jsonplaceholder.typicode.com/posts/1/comments"),
		).
		BrokenAssertBodyT(func(t cute.T, body []byte) error {
			return errors.New("example broken error")
		}).
		ExpectStatus(http.StatusCreated).
		ExecuteTest(context.Background(), t)

	t.Logf("You should see it")
}

func Test_Single_RepeatPolitic_Broken_Failed_Test(t *testing.T) {
	cute.NewTestBuilder().
		Title("Test_Single_RepeatPolitic_Broken_Failed_Test").
		Create().
		RequestRetry(2).
		RequestRetryOptional(false).
		RequestBuilder(
			cute.WithURI("https://jsonplaceholder.typicode.com/posts/1/comments"),
		).
		BrokenAssertBodyT(func(t cute.T, body []byte) error {
			return errors.New("example broken error")
		}).
		ExpectStatus(http.StatusCreated).
		ExecuteTest(context.Background(), t)

	t.Logf("You should see it")
}

func Test_Single_Broken_2(t *testing.T) {
	cute.NewTestBuilder().
		Title("Test_Single_Broken_2").
		Create().
		RequestBuilder(
			cute.WithURI("https://jsonplaceholder.typicode.com/posts/1/comments"),
		).
		AssertBodyT(func(t cute.T, body []byte) error {
			err := errors.New("example broken error")
			return cuteErrors.WrapBrokenError(err)
		}).
		ExpectStatus(http.StatusOK).
		NextTest().
		Create().
		RequestBuilder(
			cute.WithURI("https://jsonplaceholder.typicode.com/posts/1/comments"),
		).
		AssertBody(func(body []byte) error {
			return errors.New("it's NOT must be run")
		},
		).
		ExecuteTest(context.Background(), t)
}

func Test_Single_2_TestoRunner(t *testing.T) {
	t.Run("Single test with testo Runner", testo.Test(func(t T) {
		var (
			testMaker   = cute.NewHTTPTestMaker()
			testBuilder = testMaker.NewTestBuilder()
		)

		u, _ := url.Parse("https://jsonplaceholder.typicode.com/")
		u.Path = path.Join(u.Path, "/posts/1/comments")

		testBuilder.
			Title("Single test with testo T and repeat errors").
			Tag("single_test").
			Description("some_description").
			Create().
			RequestRetryDelay(3*time.Second). // delay before new try
			RequestRetry(3).                  // count attempts
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
	}))
}
