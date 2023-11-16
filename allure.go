package cute

func (it *cute) setAllureInformation(t allureProvider) {
	// Log main vars to allureProvider
	it.setLabelsAllure(t)
	it.setInfoAllure(t)
	it.setLinksAllure(t)
}

func (it *cute) setLinksAllure(t linksAllureProvider) {
	if it.allureLinks.issue != "" {
		t.SetIssue(it.allureLinks.issue)
	}
	if it.allureLinks.testCase != "" {
		t.SetTestCase(it.allureLinks.testCase)
	}
	if it.allureLinks.link != nil {
		t.Link(it.allureLinks.link)
	}
	if it.allureLinks.tmsLink != "" {
		t.TmsLink(it.allureLinks.tmsLink)
	}
	if len(it.allureLinks.tmsLinks) > 0 {
		t.TmsLinks(it.allureLinks.tmsLinks...)
	}
}

func (it *cute) setLabelsAllure(t labelsAllureProvider) {
	if it.allureLabels.id != "" {
		t.ID(it.allureLabels.id)
	}
	if it.allureLabels.suiteLabel != "" {
		t.AddSuiteLabel(it.allureLabels.suiteLabel)
	}
	if it.allureLabels.subSuite != "" {
		t.AddSubSuite(it.allureLabels.subSuite)
	}
	if it.allureLabels.parentSuite != "" {
		t.AddParentSuite(it.allureLabels.parentSuite)
	}
	if it.allureLabels.story != "" {
		t.Story(it.allureLabels.story)
	}
	if it.allureLabels.tag != "" {
		t.Tag(it.allureLabels.tag)
	}
	if it.allureLabels.allureID != "" {
		t.AllureID(it.allureLabels.allureID)
	}
	if it.allureLabels.severity != "" {
		t.Severity(it.allureLabels.severity)
	}
	if it.allureLabels.owner != "" {
		t.Owner(it.allureLabels.owner)
	}
	if it.allureLabels.lead != "" {
		t.Lead(it.allureLabels.lead)
	}
	if it.allureLabels.label != nil {
		t.Label(it.allureLabels.label)
	}
	if len(it.allureLabels.labels) != 0 {
		t.Labels(it.allureLabels.labels...)
	}
	if it.allureLabels.feature != "" {
		t.Feature(it.allureLabels.feature)
	}
	if it.allureLabels.epic != "" {
		t.Epic(it.allureLabels.epic)
	}
	if len(it.allureLabels.tags) != 0 {
		t.Tags(it.allureLabels.tags...)
	}
	if it.allureLabels.layer != "" {
		it.Layer(it.allureLabels.layer)
	}
}

func (it *cute) setInfoAllure(t infoAllureProvider) {
	if it.allureInfo.title != "" {
		t.Title(it.allureInfo.title)
	}
	if it.allureInfo.description != "" {
		t.Description(it.allureInfo.description)
	}
	if it.allureInfo.stage != "" {
		t.Stage(it.allureInfo.stage)
	}
}
