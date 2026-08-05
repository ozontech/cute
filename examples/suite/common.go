package suite

import (
	"net/url"

	"github.com/ozontech/testo"
	"github.com/stretchr/testify/require"

	"github.com/ozontech/cute"
)

type ExampleSuite struct {
	testo.Suite[T]
	host *url.URL

	testMaker *cute.HTTPTestMaker
}

func (i *ExampleSuite) BeforeAll(t T) {
	// Prepare http test builder
	i.testMaker = cute.NewHTTPTestMaker()

	// Preparing host
	host, err := url.Parse("https://jsonplaceholder.typicode.com/")
	require.NoError(t, err)

	i.host = host
}

func (i *ExampleSuite) BeforeEach(t T) {
	t.Feature("ExampleSuite")
	t.Tags("some_global_tag")
}
