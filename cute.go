package cute

import (
	"context"
	"strings"
	"testing"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/core/allure_manager/manager"
	"github.com/ozontech/allure-go/pkg/framework/core/common"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type cute struct {
	baseProps *HTTPTestMaker

	parallel bool

	allureInfo   *allureInformation
	allureLinks  *allureLinks
	allureLabels *allureLabels

	countTests int // Общее количество тестов.

	isTableTest bool
	tests       []*Test
}

type allureInformation struct {
	title       string
	description string
	stage       string
}

type allureLabels struct {
	id          string
	feature     string
	epic        string
	tag         string
	tags        []string
	suiteLabel  string
	subSuite    string
	parentSuite string
	story       string
	severity    allure.SeverityType
	owner       string
	lead        string
	label       *allure.Label
	labels      []*allure.Label
	allureID    string
	layer       string
}

type allureLinks struct {
	issue    string
	testCase string
	link     *allure.Link
	tmsLink  string
	tmsLinks []string
}

func (qt *cute) ExecuteTest(ctx context.Context, t tProvider) []ResultsHTTPBuilder {
	var internalT allureProvider

	if t == nil {
		panic("could not start test without testing.T")
	}

	stepCtx, isStepCtx := t.(provider.StepCtx)
	if isStepCtx {
		return qt.executeTestsInsideStep(ctx, stepCtx)
	}

	tOriginal, ok := t.(*testing.T)
	if ok {
		newT := createAllureT(tOriginal)
		if !qt.isTableTest {
			defer newT.FinishTest() //nolint
		}

		internalT = newT
	}

	allureT, ok := t.(provider.T)
	if ok {
		internalT = allureT
	}

	if qt.parallel {
		internalT.Parallel()
	}

	return qt.executeTests(ctx, internalT)
}

func createAllureT(t *testing.T) *common.Common {
	var (
		newT        = common.NewT(t)
		callers     = strings.Split(t.Name(), "/")
		providerCfg = manager.NewProviderConfig().
				WithFullName(t.Name()).
				WithPackageName("package").
				WithSuiteName(t.Name()).
				WithRunner(callers[0])
		newProvider = manager.NewProvider(providerCfg)
	)

	newProvider.NewTest(t.Name(), "package")

	newT.SetProvider(newProvider)
	newT.Provider.TestContext()

	return newT
}

// executeTestsInsideStep is method for run group of tests inside provider.StepCtx
func (qt *cute) executeTestsInsideStep(ctx context.Context, stepCtx provider.StepCtx) []ResultsHTTPBuilder {
	var (
		res = make([]ResultsHTTPBuilder, 0)
	)

	// Cycle for change number of Test
	for i := 0; i <= qt.countTests; i++ {
		currentTest := qt.tests[i]

		result := currentTest.executeInsideStep(ctx, stepCtx)

		// Remove from base struct all asserts
		currentTest.clearFields()

		res = append(res, result)
	}

	return res
}

// executeTests is method for run tests
// It's could be table tests or usual tests
func (qt *cute) executeTests(ctx context.Context, allureProvider allureProvider) []ResultsHTTPBuilder {
	var (
		res = make([]ResultsHTTPBuilder, 0)
	)

	// Cycle for change number of Test
	for i := 0; i <= qt.countTests; i++ {
		currentTest := qt.tests[i]

		// Execute by new T for table tests
		if qt.isTableTest {
			tableTestName := currentTest.Name

			allureProvider.Run(tableTestName, func(inT provider.T) {
				// Set current test name
				inT.Title(tableTestName)

				res = append(res, qt.executeSingleTest(ctx, inT, currentTest))
			})
		} else {
			currentTest.Name = allureProvider.Name()

			// set labels
			qt.setAllureInformation(allureProvider)

			res = append(res, qt.executeSingleTest(ctx, allureProvider, currentTest))
		}
	}

	return res
}

func (qt *cute) executeSingleTest(ctx context.Context, allureProvider allureProvider, currentTest *Test) ResultsHTTPBuilder {
	allureProvider.Logf("Test start %v", currentTest.Name)

	resT := currentTest.executeInsideAllure(ctx, allureProvider)

	switch resT.GetResultState() {
	case ResultStateBroken:
		allureProvider.BrokenNow()
		allureProvider.Logf("Test broken %v", currentTest.Name)
	case ResultStateFail:
		allureProvider.Fail()
		allureProvider.Logf("Test failed %v", currentTest.Name)
	case resultStateFailNow:
		allureProvider.FailNow()
		allureProvider.Logf("Test failed %v", currentTest.Name)
	case ResultStateSuccess:
		allureProvider.Logf("Test finished %v", currentTest.Name)
	}

	// Remove from base struct all asserts
	currentTest.clearFields()

	return resT
}
