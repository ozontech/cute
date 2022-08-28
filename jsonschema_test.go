package cute

import (
	"testing"

	cuteErrors "github.com/ozontech/cute/errors"
	"github.com/stretchr/testify/require"
)

func TestValidateJSONSchemaEmptySchema(t *testing.T) {
	var (
		tBuilder = createDefaultTest(nil)
	)

	errs := tBuilder.validateJSONSchema(nil, []byte{})
	require.Len(t, errs, 0)
}

func TestValidateJSONSchemaFromString(t *testing.T) {
	var (
		tBuilder = createDefaultTest(nil)
		tempT    = createAllureT(t)
	)

	body := []byte(`
	{
		"firstName": "Boris",
		"lastName": "Britva",
		"age": 77
	}
	`)

	tBuilder.Expect.JSONSchema.String = `
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
		tBuilder = createDefaultTest(nil)
		tempT    = createAllureT(t)
	)

	body := []byte(`
	{
		"firstName": "Boris",
		"lastName": "Britva",
		"age": "1"
	}
	`)

	tBuilder.Expect.JSONSchema.String = `
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

	expectedError := errs[0].(cuteErrors.WithFields)
	require.Equal(t, "integer", expectedError.GetFields()["Expected"])
	require.Equal(t, "string", expectedError.GetFields()["Actual"])
	require.Equal(t, "age", expectedError.GetFields()["Field"])
	require.Equal(t, "(root).age", expectedError.GetFields()["Path"])
}

func TestValidateJSONSchemaFromByteWithTwoError(t *testing.T) {
	var (
		tBuilder = createDefaultTest(nil)
		tempT    = createAllureT(t)
	)

	body := []byte(`
	{
		"firstName": "Boris",
		"lastName": "Britva",
		"age": "1"
	}
	`)

	tBuilder.Expect.JSONSchema.String = `
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

		expectedError := err.(cuteErrors.WithFields)
		require.NotEmpty(t, expectedError.GetFields()["Actual"])
		require.NotEmpty(t, expectedError.GetFields()["Expected"])
	}
}
