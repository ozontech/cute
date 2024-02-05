package cute

import (
	"fmt"

	"github.com/ozontech/cute/errors"
	"github.com/xeipuuv/gojsonschema"
)

// Validate is a function to validate json by json schema.
// Automatically add information about validation to allure.
func (it *Test) validateJSONSchema(t internalT, body []byte) []error {
	var (
		expect gojsonschema.JSONLoader
	)

	switch {
	case it.Expect.JSONSchema.String != "":
		expect = gojsonschema.NewStringLoader(it.Expect.JSONSchema.String)
	case it.Expect.JSONSchema.Byte != nil:
		expect = gojsonschema.NewBytesLoader(it.Expect.JSONSchema.Byte)
	case it.Expect.JSONSchema.File != "":
		expect = gojsonschema.NewReferenceLoader(it.Expect.JSONSchema.File)
	default:
		return nil
	}

	return executeWithStep(t, "Validate body by JSON schema", func(t T) []error {
		return checkJSONSchema(expect, body)
	})
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
