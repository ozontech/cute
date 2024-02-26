package cute

import (
	"fmt"

	"github.com/ozontech/allure-go/pkg/allure"
)

func (qt *cute) Parallel() AllureBuilder {
	qt.parallel = true

	return qt
}

func (qt *cute) Title(title string) AllureBuilder {
	qt.allureInfo.title = title

	return qt
}

func (qt *cute) Epic(epic string) AllureBuilder {
	qt.allureLabels.epic = epic

	return qt
}

func (qt *cute) Titlef(format string, args ...interface{}) AllureBuilder {
	qt.allureInfo.title = fmt.Sprintf(format, args...)

	return qt
}

func (qt *cute) Descriptionf(format string, args ...interface{}) AllureBuilder {
	qt.allureInfo.description = fmt.Sprintf(format, args...)

	return qt
}

func (qt *cute) Stage(stage string) AllureBuilder {
	qt.allureInfo.stage = stage

	return qt
}

func (qt *cute) Stagef(format string, args ...interface{}) AllureBuilder {
	qt.allureInfo.stage = fmt.Sprintf(format, args...)

	return qt
}

func (qt *cute) Layer(value string) AllureBuilder {
	qt.allureLabels.layer = value

	return qt
}

func (qt *cute) TmsLink(tmsLink string) AllureBuilder {
	qt.allureLinks.tmsLink = tmsLink

	return qt
}

func (qt *cute) TmsLinks(tmsLinks ...string) AllureBuilder {
	qt.allureLinks.tmsLinks = append(qt.allureLinks.tmsLinks, tmsLinks...)

	return qt
}

func (qt *cute) SetIssue(issue string) AllureBuilder {
	qt.allureLinks.issue = issue

	return qt
}

func (qt *cute) SetTestCase(testCase string) AllureBuilder {
	qt.allureLinks.testCase = testCase

	return qt
}

func (qt *cute) Link(link *allure.Link) AllureBuilder {
	qt.allureLinks.link = link

	return qt
}

func (qt *cute) ID(value string) AllureBuilder {
	qt.allureLabels.id = value

	return qt
}

func (qt *cute) AllureID(value string) AllureBuilder {
	qt.allureLabels.allureID = value

	return qt
}

func (qt *cute) AddSuiteLabel(value string) AllureBuilder {
	qt.allureLabels.suiteLabel = value

	return qt
}

func (qt *cute) AddSubSuite(value string) AllureBuilder {
	qt.allureLabels.subSuite = value

	return qt
}

func (qt *cute) AddParentSuite(value string) AllureBuilder {
	qt.allureLabels.parentSuite = value

	return qt
}

func (qt *cute) Story(value string) AllureBuilder {
	qt.allureLabels.story = value

	return qt
}

func (qt *cute) Tag(value string) AllureBuilder {
	qt.allureLabels.tag = value

	return qt
}

func (qt *cute) Severity(value allure.SeverityType) AllureBuilder {
	qt.allureLabels.severity = value

	return qt
}

func (qt *cute) Owner(value string) AllureBuilder {
	qt.allureLabels.owner = value

	return qt
}

func (qt *cute) Lead(value string) AllureBuilder {
	qt.allureLabels.lead = value

	return qt
}

func (qt *cute) Label(label *allure.Label) AllureBuilder {
	qt.allureLabels.label = label

	return qt
}

func (qt *cute) Labels(labels ...*allure.Label) AllureBuilder {
	qt.allureLabels.labels = labels

	return qt
}

func (qt *cute) Description(description string) AllureBuilder {
	qt.allureInfo.description = description

	return qt
}

func (qt *cute) Tags(tags ...string) AllureBuilder {
	qt.allureLabels.tags = tags

	return qt
}

func (qt *cute) Feature(feature string) AllureBuilder {
	qt.allureLabels.feature = feature

	return qt
}
