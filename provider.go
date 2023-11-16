package cute

import (
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type T interface {
	tProvider
	logProvider
	stepProvider
	attachmentProvider
	parametersProvider
}

type allureProvider interface {
	internalT
	Parallel()
	Run(testName string, testBody func(provider.T), tags ...string) (res *allure.Result)

	infoAllureProvider
	labelsAllureProvider
	linksAllureProvider
}

type internalT interface {
	tProvider
	logProvider
	stepProvider
	attachmentProvider
}

type tProvider interface {
	Fail()
	FailNow()

	Name() string

	Log(args ...interface{})
	Logf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

type logProvider interface {
	LogStep(args ...interface{})
	LogfStep(format string, args ...interface{})
}

type stepProvider interface {
	Step(step *allure.Step)
	WithNewStep(stepName string, step func(ctx provider.StepCtx), params ...*allure.Parameter)
}

type attachmentProvider interface {
	WithAttachments(attachments ...*allure.Attachment)
	WithNewAttachment(name string, mimeType allure.MimeType, content []byte)
}

type parametersProvider interface {
	WithParameters(parameters ...*allure.Parameter)
	WithNewParameters(kv ...interface{})
}

type infoAllureProvider interface {
	Title(args ...interface{})
	Titlef(format string, args ...interface{})

	Description(args ...interface{})
	Descriptionf(format string, args ...interface{})

	Stage(args ...interface{})
	Stagef(format string, args ...interface{})
}

type labelsAllureProvider interface {
	ID(value string)
	AllureID(value string)
	Epic(value string)
	Layer(value string)
	AddSuiteLabel(value string)
	AddSubSuite(value string)
	AddParentSuite(value string)
	Feature(value string)
	Story(value string)
	Tag(value string)
	Tags(values ...string)
	Severity(value allure.SeverityType)
	Owner(value string)
	Lead(value string)
	Label(label *allure.Label)
	Labels(labels ...*allure.Label)
}

type linksAllureProvider interface {
	SetIssue(issue string)
	SetTestCase(testCase string)
	Link(link *allure.Link)
	TmsLink(tmsCase string)
	TmsLinks(tmsCases ...string)
}
