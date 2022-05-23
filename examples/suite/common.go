package suite

import (
	"net/url"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/require"

	"github.com/ozontech/cute"
)

type ExampleSuite struct {
	suite.Suite
	host *url.URL

	testMaker *cute.HTTPTestMaker
}

func (i *ExampleSuite) BeforeAll(t provider.T) {
	// Prepare http test builder
	i.testMaker = cute.NewHTTPTestMaker()

	// Preparing host
	host, err := url.Parse("https://jsonplaceholder.typicode.com/")
	require.NoError(t, err)

	i.host = host
}

func (i *ExampleSuite) BeforeEach(t provider.T) {
	t.Feature("ExampleSuite")
	t.Tags("some_global_tag")
}
