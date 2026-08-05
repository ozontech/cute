package cute

import (
	"context"
	"testing"

	"github.com/ozontech/testo"
	"github.com/ozontech/testo/testoplugin"
	"github.com/ozontech/testo/testoreflect"

	allure "github.com/ozontech/testo-allure"
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
	severity    *allure.Severity
	owner       string
	lead        string
	label       *allure.Label
	labels      []allure.Label
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
	if t == nil {
		panic("could not start test without testing.T")
	}

	switch tt := t.(type) {
	case T:
		if qt.parallel {
			tt.Parallel()
		}

		return qt.executeTests(ctx, tt)
	case *testing.T:
		var (
			res  []ResultsHTTPBuilder
			opts []testoplugin.Option
		)

		if qt.isTableTest {
			// The wrapper test is only a launcher for table tests, each of which
			// produces its own allure result. Keep the wrapper out of the report.
			opts = append(opts, allure.WithExcluded(true))
		}

		testo.RunTest(tt, func(inT defaultT) {
			if qt.parallel {
				inT.Parallel()
			}

			res = qt.executeTests(ctx, inT)
		}, opts...)

		if res == nil {
			// executeTests died fatally (e.g. broken assert -> FailNow) before
			// returning. Pre-migration this stopped the caller's goroutine too,
			// so abort here instead of returning a nil result slice.
			tt.FailNow()
		}

		return res
	default:
		panic("t must be *testing.T or cute.T (a testo T with the allure plugin)")
	}
}

// runSeparateTest runs f as a new child test with its own allure result,
// mirroring allure-go's Run semantics (a separate test in the report, not a step).
func runSeparateTest(t T, name string, f func(t T)) {
	testo.Reflect(t).TestingT.Run(name, func(rt *testing.T) {
		testo.RunTest(rt, func(inT defaultT) {
			inT.Title(name)

			f(inT)
		})
	})
}

// executeTests is method for run tests
// It's could be table tests or usual tests
func (qt *cute) executeTests(ctx context.Context, t T) []ResultsHTTPBuilder {
	var (
		res = make([]ResultsHTTPBuilder, 0)
	)

	// Cycle for change number of Test
	for i := 0; i <= qt.countTests; i++ {
		currentTest := qt.tests[i]

		// Execute in a separate test for table tests
		if qt.isTableTest {
			runSeparateTest(t, currentTest.Name, func(inT T) {
				res = append(res, qt.executeInsideAllure(ctx, inT, currentTest))
			})
		} else {
			currentTest.Name = t.Name()

			// set labels, but only when running as a test, not inside a step
			if info, ok := testo.Reflect(t).Test.(testoreflect.RegularTestInfo); !ok || !info.IsSubtest {
				qt.setAllureInformation(t)
			}

			res = append(res, qt.executeInsideAllure(ctx, t, currentTest))
		}
	}

	return res
}

// executeInsideAllure is method for run test inside allure
// It's could be table tests or usual tests
func (qt *cute) executeInsideAllure(ctx context.Context, t T, currentTest *Test) ResultsHTTPBuilder {
	resT := currentTest.executeInsideAllure(ctx, t)

	// Remove from base struct all asserts
	currentTest.clearFields()

	return resT
}
