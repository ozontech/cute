package cute

import (
	"fmt"
	"os"
	"strings"

	allure "github.com/ozontech/testo-allure"
)

// Env keys with URL patterns for links, kept compatible with allure-go.
// A pattern must contain exactly one "%s", which is replaced by the ID.
const (
	issuePatternEnvKey    = "ALLURE_ISSUE_PATTERN"
	testCasePatternEnvKey = "ALLURE_TESTCASE_PATTERN"
	tmsPatternEnvKey      = "ALLURE_LINK_TMS_PATTERN"

	testCaseLinkType allure.LinkType = "test_case"
)

func linkFromPattern(envKey, name, id string, linkType allure.LinkType) allure.Link {
	pattern := os.Getenv(envKey)
	if !strings.Contains(pattern, "%s") {
		pattern = "%s"
	}

	return allure.Link{
		Name: name,
		URL:  fmt.Sprintf(pattern, id),
		Type: linkType,
	}
}

func (qt *cute) setAllureInformation(t T) {
	// Log main vars to allure
	qt.setLabelsAllure(t)
	qt.setInfoAllure(t)
	qt.setLinksAllure(t)
}

func (qt *cute) setLinksAllure(t T) {
	if issue := qt.allureLinks.issue; issue != "" {
		t.Links(linkFromPattern(issuePatternEnvKey, fmt.Sprintf("Issue[%s]", issue), issue, allure.LinkTypeIssue))
	}

	if testCase := qt.allureLinks.testCase; testCase != "" {
		t.Links(linkFromPattern(testCasePatternEnvKey, fmt.Sprintf("TestCase[%s]", testCase), testCase, testCaseLinkType))
	}

	if qt.allureLinks.link != nil {
		t.Links(*qt.allureLinks.link)
	}

	if tmsLink := qt.allureLinks.tmsLink; tmsLink != "" {
		t.Links(linkFromPattern(tmsPatternEnvKey, tmsLink, tmsLink, allure.LinkTypeTMS))
	}

	for _, tmsLink := range qt.allureLinks.tmsLinks {
		t.Links(linkFromPattern(tmsPatternEnvKey, tmsLink, tmsLink, allure.LinkTypeTMS))
	}
}

func (qt *cute) setLabelsAllure(t T) {
	if qt.allureLabels.id != "" {
		t.Labels(allure.NewLabel("id", qt.allureLabels.id))
	}

	if qt.allureLabels.suiteLabel != "" {
		t.Labels(allure.NewLabel("suite", qt.allureLabels.suiteLabel))
	}

	if qt.allureLabels.subSuite != "" {
		t.Labels(allure.NewLabel("subSuite", qt.allureLabels.subSuite))
	}

	if qt.allureLabels.parentSuite != "" {
		t.Labels(allure.NewLabel("parentSuite", qt.allureLabels.parentSuite))
	}

	if qt.allureLabels.story != "" {
		t.Story(qt.allureLabels.story)
	}

	if qt.allureLabels.tag != "" {
		t.Tags(qt.allureLabels.tag)
	}

	if qt.allureLabels.allureID != "" {
		t.ID(qt.allureLabels.allureID)
	}

	if qt.allureLabels.severity != nil {
		t.Severity(*qt.allureLabels.severity)
	}

	if qt.allureLabels.owner != "" {
		t.Owner(qt.allureLabels.owner)
	}

	if qt.allureLabels.lead != "" {
		t.Labels(allure.NewLabel("lead", qt.allureLabels.lead))
	}

	if qt.allureLabels.label != nil {
		t.Labels(*qt.allureLabels.label)
	}

	if len(qt.allureLabels.labels) != 0 {
		t.Labels(qt.allureLabels.labels...)
	}

	if qt.allureLabels.feature != "" {
		t.Feature(qt.allureLabels.feature)
	}

	if qt.allureLabels.epic != "" {
		t.Epic(qt.allureLabels.epic)
	}

	if len(qt.allureLabels.tags) != 0 {
		t.Tags(qt.allureLabels.tags...)
	}

	if qt.allureLabels.layer != "" {
		t.Labels(allure.NewLabel("layer", qt.allureLabels.layer))
	}
}

func (qt *cute) setInfoAllure(t T) {
	if qt.allureInfo.title != "" {
		t.Title(qt.allureInfo.title)
	}

	if qt.allureInfo.description != "" {
		t.Description(qt.allureInfo.description)
	}

	if qt.allureInfo.stage != "" {
		t.Labels(allure.NewLabel("stage", qt.allureInfo.stage))
	}
}
