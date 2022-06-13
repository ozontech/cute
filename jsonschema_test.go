package cute

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/core/common"
	cuteErrors "github.com/ozontech/cute/errors"
	"github.com/stretchr/testify/require"
)

func TestValidateJSONSchemaEmptySchema(t *testing.T) {
	var (
		a        = HTTPTestMaker{}
		tBuilder = a.NewTestBuilder().(*test)
	)

	errs := tBuilder.validateJSONSchema(nil, []byte{})
	require.Len(t, errs, 0)
}

func TestValidateJSONSchemaFromString(t *testing.T) {
	var (
		a        = HTTPTestMaker{}
		tBuilder = a.NewTestBuilder().(*test)
		tempT    = common.NewT(t, "package", t.Name())
	)
	tempT.NewTest(t.Name(), "package")
	tempT.TestContext()

	body := []byte(`
	{
		"firstName": "Boris",
		"lastName": "Britva",
		"age": 77
	}
	`)

	tBuilder.expect.jsSchemaString = `
{
  "$id": "https://example.com/person.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Person",
  "type": "object",
  "properties": {
    "firstName": {
      "type": "string"
    },
    "lastName": {
      "type": "string"
    },
    "age": {
      "type": "integer",
      "minimum": 0
    }
  }
}
	`

	errs := tBuilder.validateJSONSchema(tempT, body)
	require.Len(t, errs, 0)
}

func TestValidateJSONSchemaFromStringWithError(t *testing.T) {
	var (
		a        = HTTPTestMaker{}
		tBuilder = a.NewTestBuilder().(*test)
		tempT    = common.NewT(t, "package", t.Name())
	)
	tempT.NewTest(t.Name(), "package")
	tempT.TestContext()

	body := []byte(`
	{
		"firstName": "Boris",
		"lastName": "Britva",
		"age": "1"
	}
	`)

	tBuilder.expect.jsSchemaString = `
	{
	  "$id": "https://example.com/person.schema.json",
	  "$schema": "https://json-schema.org/draft/2020-12/schema",
	  "title": "Person",
	  "type": "object",
	  "properties": {
	    "firstName": {
	      "type": "string"
	    },
	    "lastName": {
	      "type": "string"
	    },
	    "age": {
	      "type": "integer"
	    }
	  }
	}
	`

	errs := tBuilder.validateJSONSchema(tempT, body)
	require.Len(t, errs, 1)
	require.Error(t, errs[0])

	errWithName := errs[0].(cuteErrors.WithNameError)
	require.NotEmpty(t, errWithName.GetName())

	expectedError := errs[0].(cuteErrors.ExpectedError)
	require.Equal(t, "integer", expectedError.GetExpected())
	require.Equal(t, "string", expectedError.GetActual())
}

func TestValidateJSONSchemaFromByteWithTwoError(t *testing.T) {
	var (
		a        = HTTPTestMaker{}
		tBuilder = a.NewTestBuilder().(*test)
		tempT    = common.NewT(t, "package", t.Name())
	)
	tempT.NewTest(t.Name(), "package")
	tempT.TestContext()

	body := []byte(`
	{
		"firstName": "Boris",
		"lastName": "Britva",
		"age": "1"
	}
	`)

	tBuilder.expect.jsSchemaString = `
	{
	  "$id": "https://example.com/person.schema.json",
	  "$schema": "https://json-schema.org/draft/2020-12/schema",
	  "title": "Person",
	  "type": "object",
	  "properties": {
	    "firstName": {
	      "type": "string"
	    },
	    "lastName": {
	      "type": "integer"
	    },
	    "age": {
	      "type": "integer"
	    }
	  }
	}
	`

	errs := tBuilder.validateJSONSchema(tempT, body)
	require.Len(t, errs, 2)

	for _, err := range errs {
		errWithName := err.(cuteErrors.WithNameError)
		require.NotEmpty(t, errWithName.GetName())

		expectedError := err.(cuteErrors.ExpectedError)
		require.NotEmpty(t, expectedError.GetExpected())
		require.NotEmpty(t, expectedError.GetActual())
	}
}
