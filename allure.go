package cute

func (qt *cute) setAllureInformation(t allureProvider) {
	// Log main vars to allureProvider
	qt.setLabelsAllure(t)
	qt.setInfoAllure(t)
	qt.setLinksAllure(t)
}

func (qt *cute) setLinksAllure(t linksAllureProvider) {
	if qt.allureLinks.issue != "" {
		t.SetIssue(qt.allureLinks.issue)
	}

	if qt.allureLinks.testCase != "" {
		t.SetTestCase(qt.allureLinks.testCase)
	}

	if qt.allureLinks.link != nil {
		t.Link(qt.allureLinks.link)
	}

	if qt.allureLinks.tmsLink != "" {
		t.TmsLink(qt.allureLinks.tmsLink)
	}

	if len(qt.allureLinks.tmsLinks) > 0 {
		t.TmsLinks(qt.allureLinks.tmsLinks...)
	}
}

func (qt *cute) setLabelsAllure(t labelsAllureProvider) {
	if qt.allureLabels.id != "" {
		t.ID(qt.allureLabels.id)
	}

	if qt.allureLabels.suiteLabel != "" {
		t.AddSuiteLabel(qt.allureLabels.suiteLabel)
	}

	if qt.allureLabels.subSuite != "" {
		t.AddSubSuite(qt.allureLabels.subSuite)
	}

	if qt.allureLabels.parentSuite != "" {
		t.AddParentSuite(qt.allureLabels.parentSuite)
	}

	if qt.allureLabels.story != "" {
		t.Story(qt.allureLabels.story)
	}

	if qt.allureLabels.tag != "" {
		t.Tag(qt.allureLabels.tag)
	}

	if qt.allureLabels.allureID != "" {
		t.AllureID(qt.allureLabels.allureID)
	}

	if qt.allureLabels.severity != "" {
		t.Severity(qt.allureLabels.severity)
	}

	if qt.allureLabels.owner != "" {
		t.Owner(qt.allureLabels.owner)
	}

	if qt.allureLabels.lead != "" {
		t.Lead(qt.allureLabels.lead)
	}

	if qt.allureLabels.label != nil {
		t.Label(qt.allureLabels.label)
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
		t.Layer(qt.allureLabels.layer)
	}
}

func (qt *cute) setInfoAllure(t infoAllureProvider) {
	if qt.allureInfo.title != "" {
		t.Title(qt.allureInfo.title)
	}

	if qt.allureInfo.description != "" {
		t.Description(qt.allureInfo.description)
	}

	if qt.allureInfo.stage != "" {
		t.Stage(qt.allureInfo.stage)
	}
}
