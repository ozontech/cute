//go:build example
// +build example

package suite

import (
	"os"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestExampleSuite(t *testing.T) {
	os.Setenv("ALLURE_OUTPUT_PATH", "../") // custom, read Readme.md for more info
	suite.RunSuite(t, new(ExampleSuite))
}
