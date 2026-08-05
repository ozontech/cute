package cute

import (
	"testing"

	cuteErrors "github.com/ozontech/cute/errors"
	"github.com/stretchr/testify/require"
)

func TestValidateJSONSchemaEmptySchema(t *testing.T) {
	var (
		tBuilder = createDefaultTest(&HTTPTestMaker{middleware: new(Middleware)})
	)

	tBuilder.initEmptyFields()

	errs := tBuilder.validateJSONSchema(nil, []byte{})
	require.Len(t, errs, 0)
}

func TestValidateJSONSchemaFromString(t *testing.T) {
	var (
		tBuilder = createDefaultTest(&HTTPTestMaker{middleware: new(Middleware)})
	)

	tBuilder.initEmptyFields()

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

	runCuteTest(t, func(ct defaultT) {
		errs := tBuilder.validateJSONSchema(ct, body)
		require.Len(ct, errs, 0)
	})
}

func TestValidateJSONSchemaFromStringWithError(t *testing.T) {
	var (
		tBuilder = createDefaultTest(&HTTPTestMaker{middleware: new(Middleware)})
	)

	tBuilder.initEmptyFields()

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

	runCuteTest(t, func(ct defaultT) {
		errs := tBuilder.validateJSONSchema(ct, body)
		require.Len(ct, errs, 1)
		require.Error(ct, errs[0])

		errWithName := errs[0].(cuteErrors.WithNameError)
		require.NotEmpty(ct, errWithName.GetName())

		expectedError := errs[0].(cuteErrors.WithFields)
		require.Equal(ct, "integer", expectedError.GetFields()["Expected"])
		require.Equal(ct, "string", expectedError.GetFields()["Actual"])
		require.Equal(ct, "age", expectedError.GetFields()["Field"])
		require.Equal(ct, "(root).age", expectedError.GetFields()["Path"])
	})
}

func TestValidateJSONSchemaFromByteWithTwoError(t *testing.T) {
	var (
		tBuilder = createDefaultTest(&HTTPTestMaker{middleware: new(Middleware)})
	)

	tBuilder.initEmptyFields()

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

	runCuteTest(t, func(ct defaultT) {
		errs := tBuilder.validateJSONSchema(ct, body)
		require.Len(ct, errs, 2)

		for _, err := range errs {
			errWithName := err.(cuteErrors.WithNameError)
			require.NotEmpty(ct, errWithName.GetName())

			expectedError := err.(cuteErrors.WithFields)
			require.NotEmpty(ct, expectedError.GetFields()["Actual"])
			require.NotEmpty(ct, expectedError.GetFields()["Expected"])
		}
	})
}
