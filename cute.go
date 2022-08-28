package cute

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/core/allure_manager/manager"
	"github.com/ozontech/allure-go/pkg/framework/core/common"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type cute struct {
	httpClient *http.Client

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
	label       allure.Label
	labels      []allure.Label
}

type allureLinks struct {
	issue    string
	testCase string
	link     allure.Link
}

func (it *cute) ExecuteTest(ctx context.Context, t testing.TB) []ResultsHTTPBuilder {
	var internalT allureProvider

	tOriginal, ok := t.(*testing.T)
	if ok {
		newT := createAllureT(tOriginal)
		defer newT.FinishTest()

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

	// set labels
	it.setAllureInformation(allureProvider)

	// Cycle for change number of Test
	for i := 0; i <= it.countTests; i++ {
		currentTest := it.tests[i]

		// Execute by new T for table tests
		if it.isTableTest {
			tableTestName := currentTest.Name

			allureProvider.Run(tableTestName, func(inT provider.T) {
				inT.Logf("Test start %v", tableTestName)
				resT := currentTest.execute(ctx, inT)
				res = append(res, resT)
				inT.Logf("Test finished %v", tableTestName)
			})
		} else {
			currentTest.Name = allureProvider.Name()

			allureProvider.Logf("Test start %v", currentTest.Name)
			resT := currentTest.execute(ctx, allureProvider)
			res = append(res, resT)
			allureProvider.Logf("Test finished %v", currentTest.Name)
		}
	}

	return res
}
