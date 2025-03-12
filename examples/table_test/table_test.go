//go:build example
// +build example

package table_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ozontech/cute"
	"github.com/ozontech/cute/asserts/json"
	"github.com/ozontech/cute/errors"
)

func init() {
	os.Setenv("ALLURE_OUTPUT_PATH", "../") // custom, read Readme.md for more info
}

func Test_Table(t *testing.T) {
	u, _ := url.Parse("https://jsonplaceholder.typicode.com/posts/1/comments")

	req, _ := http.NewRequest(http.MethodPost, u.String(), nil)
	req.Header = map[string][]string{
		"some_auth_token": []string{fmt.Sprint(11111)},
	}
	cute.NewTestBuilder().
		Title("Example put tests in table test").
		Tag("table_test").
		CreateTableTest().
		PutNewTest(
			"Execute validation 1",
			req,
			&cute.Expect{
				Code: 201,
			}).
		PutNewTest(
			"Execute validation 2",
			req,
			&cute.Expect{
				AssertBody: []cute.AssertBody{
					json.Equal("$[0].email", "Eliseo@gardner.biz"),
					json.Present("$[1].name"),
				},
			},
		).
		ExecuteTest(context.Background(), t)
}

func Test_Table_Array(t *testing.T) {
	tests := []*cute.Test{
		{
			Name:       "Create something",
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("https://jsonplaceholder.typicode.com/posts/1/comments"),
					cute.WithMethod(http.MethodPost),
				},
			},
			Expect: &cute.Expect{
				Code: 200,
			},
		},
		{
			Name:       "Delete something",
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("https://jsonplaceholder.typicode.com/posts/1/comments"),
					cute.WithMethod(http.MethodGet),
				},
			},
			Expect: &cute.Expect{
				Code: 200,
				AssertBody: []cute.AssertBody{
					json.Equal("$[0].email", "Eliseo@gardner.biz"),
					json.Present("$[1].name"),
					func(body []byte) error {
						return errors.NewAssertError("example error", "example message", nil, nil)
					},
				},
			},
		},
	}

	cute.NewTestBuilder().
		Title("Example table test").
		Tag("table_test").
		Description("Execute array tests").
		CreateTableTest().
		PutTests(tests...).
		ExecuteTest(context.Background(), t)
}

func Test_One_Execute(t *testing.T) {
	test := &cute.Test{
		Name: "test_1",
		Request: &cute.Request{
			Builders: []cute.RequestBuilder{
				cute.WithURI("https://jsonplaceholder.typicode.com/posts/1/comments"),
				cute.WithMethod(http.MethodGet),
			},
		},
		Expect: nil,
	}

	test.Execute(context.Background(), t)
}

func Test_Array_Retry_OptionalFirstTries(t *testing.T) {
	tests := []*cute.Test{
		{
			Name: "test_1",

			Retry: &cute.Retry{
				MaxAttempts: 10,
				Delay:       1 * time.Second,
			},
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("https://httpstat.us/Random/201,202"),
					cute.WithMethod(http.MethodGet),
				},
			},
			Expect: &cute.Expect{
				Code: 201,
			},
		},
		{
			Name: "test_2",
			Retry: &cute.Retry{
				MaxAttempts: 10,
				Delay:       1 * time.Second,
			},
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("https://httpstat.us/Random/403,404"),
					cute.WithMethod(http.MethodGet),
					cute.WithMarshalBody([]byte("{\"test\":\"abc\"}")),
				},
			},
			Expect: &cute.Expect{
				Code: 404,
			},
		},
	}

	for _, test := range tests {
		test.Execute(context.Background(), t)
	}
}

func Test_Array_Retry_OptionalFirstTries_UltimatelyFailing(t *testing.T) {
	tests := []*cute.Test{
		{
			Name: "test_1",

			Retry: &cute.Retry{
				MaxAttempts: 4,
				Delay:       1 * time.Second,
			},
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("https://httpstat.us/Random/202,200"),
					cute.WithMethod(http.MethodGet),
				},
			},
			Expect: &cute.Expect{
				Code: 201,
			},
		},
		{
			Name: "test_2",
			Retry: &cute.Retry{
				MaxAttempts: 3,
				Delay:       1 * time.Second,
			},
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("https://httpstat.us/Random/403,401"),
					cute.WithMethod(http.MethodGet),
				},
			},
			Expect: &cute.Expect{
				Code: 404,
			},
		},
	}

	for _, test := range tests {
		test.Execute(context.Background(), t)
	}
}

func Test_Array_TimeoutRetry(t *testing.T) {
	var executeTimeout = 3000

	tests := []*cute.Test{
		{
			Retry: &cute.Retry{
				MaxAttempts: 2,
			},
			Name: "test_timeout",
			Middleware: &cute.Middleware{
				Before: []cute.BeforeExecute{
					cute.BeforeExecute(func(request *http.Request) error {
						query := request.URL.Query()
						query.Set("sleep", strconv.Itoa(executeTimeout))
						request.URL.RawQuery = query.Encode()
						executeTimeout = executeTimeout - 1000
						return nil
					}),
				},
			},
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("https://httpstat.us/202?sleep=3000"),
					cute.WithBody([]byte("{\"test\":\"abc\"}")),
					cute.WithMethod(http.MethodGet),
				},
			},
			Expect: &cute.Expect{
				Code:        202,
				ExecuteTime: 3 * time.Second,
			},
		},
	}

	for _, test := range tests {
		test.Execute(context.Background(), t)
	}
}

func Test_Array(t *testing.T) {
	tests := []*cute.Test{
		{
			Name:       "test_1",
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("https://jsonplaceholder.typicode.com/posts/1/comments"),
					cute.WithMethod(http.MethodPost),
				},
			},
			Expect: &cute.Expect{
				Code: 201,
			},
		},
		{
			Name:       "test_2",
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("https://jsonplaceholder.typicode.com/posts/1/comments"),
					cute.WithMethod(http.MethodGet),
				},
			},
			Expect: &cute.Expect{
				Code: 200,
				AssertBody: []cute.AssertBody{
					json.Equal("$[0].email", "Eliseo@gardner.biz"),
					json.Present("$[1].name"),
					func(body []byte) error {
						return errors.NewAssertError("example error", "example message", nil, nil)
					},
				},
			},
		},
	}

	for _, test := range tests {
		test.Execute(context.Background(), t)
	}
}

func Test_Array_All_Parallel(t *testing.T) {
	tests := []*cute.Test{
		{
			Name:       "test_201",
			Parallel:   true,
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("https://httpstat.us/201"),
					cute.WithMethod(http.MethodGet),
				},
			},
			Expect: &cute.Expect{
				Code: 201,
			},
		},
		{
			Name:       "test_200_delay_5s",
			Parallel:   true,
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("https://httpstat.us/200?sleep=5000"),
					cute.WithMethod(http.MethodGet),
				},
			},
			Expect: &cute.Expect{
				Code: 200,
			},
		},
		{
			Name:       "test_202_delay_3s",
			Parallel:   true,
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("https://httpstat.us/202?sleep=3000"),
					cute.WithMethod(http.MethodGet),
				},
			},
			Expect: &cute.Expect{
				Code: 202,
			},
		},
		{
			Name:       "test_203",
			Parallel:   true,
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("https://httpstat.us/203"),
					cute.WithMethod(http.MethodGet),
				},
			},
			Expect: &cute.Expect{
				Code: 203,
			},
		},
	}

	for _, test := range tests {
		test.Execute(context.Background(), t)
	}
}

func Test_Array_Some_Parallel(t *testing.T) {
	tests := []*cute.Test{
		{
			Name:       "test_parallel_1",
			Parallel:   true,
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("https://httpstat.us/201?sleep=1000"),
					cute.WithMethod(http.MethodGet),
				},
			},
			Expect: &cute.Expect{
				Code: 201,
			},
		},
		{
			Name:       "test_parallel_2",
			Parallel:   true,
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("https://httpstat.us/202?sleep=1000"),
					cute.WithMethod(http.MethodGet),
				},
			},
			Expect: &cute.Expect{
				Code: 202,
			},
		},
		{
			Name:       "test_1_sequential",
			Parallel:   false,
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("https://jsonplaceholder.typicode.com/posts/1/comments"),
					cute.WithMethod(http.MethodPost),
				},
			},
			Expect: &cute.Expect{
				Code: 201,
			},
		},
		{
			Name:       "test_2_sequential",
			Parallel:   false,
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("https://jsonplaceholder.typicode.com/posts/1/comments"),
					cute.WithMethod(http.MethodGet),
				},
			},
			Expect: &cute.Expect{
				Code: 200,
				AssertBody: []cute.AssertBody{
					json.Equal("$[0].email", "Eliseo@gardner.biz"),
					json.Present("$[1].name"),
				},
			},
		},
	}

	for _, test := range tests {
		test.Execute(context.Background(), t)
	}
}

func Test_Array_Retry(t *testing.T) {
	tests := []*cute.Test{
		{
			Name:     "test_1",
			Parallel: true,
			Retry: &cute.Retry{
				MaxAttempts: 10,
				Delay:       1,
			},
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("https://httpstat.us/Random/201,202"),
					cute.WithMethod(http.MethodGet),
				},
			},
			Expect: &cute.Expect{
				Code: 201,
			},
		},
		{
			Name:     "test_2",
			Parallel: true,
			Retry: &cute.Retry{
				MaxAttempts: 10,
				Delay:       1,
			},
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("https://httpstat.us/Random/403,404"),
					cute.WithMethod(http.MethodGet),
				},
			},
			Expect: &cute.Expect{
				Code: 404,
			},
		},
	}

	for _, test := range tests {
		test.Execute(context.Background(), t)
	}
}

func Test_Array_Timeout(t *testing.T) {
	tests := []*cute.Test{
		{
			Name:       "test_timeout",
			Middleware: nil,
			Request: &cute.Request{
				Builders: []cute.RequestBuilder{
					cute.WithURI("https://httpstat.us/202?sleep=3000"),
					cute.WithMethod(http.MethodGet),
				},
			},
			Expect: &cute.Expect{
				Code:        202,
				ExecuteTime: 2 * time.Second,
			},
		},
	}

	for _, test := range tests {
		test.Execute(context.Background(), t)
	}
}
