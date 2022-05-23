package json

import (
	"encoding/json"
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/ozontech/cute"
)

// Contains is a function to assert that a jsonpath expression extracts a value in an array
// About expression - https://goessner.net/articles/JsonPath/
func Contains(expression string, expect interface{}) cute.AssertBody {
	return func(body []byte) error {
		return contains(body, expression, expect)
	}
}

// Equal is a function to assert that a jsonpath expression matches the given value
// About expression - https://goessner.net/articles/JsonPath/
func Equal(expression string, expect interface{}) cute.AssertBody {
	return func(body []byte) error {
		return equal(body, expression, expect)
	}
}

// NotEqual is a function to check json path expression value is not equal to given value
// About expression - https://goessner.net/articles/JsonPath/
func NotEqual(expression string, expect interface{}) cute.AssertBody {
	return func(body []byte) error {
		return notEqual(body, expression, expect)
	}
}

// Length is a function to asserts that value is the expected length
// About expression - https://goessner.net/articles/JsonPath/
func Length(expression string, expectLength int) cute.AssertBody {
	return func(body []byte) error {
		return length(body, expression, expectLength)
	}
}

// GreaterThan is a function to asserts that value is greater than the given length
// About expression - https://goessner.net/articles/JsonPath/
func GreaterThan(expression string, minimumLength int) cute.AssertBody {
	return func(body []byte) error {
		return greaterThan(body, expression, minimumLength)
	}
}

// LessThan is a function to asserts that value is less than the given length
// About expression - https://goessner.net/articles/JsonPath/
func LessThan(expression string, maximumLength int) cute.AssertBody {
	return func(body []byte) error {
		return lessThan(body, expression, maximumLength)
	}
}

// Present is a function to asserts that value is present
// About expression - https://goessner.net/articles/JsonPath/
func Present(expression string) cute.AssertBody {
	return func(body []byte) error {
		return present(body, expression)
	}
}

// NotPresent is a function to asserts that value is not present
// About expression - https://goessner.net/articles/JsonPath/
func NotPresent(expression string) cute.AssertBody {
	return func(body []byte) error {
		return notPresent(body, expression)
	}
}

// GetValueFromJSON is function for get value from json
// TODO create tests
func GetValueFromJSON(js []byte, expression string) (interface{}, error) {
	v := interface{}(nil)

	err := json.Unmarshal(js, &v)
	if err != nil {
		return nil, err
	}

	value, err := jsonpath.Get(expression, v)
	if err != nil {
		return nil, fmt.Errorf("evaluating '%s' resulted in error: '%s'", expression, err)
	}

	return value, nil
}
