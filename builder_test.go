package cute

import (
	"net/http"
	"testing"
	"time"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/stretchr/testify/require"
)

func TestNewTestBuilder(t *testing.T) {
	var (
		maker = NewHTTPTestMaker()
		ht    = maker.NewTestBuilder().(*cute)
	)

	require.NotNil(t, ht.tests)
	require.Len(t, ht.tests, 1)
	require.NotNil(t, ht.tests[0].request)
	require.NotNil(t, ht.tests[0].middleware)
	require.NotNil(t, ht.tests[0].allureStep)
	require.NotNil(t, ht.allureInfo)
	require.NotNil(t, ht.httpClient)
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
		owner          = "Owner"
		lead           = "Lead"
		label          = allure.Label{"kek", "lol"}
		setIssue       = "SetIssue"
		setTestCase    = "SetTestCase"
		repeatCount    = 10
		repeatDelay    = time.Duration(10)
		link           = allure.Link{
			Name: "link",
			Type: "type",
			URL:  "http://go.go",
		}
		labels = []allure.Label{
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
		Owner(owner).
		Lead(lead).
		Label(label).
		Labels(labels...).
		SetIssue(setIssue).
		SetTestCase(setTestCase).
		Link(link).
		Description(desc).
		CreateWithStep().
		StepName(stepName).
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
		AssertResponseT(assertResponseT...)

	resHt := ht.(*cute)
	require.Equal(t, title, resHt.allureInfo.title)
	require.Equal(t, tags, resHt.allureLabels.tags)
	require.Equal(t, desc, resHt.allureInfo.description)
	require.Equal(t, feature, resHt.allureLabels.feature)
	require.Equal(t, epic, resHt.allureLabels.epic)
	require.Equal(t, stepName, resHt.tests[0].allureStep.name)
	require.Equal(t, req, resHt.tests[0].request.base)
	require.Equal(t, executeTime, resHt.tests[0].expect.ExecuteTime)
	require.Equal(t, status, resHt.tests[0].expect.Code)
	require.Equal(t, schemaBt, resHt.tests[0].expect.JSONSchemaByte)
	require.Equal(t, schemaStg, resHt.tests[0].expect.JSONSchemaString)
	require.Equal(t, schemaFile, resHt.tests[0].expect.JSONSchemaFile)
	require.Equal(t, id, resHt.allureLabels.id)
	require.Equal(t, addSuiteLabel, resHt.allureLabels.suiteLabel)
	require.Equal(t, addSubSuite, resHt.allureLabels.subSuite)
	require.Equal(t, addParentSuite, resHt.allureLabels.parentSuite)
	require.Equal(t, story, resHt.allureLabels.story)
	require.Equal(t, tag, resHt.allureLabels.tag)
	require.Equal(t, owner, resHt.allureLabels.owner)
	require.Equal(t, lead, resHt.allureLabels.lead)
	require.Equal(t, label, resHt.allureLabels.label)
	require.Equal(t, setIssue, resHt.allureLinks.issue)
	require.Equal(t, setTestCase, resHt.allureLinks.testCase)
	require.Equal(t, link, resHt.allureLinks.link)
	require.Equal(t, repeatCount, resHt.tests[0].request.repeat.count)
	require.Equal(t, repeatDelay, resHt.tests[0].request.repeat.delay)

	require.Equal(t, len(assertHeaders), len(resHt.tests[0].expect.AssertHeaders))
	require.Equal(t, len(assertHeadersT), len(resHt.tests[0].expect.AssertHeadersT))

	require.Equal(t, len(assertBody), len(resHt.tests[0].expect.AssertBody))
	require.Equal(t, len(assertBodyT), len(resHt.tests[0].expect.AssertBodyT))

	require.Equal(t, len(assertResponse), len(resHt.tests[0].expect.AssertResponse))
	require.Equal(t, len(assertResponseT), len(resHt.tests[0].expect.AssertResponseT))
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
