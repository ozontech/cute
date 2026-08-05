package suite

import (
	"github.com/ozontech/testo"
	allure "github.com/ozontech/testo-allure"
)

// T is the canonical testo T for the suite examples.
// It satisfies cute.T thanks to the embedded *testo.T and *allure.PluginAllure.
type T struct {
	*testo.T
	*allure.PluginAllure
}
