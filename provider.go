package cute

import (
	"github.com/ozontech/testo"

	allure "github.com/ozontech/testo-allure"
)

// T is the test handle passed to user-defined asserts and middleware.
//
// It is satisfied by any testo T struct that embeds *testo.T and
// *allure.PluginAllure:
//
//	type T struct {
//		*testo.T
//		*allure.PluginAllure
//	}
//
// Inside asserts you may create allure steps with allure.Step(t, ...),
// add attachments with t.Attach(...), parameters with t.Parameters(...),
// log information, etc.
type T interface {
	testo.CommonT
	allure.Interface
}

// internalT is an internal alias kept for readability of the execution flow.
type internalT = T

// defaultT is the T used when a raw *testing.T is passed to ExecuteTest.
type defaultT struct {
	*testo.T
	*allure.PluginAllure
}

// tProvider is a minimal testing handle accepted by ExecuteTest and Execute.
// It is satisfied by *testing.T and by any cute.T (testo T with allure plugin).
type tProvider interface {
	Fail()
	FailNow()

	Name() string

	Log(args ...interface{})
	Logf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}
