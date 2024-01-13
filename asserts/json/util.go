package json

import (
	"bytes"
	"fmt"
	"github.com/ohler55/ojg/oj"
	"reflect"
	"strings"

	"github.com/ozontech/cute/errors"
)

// Contains is a function to assert that a jsonpath expression extracts a value in an array
// Given the response is {"first": 777, "second": [{"key_1": "some_key", "value": "some_value"}]}, we can assert on the result like so `$.second[? @.key_1=="some_key"].value`, "some_value"
// About expression - https://goessner.net/articles/JsonPath/
func contains(data []byte, expression string, expect interface{}) error {
	values, err := GetValueFromJSON(data, expression)
	if err != nil {
		return err
	}

	for _, value := range values {
		ok, found := insideArray(value, expect)
		if !ok {
			return errors.NewAssertError("Contains", fmt.Sprintf("on path %v. %v could not be applied builtin len()", expression, expect), nil, nil)
		}

		if !found {
			return errors.NewAssertError("Contains", fmt.Sprintf("on path %v. expect %v, but actual %v", expression, expect, value), value, expect)
		}
	}

	return nil
}

// Equal is a function to assert that a jsonpath expression matches the given value
// About expression - https://goessner.net/articles/JsonPath/
func equal(data []byte, expression string, expect interface{}) error {
	values, err := GetValueFromJSON(data, expression)
	if err != nil {
		return err
	}

	for _, value := range values {
		if !objectsAreEqual(value, expect) {
			return errors.NewAssertError("Equal", fmt.Sprintf("on path %v. expect %v, but actual %v", expression, expect, value), value, expect)
		}
	}

	return nil
}

// NotEqual is a function to check json path expression value is not equal to given value
// About expression - https://goessner.net/articles/JsonPath/
func notEqual(data []byte, expression string, expect interface{}) error {
	values, err := GetValueFromJSON(data, expression)
	if err != nil {
		return err
	}

	for _, value := range values {
		if objectsAreEqual(value, expect) {
			return errors.NewAssertError("NotEqual", fmt.Sprintf("on path %v. expect %v, but actual %v", expression, expect, value), value, expect)
		}
	}

	return nil
}

// EqualJSON is a function to check json path expression value is equal to given json
// About expression - https://goessner.net/articles/JsonPath/
func equalJSON(data []byte, expression string, expect []byte) error {
	values, err := GetValueFromJSON(data, expression)
	if err != nil {
		return err
	}

	obj, err := oj.Parse(expect)
	if err != nil {
		return fmt.Errorf("could not parse json in EqualJSON error: '%s'", err)
	}

	for _, value := range values {
		if !objectsAreEqual(value, obj) {
			return errors.NewAssertError(
				"EqualJSON",
				fmt.Sprintf("on path %v. expect %v (from json %v), but actual %v", expression, obj, expect, value),
				value,
				obj)
		}
	}

	return nil
}

// NotEqualJSON is a function to check json path expression value is not equal to given json
// About expression - https://goessner.net/articles/JsonPath/
func notEqualJSON(data []byte, expression string, expect []byte) error {
	values, err := GetValueFromJSON(data, expression)
	if err != nil {
		return err
	}

	obj, err := oj.Parse(expect)
	if err != nil {
		return fmt.Errorf("could not parse json in NotEqualJSON error: '%s'", err)
	}

	for _, value := range values {
		if objectsAreEqual(value, obj) {
			return errors.NewAssertError(
				"NotEqualJSON",
				fmt.Sprintf("on path %v. expect %v (from json %v), but actual %v", expression, obj, expect, value),
				value,
				obj)
		}
	}

	return nil
}

// Length is a function to asserts that value is the expected length
// About expression - https://goessner.net/articles/JsonPath/
func length(data []byte, expression string, expectLength int) error {
	values, err := GetValueFromJSON(data, expression)
	if err != nil {
		return err
	}

	for _, value := range values {
		v := reflect.ValueOf(value)
		if v.Len() != expectLength {
			return errors.NewAssertError("Length", fmt.Sprintf("on path %v. expect lenght %v, but actual %v", expression, expectLength, v.Len()), v.Len(), expectLength)
		}
	}

	return nil
}

// GreaterThan is a function to asserts that value is greater than the given length
// About expression - https://goessner.net/articles/JsonPath/
func greaterThan(data []byte, expression string, minimumLength int) error {
	values, err := GetValueFromJSON(data, expression)
	if err != nil {
		return err
	}

	for _, value := range values {
		v := reflect.ValueOf(value)
		if v.Len() <= minimumLength {
			return errors.NewAssertError("GreaterThan", fmt.Sprintf("on path %v. %v is greater than %v", expression, v.Len(), minimumLength), v.Len(), minimumLength)
		}
	}

	return nil
}

// GreaterOrEqualThan is a function to asserts that value is greater or equal than the given length
// About expression - https://goessner.net/articles/JsonPath/
func greaterOrEqualThan(data []byte, expression string, minimumLength int) error {
	values, err := GetValueFromJSON(data, expression)
	if err != nil {
		return err
	}

	for _, value := range values {
		v := reflect.ValueOf(value)
		if v.Len() < minimumLength {
			return errors.NewAssertError("GreaterOrEqualThan", fmt.Sprintf("on path %v. %v is greater or equal than %v", expression, v.Len(), minimumLength), v.Len(), minimumLength)
		}
	}

	return nil
}

// LessThan is a function to asserts that value is less than the given length
// About expression - https://goessner.net/articles/JsonPath/
func lessThan(data []byte, expression string, maximumLength int) error {
	values, err := GetValueFromJSON(data, expression)
	if err != nil {
		return err
	}

	for _, value := range values {
		v := reflect.ValueOf(value)
		if v.Len() >= maximumLength {
			return errors.NewAssertError("LessThan", fmt.Sprintf("on path %v. %v is less than %v", expression, v.Len(), maximumLength), v.Len(), maximumLength)
		}
	}

	return nil
}

// LessOrEqualThan is a function to asserts that value is less or equal than the given length
// About expression - https://goessner.net/articles/JsonPath/
func lessOrEqualThan(data []byte, expression string, maximumLength int) error {
	values, err := GetValueFromJSON(data, expression)
	if err != nil {
		return err
	}

	for _, value := range values {
		v := reflect.ValueOf(value)
		if v.Len() > maximumLength {
			return errors.NewAssertError("LessThan", fmt.Sprintf("on path %v. %v is less or equal than %v", expression, v.Len(), maximumLength), v.Len(), maximumLength)
		}
	}

	return nil
}

// notEmpty is a function to asserts that value is not empty (!= 0, != null)
// About expression - https://goessner.net/articles/JsonPath/
func notEmpty(data []byte, expression string) error {
	values, _ := GetValueFromJSON(data, expression)

	if len(values) == 0 {
		return errors.NewAssertError("NotEmpty", fmt.Sprintf("on path %v. value is not present", expression), nil, nil)
	}

	for _, value := range values {
		if isEmpty(value) {
			return errors.NewAssertError("NotEmpty", fmt.Sprintf("on path %v. value is not present", expression), nil, nil)
		}
	}

	return nil
}

// Present is a function to asserts that value is present
// value can be 0 or null
// About expression - https://goessner.net/articles/JsonPath/
func present(data []byte, expression string) error {
	values, err := GetValueFromJSON(data, expression)
	if err != nil || len(values) == 0 {
		return errors.NewAssertError("Present", fmt.Sprintf("on path %v. value not present", expression), nil, nil)
	}

	return nil
}

// NotPresent is a function to asserts that value is not present
// About expression - https://goessner.net/articles/JsonPath/
func notPresent(data []byte, expression string) error {
	values, _ := GetValueFromJSON(data, expression)

	for _, value := range values {
		if !isEmpty(value) {
			return errors.NewAssertError("NotPresent", fmt.Sprintf("on path %v. value present", expression), nil, nil)
		}
	}

	return nil
}

func objectsAreEqual(expect, actual interface{}) bool {
	if reflect.DeepEqual(expect, actual) {
		return true
	}

	if expect == nil || actual == nil {
		return expect == actual
	}

	if fmt.Sprintf("%v", expect) == fmt.Sprintf("%v", actual) {
		return true
	}

	exp, ok := expect.([]byte)
	if !ok {
		return reflect.DeepEqual(expect, actual)
	}

	act, ok := actual.([]byte)
	if !ok {
		return false
	}

	if exp == nil || act == nil {
		return exp == nil && act == nil
	}

	return bytes.Equal(exp, act)
}

func isEmpty(object interface{}) bool {
	if object == nil {
		return true
	}

	objValue := reflect.ValueOf(object)

	switch objValue.Kind() {
	case reflect.Ptr:
		if objValue.IsNil() {
			return true
		}
		deref := objValue.Elem().Interface()
		return isEmpty(deref)
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return objValue.Len() == 0
	default:
		zero := reflect.Zero(objValue.Type())
		return reflect.DeepEqual(object, zero.Interface())
	}
}

func insideArray(list interface{}, element interface{}) (ok, found bool) {
	var (
		listValue    = reflect.ValueOf(list)
		elementValue = reflect.ValueOf(element)
	)

	defer func() {
		if err := recover(); err != nil {
			ok = false
			found = false
		}
	}()

	if reflect.TypeOf(list).Kind() == reflect.String {
		return true, strings.Contains(listValue.String(), elementValue.String())
	}

	if reflect.TypeOf(list).Kind() == reflect.Map {
		mapKeys := listValue.MapKeys()

		for i := 0; i < len(mapKeys); i++ {
			if objectsAreEqual(mapKeys[i].Interface(), element) {
				return true, true
			}
		}

		return true, false
	}

	for i := 0; i < listValue.Len(); i++ {
		if objectsAreEqual(listValue.Index(i).Interface(), element) {
			return true, true
		}
	}

	return true, false
}
