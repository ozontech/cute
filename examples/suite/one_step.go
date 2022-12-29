package suite

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/cute"
	"github.com/ozontech/cute/asserts/headers"
	"github.com/ozontech/cute/asserts/json"
	"github.com/ozontech/cute/examples"
)

/*
	Example testing HTTP GET and validate body.
	Validate:
		1) Execute time
		2) Status code
		3) Validate body by json schema
		4) Validate fields in json

Response:
[
  {
    "postId": 1,
    "id": 1,
    "name": "id labore ex et quam laborum",
    "email": "Eliseo@gardner.biz",
    "body": "laudantium enim quasi est quidem magnam voluptate ipsam eos\ntempora quo necessitatibus\ndolor quam autem quasi\nreiciendis et nam sapiente accusantium"
  },
  {
    "postId": 1,
    "id": 2,
    "name": "quo vero reiciendis velit similique earum",
    "email": "Jayne_Kuhic@sydney.com",
    "body": "est natus enim nihil est dolore omnis voluptatem numquam\net omnis occaecati quod ullam at\nvoluptatem error expedita pariatur\nnihil sint nostrum voluptatem reiciendis et"
  },
  {
    "postId": 1,
    "id": 3,
    "name": "odio adipisci rerum aut animi",
    "email": "Nikita@garfield.biz",
    "body": "quia molestiae reprehenderit quasi aspernatur\naut expedita occaecati aliquam eveniet laudantium\nomnis quibusdam delectus saepe quia accusamus maiores nam est\ncum et ducimus et vero voluptates excepturi deleniti ratione"
  },
  {
    "postId": 1,
    "id": 4,
    "name": "alias odio sit",
    "email": "Lew@alysha.tv",
    "body": "non et atque\noccaecati deserunt quas accusantium unde odit nobis qui voluptatem\nquia voluptas consequuntur itaque dolor\net qui rerum deleniti ut occaecati"
  },
  {
    "postId": 1,
    "id": 5,
    "name": "vero eaque aliquid doloribus et culpa",
    "email": "Hayden@althea.biz",
    "body": "harum non quasi et ratione\ntempore iure ex voluptates in ratione\nharum architecto fugit inventore cupiditate\nvoluptates magni quo et"
  }
]

*/

func (i *ExampleSuite) Test_OneStep(t provider.T) {
	var (
		testBuilder = i.testMaker.NewTestBuilder()
	)

	u, _ := url.Parse(i.host.String())
	u.Path = path.Join(u.Path, "/posts/1/comments")

	testBuilder.
		Title("Test with one step").
		Tags("one_stp", "some_local_tag", "suite", "json").
		Feature("some_feature").
		Epic("some_epic").
		Description("some_description").
		Parallel().
		CreateStep("Example GET json request").
		AfterExecuteT(func(t cute.T, resp *http.Response, errs []error) error {
			if len(errs) != 0 {
				return nil
			}

			/*
			 Implement some logic
			*/

			return nil
		},

			// After failed test
			func(t cute.T, resp *http.Response, errs []error) error {
				if len(errs) == 0 {
					return nil
				}

				/*
				 Implement some logic
				*/

				return nil
			},
		).
		RequestBuilder(
			cute.WithHeaders(map[string][]string{
				"some_header":       []string{"something"},
				"some_array_header": []string{"1", "2", "3", "some_thing"},
			}),
			cute.WithURL(u),
			cute.WithMethod(http.MethodGet),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectJSONSchemaFile("file://./resources/example_valid_request.json").
		ExpectStatus(http.StatusOK).
		AssertBody(
			json.Equal("$[0].email", "Eliseo@gardner.biz"),
			json.Present("$[1].name"),
			json.NotPresent("$[1].some_not_present"),
			json.LengthGreaterThan("$", 3),
			json.Length("$", 5),
			json.LengthLessThan("$", 100),
			json.NotEqual("$[3].name", "kekekekeke"),

			// Custom assert body
			examples.CustomAssertBody(),
		).
		AssertBodyT(
			// Custom assert body with testing.tb
			examples.CustomAssertBodyT(),

			func(t cute.T, body []byte) error {
				/*
					Implement here logic with TB
				*/
				time.Sleep(5 * time.Second)
				return nil
			},
		).
		AssertHeaders(
			headers.Present("Content-Type"),

			// Custom assert headers
			examples.CustomAssertHeaders(),
		).
		AssertResponse(
			examples.CustomAssertResponse(),
		).
		ExecuteTest(context.Background(), t)
}
