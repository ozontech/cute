package suite

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/cute"
	"github.com/ozontech/cute/examples"
)

/*

	Example testing HTTP POST, validate body and make second request.
	Validate:
		1) Execute time
		2) Status code

*/

func (i *ExampleSuite) TestExample_TwoSteps(t provider.T) {
	var (
		testBuilder = i.testMaker.NewTestBuilder()

		// Request body
		r = `
{
    "result": {
        "author": "Yours Truly",
        "date": "15.11.1993",
        "slides": [
            {
                "title": "Beer",
                "type": "drink"
            },
            {
                "title": "Apple",
                "type": "fruit"
            },
            {
                "title": "Orange",
                "type": "fruit"
            }
        ],
        "Info": {
            "shop": "BigShopPlus",
            "address": "address"
        },
        "title": "Sample Show"
    }
}
	`
	)

	u, _ := url.Parse(i.host.String())
	u.Path = path.Join(u.Path, "/posts/1/comments")

	req, _ := http.NewRequest(http.MethodPost, u.String(), ioutil.NopCloser(strings.NewReader(r)))
	req.Header = map[string][]string{
		"some_auth_token": []string{fmt.Sprint(11111)},
	}

	testBuilder.
		Title("TestExample_TwoSteps").
		Tags("TestExample_TwoSteps", "some_tag").
		Parallel().
		CreateWithStep().

		// CreateWithStep first step

		StepName("Creat entry /posts/1").
		Request(req).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusCreated).
		AssertBody(
			// Custom assert body
			examples.CustomAssertBody(),
		).
		ExecuteTest(context.Background(), t).

		// CreateWithStep second step for delete

		NextTestWithStep().
		StepName("Delete entry").
		RequestBuilder(
			cute.WithURL(u),
			cute.WithMethod(http.MethodDelete),
			cute.WithHeaders(map[string][]string{
				"some_auth_token": []string{fmt.Sprint(11111)},
			}),
		).
		ExecuteTest(context.Background(), t)
}
