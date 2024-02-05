package cute

import (
	"net/http"
	"testing"
	"time"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/stretchr/testify/require"
)

func TestBuilderAfterTest(t *testing.T) {
	var (
		maker = NewHTTPTestMaker()
	)

	ht := maker.NewTestBuilder().
		Create().
		RequestBuilder().
		NextTest().
		AfterTestExecute(
			func(response *http.Response, errors []error) error {
				return nil
			},
			func(response *http.Response, errors []error) error {
				return nil
			}).
		AfterTestExecuteT(
			func(t T, response *http.Response, errors []error) error {

				return nil
			},
			func(t T, response *http.Response, errors []error) error {

				return nil
			},
			func(t T, response *http.Response, errors []error) error {

				return nil
			},
		)

	res := ht.(*cute)
	require.Len(t, res.tests[0].Middleware.After, 2)
	require.Len(t, res.tests[0].Middleware.AfterT, 3)
}

func TestBuilderAfterTestTwoStep(t *testing.T) {
	var (
		maker = NewHTTPTestMaker(
			WithMiddlewareBefore(
				func(request *http.Request) error {
					return nil
				},
				func(request *http.Request) error {
					return nil
				},
			),
			WithMiddlewareBeforeT(
				func(t T, request *http.Request) error {
					return nil
				},
			),
			WithMiddlewareAfter(
				func(response *http.Response, errors []error) error {
					return nil
				},
			),
			WithMiddlewareAfterT(
				func(t T, response *http.Response, errors []error) error {
					return nil
				},
				func(t T, response *http.Response, errors []error) error {
					return nil
				},
				func(t T, response *http.Response, errors []error) error {
					return nil
				},
			),
		)
	)

	ht :=
		maker.NewTestBuilder().
			Create().
			RequestBuilder().
			NextTest().
			AfterTestExecute(
				func(response *http.Response, errors []error) error {
					return nil
				},
				func(response *http.Response, errors []error) error {
					return nil
				}).
			AfterTestExecuteT(
				func(t T, response *http.Response, errors []error) error {

					return nil
				},
				func(t T, response *http.Response, errors []error) error {

					return nil
				},
				func(t T, response *http.Response, errors []error) error {

					return nil
				},
			).
			Create().
			AfterExecute(
				func(response *http.Response, errors []error) error {

					return nil
				},
			).
			AfterExecuteT(
				func(t T, response *http.Response, errors []error) error {

					return nil
				}).
			RequestBuilder().
			NextTest().
			AfterTestExecute(
				func(response *http.Response, errors []error) error {

					return nil
				},
			)

	res := ht.(*cute)
	require.Len(t, res.tests[0].Middleware.After, 2+1)
	require.Len(t, res.tests[0].Middleware.Before, 2)
	require.Len(t, res.tests[0].Middleware.BeforeT, 1)
	require.Len(t, res.tests[0].Middleware.AfterT, 3+3)

	require.Len(t, res.tests[1].Middleware.After, 2+1)
	require.Len(t, res.tests[1].Middleware.AfterT, 1+3)
	require.Len(t, res.tests[1].Middleware.Before, 2)
	require.Len(t, res.tests[1].Middleware.BeforeT, 1)
}

func TestNewTestBuilder(t *testing.T) {
	var (
		maker = NewHTTPTestMaker()
		ht    = maker.NewTestBuilder().(*cute)
	)

	require.NotNil(t, ht.tests)
	require.Len(t, ht.tests, 1)
	require.NotNil(t, ht.tests[0].Request)
	require.NotNil(t, ht.tests[0].Middleware)
	require.NotNil(t, ht.tests[0].AllureStep)
	require.NotNil(t, ht.allureInfo)
	require.NotNil(t, ht.baseProps.httpClient)
}

func TestHTTPTestMaker(t *testing.T) {
	var (
		maker          = NewHTTPTestMaker()
		ht             = maker.NewTestBuilder()
		title          = "title"
		epic           = "epic"
		desc           = "desc"
		feature        = "feature"
		tags           = []string{"tag_1", "tag_2"}
		stepName       = "stepname"
		req, _         = http.NewRequest(http.MethodGet, "https://site.go", nil)
		executeTime    = time.Duration(10)
		status         = 400
		schemaStg      = "some_json_schema"
		schemaBt       = []byte("some_json_schema")
		schemaFile     = "file_path"
		id             = "ID"
		addSuiteLabel  = "AddSuiteLabel"
		addSubSuite    = "AddSubSuite"
		addParentSuite = "AddParentSuite"
		story          = "Story"
		tag            = "Tag"
		allureID       = "AllureID"
		owner          = "Owner"
		lead           = "Lead"
		label          = &allure.Label{Name: "kek", Value: "lol"}
		setIssue       = "SetIssue"
		setTestCase    = "SetTestCase"
		repeatCount    = 10
		repeatDelay    = time.Duration(10)
		link           = &allure.Link{
			Name: "link",
			Type: "type",
			URL:  "http://go.go",
		}
		labels = []*allure.Label{
			{
				Name:  "label_1",
				Value: "value_1",
			},
			{
				Name:  "label_2",
				Value: "value_2",
			},
		}

		assertHeaders = []AssertHeaders{
			func(headers http.Header) error {
				return nil
			},
		}
		assertHeadersT = []AssertHeadersT{
			func(t T, headers http.Header) error {
				return nil
			},
			func(t T, headers http.Header) error {
				return nil
			},
		}

		assertBody = []AssertBody{
			func(body []byte) error {
				return nil
			},
		}
		assertBodyT = []AssertBodyT{
			func(t T, body []byte) error {
				return nil
			},
			func(t T, body []byte) error {
				return nil
			},
		}

		assertResponse = []AssertResponse{
			func(resp *http.Response) error {
				return nil
			},
		}
		assertResponseT = []AssertResponseT{
			func(t T, resp *http.Response) error {
				return nil
			},
			func(t T, resp *http.Response) error {
				return nil
			},
		}
		after = []AfterExecute{
			func(response *http.Response, errors []error) error {

				return nil
			},
			func(response *http.Response, errors []error) error {

				return nil
			},
		}
		afterT = []AfterExecuteT{
			func(t T, response *http.Response, errors []error) error {

				return nil
			},
			func(t T, response *http.Response, errors []error) error {

				return nil
			},
		}
	)

	ht.
		Title(title).
		Tags(tags...).
		Epic(epic).
		Feature(feature).
		ID(id).
		AddSuiteLabel(addSuiteLabel).
		AddSubSuite(addSubSuite).
		AddParentSuite(addParentSuite).
		Story(story).
		Tag(tag).
		Severity(allure.CRITICAL).
		AllureID(allureID).
		Owner(owner).
		Lead(lead).
		Label(label).
		Labels(labels...).
		SetIssue(setIssue).
		SetTestCase(setTestCase).
		Link(link).
		Description(desc).
		CreateStep(stepName).
		RequestRepeat(repeatCount).
		RequestRepeatDelay(repeatDelay).
		Request(req).
		ExpectExecuteTimeout(executeTime).
		ExpectStatus(status).
		ExpectJSONSchemaByte(schemaBt).
		ExpectJSONSchemaString(schemaStg).
		ExpectJSONSchemaFile(schemaFile).
		AssertHeaders(assertHeaders...).
		AssertHeadersT(assertHeadersT...).
		AssertBody(assertBody...).
		AssertBodyT(assertBodyT...).
		AssertResponse(assertResponse...).
		AssertResponseT(assertResponseT...).
		After(after...).
		AfterT(afterT...)

	resHt := ht.(*cute)
	resTest := resHt.tests[0]

	require.Equal(t, title, resHt.allureInfo.title)
	require.Equal(t, tags, resHt.allureLabels.tags)
	require.Equal(t, desc, resHt.allureInfo.description)
	require.Equal(t, feature, resHt.allureLabels.feature)
	require.Equal(t, epic, resHt.allureLabels.epic)
	require.Equal(t, stepName, resTest.AllureStep.Name)
	require.Equal(t, req, resTest.Request.Base)
	require.Equal(t, executeTime, resTest.Expect.ExecuteTime)
	require.Equal(t, status, resTest.Expect.Code)
	require.Equal(t, schemaBt, resTest.Expect.JSONSchema.Byte)
	require.Equal(t, schemaStg, resTest.Expect.JSONSchema.String)
	require.Equal(t, schemaFile, resTest.Expect.JSONSchema.File)
	require.Equal(t, id, resHt.allureLabels.id)
	require.Equal(t, addSuiteLabel, resHt.allureLabels.suiteLabel)
	require.Equal(t, addSubSuite, resHt.allureLabels.subSuite)
	require.Equal(t, addParentSuite, resHt.allureLabels.parentSuite)
	require.Equal(t, story, resHt.allureLabels.story)
	require.Equal(t, tag, resHt.allureLabels.tag)
	require.Equal(t, owner, resHt.allureLabels.owner)
	require.Equal(t, lead, resHt.allureLabels.lead)
	require.Equal(t, label, resHt.allureLabels.label)
	require.Equal(t, allureID, resHt.allureLabels.allureID)
	require.Equal(t, setIssue, resHt.allureLinks.issue)
	require.Equal(t, setTestCase, resHt.allureLinks.testCase)
	require.Equal(t, link, resHt.allureLinks.link)
	require.Equal(t, repeatCount, resTest.Request.Repeat.Count)
	require.Equal(t, repeatDelay, resTest.Request.Repeat.Delay)

	require.Equal(t, len(assertHeaders), len(resTest.Expect.AssertHeaders))
	require.Equal(t, len(assertHeadersT), len(resTest.Expect.AssertHeadersT))

	require.Equal(t, len(assertBody), len(resTest.Expect.AssertBody))
	require.Equal(t, len(assertBodyT), len(resTest.Expect.AssertBodyT))

	require.Equal(t, len(assertResponse), len(resTest.Expect.AssertResponse))
	require.Equal(t, len(assertResponseT), len(resTest.Expect.AssertResponseT))

	require.Equal(t, len(after), len(resTest.Middleware.After))
	require.Equal(t, len(afterT), len(resTest.Middleware.AfterT))
}

func TestCreateDefaultTest(t *testing.T) {
	resTest := createDefaultTest(&HTTPTestMaker{httpClient: http.DefaultClient, middleware: new(Middleware)})

	require.Equal(t, &Test{
		httpClient: http.DefaultClient,
		Name:       "",
		AllureStep: new(AllureStep),
		Middleware: &Middleware{
			After:   make([]AfterExecute, 0),
			AfterT:  make([]AfterExecuteT, 0),
			Before:  make([]BeforeExecute, 0),
			BeforeT: make([]BeforeExecuteT, 0),
		},
		Request: &Request{
			Repeat: new(RequestRepeatPolitic),
		},
		Expect: &Expect{
			JSONSchema: new(ExpectJSONSchema),
		},
	}, resTest)
}

func TestCreateTableTest(t *testing.T) {
	c := &cute{}
	c.CreateTableTest()

	require.True(t, c.isTableTest)
}

func TestPutNewTest(t *testing.T) {
	tests := make([]*Test, 1)
	tests[0] = createDefaultTest(&HTTPTestMaker{httpClient: http.DefaultClient, middleware: new(Middleware)})

	var (
		c = &cute{tests: tests, baseProps: &HTTPTestMaker{
			middleware: &Middleware{},
		}}
		reqOne, _    = http.NewRequest("GET", "URL_1", nil)
		expectOne    = &Expect{Code: 200}
		reqSecond, _ = http.NewRequest("POST", "URL_1", nil)
		expectSecond = &Expect{Code: 400}
	)

	c.PutNewTest("name_1", reqOne, expectOne)
	c.PutNewTest("name_2", reqSecond, expectSecond)

	require.Equal(t, c.tests[0].Name, "name_1")
	require.Equal(t, c.tests[0].Expect, expectOne)
	require.Equal(t, c.tests[0].Request.Base, reqOne)

	require.Equal(t, c.tests[1].Name, "name_2")
	require.Equal(t, c.tests[1].Expect, expectSecond)
	require.Equal(t, c.tests[1].Request.Base, reqSecond)
}

func TestPutTests(t *testing.T) {
	var (
		tests        = createDefaultTests(&HTTPTestMaker{httpClient: http.DefaultClient, middleware: new(Middleware)})
		c            = &cute{tests: tests}
		reqOne, _    = http.NewRequest("GET", "URL_1", nil)
		expectOne    = &Expect{Code: 200}
		reqSecond, _ = http.NewRequest("POST", "URL_1", nil)
		expectSecond = &Expect{Code: 400}
	)

	tests = append(tests,
		&Test{
			Name: "name_1",
			Request: &Request{
				Base: reqOne,
			},
			Expect: expectOne,
		},
		&Test{
			Name: "name_2",
			Request: &Request{
				Base: reqSecond,
			},
			Expect: expectSecond,
		},
	)

	c.PutTests(tests...)

	require.Equal(t, c.tests[0].Name, "name_1")
	require.Equal(t, c.tests[0].Expect, expectOne)
	require.Equal(t, c.tests[0].Request.Base, reqOne)

	require.Equal(t, c.tests[1].Name, "name_2")
	require.Equal(t, c.tests[1].Expect, expectSecond)
	require.Equal(t, c.tests[1].Request.Base, reqSecond)
}

func TestCreateHTTPTestMakerWithHttpClient(t *testing.T) {
	cli := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       100,
	}

	maker := NewHTTPTestMaker(WithHTTPClient(cli))

	require.Equal(t, cli, maker.httpClient)
	require.Equal(t, time.Duration(100), maker.httpClient.Timeout)
}

type rt struct {
}

func (r *rt) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, nil
}

func TestCreateHTTPMakerOps(t *testing.T) {
	timeout := time.Second * 100
	roundTripper := &rt{}

	maker := NewHTTPTestMaker(
		WithCustomHTTPTimeout(timeout),
		WithCustomHTTPRoundTripper(roundTripper),
	)

	require.Equal(t, timeout, maker.httpClient.Timeout)
	require.Equal(t, roundTripper, maker.httpClient.Transport)
}
