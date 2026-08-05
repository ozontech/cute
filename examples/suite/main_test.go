//go:build example
// +build example

package suite

import (
	"testing"

	"github.com/ozontech/testo"
	allure "github.com/ozontech/testo-allure"
)

func TestExampleSuite(t *testing.T) {
	testo.RunSuite(t, new(ExampleSuite), allure.WithOutputDir("../")) // custom, read Readme.md for more info
}
