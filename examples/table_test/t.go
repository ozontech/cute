//go:build example
// +build example

package table_test

import (
	"testing"

	"github.com/ozontech/testo"
	allure "github.com/ozontech/testo-allure"
)

// T is the canonical testo T for the table tests examples.
// It satisfies cute.T thanks to the embedded *testo.T and *allure.PluginAllure.
type T struct {
	*testo.T
	*allure.PluginAllure
}

// run executes f as a testo test with a custom allure output dir,
// read Readme.md for more info.
func run(t *testing.T, f func(t T)) {
	testo.RunTest(t, f, allure.WithOutputDir("../"))
}
