//go:build example
// +build example

package examples

import (
	"github.com/ozontech/testo"
	allure "github.com/ozontech/testo-allure"
)

// T is the canonical testo T for the examples package.
// It satisfies cute.T thanks to the embedded *testo.T and *allure.PluginAllure.
type T struct {
	*testo.T
	*allure.PluginAllure
}
