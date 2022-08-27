//go:build example
// +build example

package table_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ozontech/allure-go/pkg/framework/core/allure_manager/manager"
	"github.com/ozontech/allure-go/pkg/framework/core/common"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/cute"
	"github.com/ozontech/cute/asserts/json"
)

func init() {
	os.Setenv("ALLURE_OUTPUT_PATH", "../") // custom, read Readme.md for more info
}

func TestTableExample(t *testing.T) {
	u, _ := url.Parse("https://jsonplaceholder.typicode.com/posts/1/comments")

	req, _ := http.NewRequest(http.MethodPost, u.String(), nil)
	req.Header = map[string][]string{
		"some_auth_token": []string{fmt.Sprint(11111)},
	}
	cute.NewTestBuilder().
		//Title("Example table test").
		CreateTableTest().
		PutTest("Execute validation 1", req, &cute.Expect{
			Code: 201,
		}).
		PutTest(
			"Execute validation 2",
			req,
			&cute.Expect{
				AssertBody: []cute.AssertBody{
					json.Equal("$[0].email", "Eliseo@gardner.biz"),
					json.Present("$[1].nam1e"),
				},
			},
		).
		ExecuteTest(context.Background(), t)
}

func TestExampleProblem(t *testing.T) {
	allureT := createAllureT(t)

	allureT.Title("Title")

	allureT.Logf("First log")

	allureT.Run("insideRun", func(inT provider.T) {
		inT.Logf("First log from inside")
		time.Sleep(time.Second)
		inT.Logf("Last log from inside")
	})

	allureT.Logf("Last log")
}

func createAllureT(t *testing.T) *common.Common {
	var (
		newT        = common.NewT(t)
		callers     = strings.Split(t.Name(), "/")
		providerCfg = manager.NewProviderConfig().
				WithFullName("t.Name()").
				WithPackageName("package").
				WithRunner(callers[0])
		newProvider = manager.NewProvider(providerCfg)
	)
	newProvider.NewTest(t.Name(), "package")

	newT.SetProvider(newProvider)
	newT.Provider.TestContext()

	return newT
}
