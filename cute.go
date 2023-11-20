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

func (it *cute) ExecuteTest(ctx context.Context, t tProvider) []ResultsHTTPBuilder {
	var internalT allureProvider

	if t == nil {
		panic("could not start test without testing.T")
	}

	stepCtx, isStepCtx := t.(provider.StepCtx)
	if isStepCtx {
		return it.executeTestInsideStep(ctx, stepCtx)
	}

	tOriginal, ok := t.(*testing.T)
	if ok {
		newT := createAllureT(tOriginal)
		if !it.isTableTest {
			defer newT.FinishTest() //nolint
		}

		internalT = newT
	}

	allureT, ok := t.(provider.T)
	if ok {
		internalT = allureT
	}

	if it.parallel {
		internalT.Parallel()
	}

	return it.executeTest(ctx, internalT)
}

func (it *cute) executeTestInsideStep(ctx context.Context, stepCtx provider.StepCtx) []ResultsHTTPBuilder {
	var (
		res = make([]ResultsHTTPBuilder, 0)
	)

	// Cycle for change number of Test
	for i := 0; i <= it.countTests; i++ {
		currentTest := it.tests[i]

		result := currentTest.executeInsideStep(ctx, stepCtx)
		res = append(res, result)
	}

	return res
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

func (it *cute) executeTest(ctx context.Context, allureProvider allureProvider) []ResultsHTTPBuilder {
	var (
		res = make([]ResultsHTTPBuilder, 0)
	)

	// Cycle for change number of Test
	for i := 0; i <= it.countTests; i++ {
		currentTest := it.tests[i]

		// Execute by new T for table tests
		if it.isTableTest {
			tableTestName := currentTest.Name

			allureProvider.Run(tableTestName, func(inT provider.T) {
				inT.Title(tableTestName) // set current test name

				resultTest := it.executeSingleTest(ctx, inT, currentTest)
				res = append(res, resultTest)
			})
		} else {
			currentTest.Name = allureProvider.Name()

			// set labels
			it.setAllureInformation(allureProvider)

			resultTest := it.executeSingleTest(ctx, allureProvider, currentTest)
			if resultTest != nil {
				res = append(res, resultTest)
			}
		}
	}

	return res
}

func (it *cute) executeSingleTest(ctx context.Context, allureProvider allureProvider, currentTest *Test) ResultsHTTPBuilder {
	allureProvider.Logf("Test start %v", currentTest.Name)

	resT := currentTest.execute(ctx, allureProvider)
	if resT.IsFailed() {
		allureProvider.Fail()
		allureProvider.Logf("Test was failed %v", currentTest.Name)

		return resT
	}

	allureProvider.Logf("Test finished %v", currentTest.Name)

	return resT
}
