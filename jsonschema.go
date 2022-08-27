package cute

import (
	"fmt"

	"github.com/ozontech/cute/errors"
	"github.com/xeipuuv/gojsonschema"
)

// Validate is a function to validate json by json schema.
// Automatically add information about validation to allure.
func (it *cute) validateJSONSchema(t internalT, body []byte) []error {
	var (
		scope  = make([]error, 0)
		expect gojsonschema.JSONLoader
	)

	switch {
	case it.tests[it.correctTest].expect.JSONSchemaString != "":
		expect = gojsonschema.NewStringLoader(it.tests[it.correctTest].expect.JSONSchemaString)
	case it.tests[it.correctTest].expect.JSONSchemaByte != nil:
		expect = gojsonschema.NewBytesLoader(it.tests[it.correctTest].expect.JSONSchemaByte)
	case it.tests[it.correctTest].expect.JSONSchemaFile != "":
		expect = gojsonschema.NewReferenceLoader(it.tests[it.correctTest].expect.JSONSchemaFile)
	default:
		return nil
	}

	executeWithStep(t, "Validate body by JSON schema", func(t T) []error {
		scope = checkJSONSchema(expect, body)

		return scope
	}, false)

	return scope
}

func checkJSONSchema(expect gojsonschema.JSONLoader, data []byte) []error {
	scope := make([]error, 0)
	validateResult, err := gojsonschema.Validate(expect, gojsonschema.NewBytesLoader(data))
	if err != nil {
		return []error{err}
	}

	if !validateResult.Valid() && len(validateResult.Errors()) > 0 {
		for _, resultError := range validateResult.Errors() {
			scope = append(
				scope,
				createJSONSchemaError(resultError),
			)
		}
	}

	return scope
}

func createJSONSchemaError(err gojsonschema.ResultError) error {
	fields := make(map[string]interface{})
	textError := ""

	if v, ok := err.Details()["context"]; ok {
		textError = fmt.Sprintf("On path: %v.", v)
		fields["Path"] = v
	}

	if v, ok := err.Details()["field"]; ok {
		textError = fmt.Sprintf("%v Error field: %v.", textError, v)
		fields["Field"] = v
	}

	textError = fmt.Sprintf("%v Error: %v.", textError, err.String())

	assertError := errors.NewAssertError(
		fmt.Sprintf("Error \"%v\"", err.Type()),
		textError,
		err.Details()["given"],
		err.Details()["expected"])

	assertError.(errors.WithFields).PutFields(fields)

	return assertError
}
