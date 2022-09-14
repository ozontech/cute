package suite

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/cute"
	"github.com/ozontech/cute/asserts/headers"
	"github.com/ozontech/cute/asserts/json"
	cuteErrors "github.com/ozontech/cute/errors"
	"github.com/ozontech/cute/examples"
)

func (i *ExampleSuite) Test_OneStep_Errors(t provider.T) {
	var (
		testBuilder = i.testMaker.NewTestBuilder()
	)

	testBuilder.
		Title("Test with errors").
		Tags("one_step", "some_local_tag", "suite", "json").
		Parallel().
		CreateStep("Example GET json request").
		RequestBuilder(
			cute.WithHeaders(map[string][]string{
				"some_header":       []string{"something"},
				"some_array_header": []string{"1", "2", "3", "some_thing"},
			}),
			cute.WithURI(i.host.String()+"/posts/1/comments"),
			cute.WithMethod(http.MethodGet),
			cute.WithMarshalBody(
				map[string]interface{}{
					"key": "value",
					"more_key": map[string]interface{}{
						"some_value": "sss",
					},
				},
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectJSONSchemaFile("file://./resources/example_valid_request.json").
		AssertBody(
			json.Equal("$[0].email", "something"),
			json.Present("$[1].not_present"),
			json.GreaterThan("$", 99999),
			json.Length("$", 0),
			// Custom assert body
			examples.CustomAssertBody(),
			examples.CustomAssertBodyWithCustomError(),
		).
		AssertHeaders(
			headers.Present("Content-Type"),

			// Custom assert headers
			examples.CustomAssertHeaders(),
		).
		AssertResponse(
			examples.CustomAssertResponse(),
		).
		AssertHeadersT(
			func(t cute.T, headers http.Header) error {
				// Example pretty print error
				return cuteErrors.NewAssertError("custom_assert", "example custom assert", "empty", "not empty") //
			},
		).
		// Example optional
		OptionalAssertBody( // example optional assert
			func(body []byte) error {
				return errors.New("some optional error from OptionalAssert")
			},
		).
		AssertBody(
			func(body []byte) error {
				return cuteErrors.NewOptionalError("some optional error from creator") // example optional error
			},
		).
		ExecuteTest(context.Background(), t)
}
