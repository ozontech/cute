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

	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

type logProvider interface {
	Log(args ...interface{})
	Logf(format string, args ...interface{})
}

type stepProvider interface {
	Step(step *allure.Step)
	WithNewStep(stepName string, step func(ctx provider.StepCtx), params ...allure.Parameter)
}

type attachmentProvider interface {
	WithAttachments(attachments ...*allure.Attachment)
	WithNewAttachment(name string, mimeType allure.MimeType, content []byte)
}

type parametersProvider interface {
	WithParameters(parameters ...allure.Parameter)
	WithNewParameters(kv ...interface{})
}

type infoAllureProvider interface {
	Title(title string)
	Description(description string)
}

type labelsAllureProvider interface {
	ID(value string)
	Epic(value string)
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
	Label(label allure.Label)
	Labels(labels ...allure.Label)
}

type linksAllureProvider interface {
	SetIssue(issue string)
	SetTestCase(testCase string)
	Link(link allure.Link)
}
