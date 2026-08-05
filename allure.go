package cute

import (
	allure "github.com/ozontech/testo-allure"
)

func (qt *cute) setAllureInformation(t T) {
	// Log main vars to allure
	qt.setLabelsAllure(t)
	qt.setInfoAllure(t)
	qt.setLinksAllure(t)
}

func (qt *cute) setLinksAllure(t T) {
	if qt.allureLinks.issue != "" {
		t.Links(allure.Issue(qt.allureLinks.issue))
	}

	if qt.allureLinks.testCase != "" {
		t.Links(allure.TMS(qt.allureLinks.testCase))
	}

	if qt.allureLinks.link != nil {
		t.Links(*qt.allureLinks.link)
	}

	if qt.allureLinks.tmsLink != "" {
		t.Links(allure.TMS(qt.allureLinks.tmsLink))
	}

	for _, tmsLink := range qt.allureLinks.tmsLinks {
		t.Links(allure.TMS(tmsLink))
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
