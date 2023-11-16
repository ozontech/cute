package json

import (
	"fmt"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
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

// LengthGreaterThan is a function to asserts that value is greater than the given length
// About expression - https://goessner.net/articles/JsonPath/
func LengthGreaterThan(expression string, minimumLength int) cute.AssertBody {
	return func(body []byte) error {
		return greaterThan(body, expression, minimumLength)
	}
}

// LengthGreaterOrEqualThan is a function to asserts that value is greater or equal than the given length
// About expression - https://goessner.net/articles/JsonPath/
func LengthGreaterOrEqualThan(expression string, minimumLength int) cute.AssertBody {
	return func(body []byte) error {
		return greaterOrEqualThan(body, expression, minimumLength)
	}
}

// LengthLessThan is a function to asserts that value is less than the given length
// About expression - https://goessner.net/articles/JsonPath/
func LengthLessThan(expression string, maximumLength int) cute.AssertBody {
	return func(body []byte) error {
		return lessThan(body, expression, maximumLength)
	}
}

// LengthLessOrEqualThan is a function to asserts that value is less or equal than the given length
// About expression - https://goessner.net/articles/JsonPath/
func LengthLessOrEqualThan(expression string, maximumLength int) cute.AssertBody {
	return func(body []byte) error {
		return lessOrEqualThan(body, expression, maximumLength)
	}
}

// Present is a function to asserts that value is present
// value can be nil or 0
// About expression - https://goessner.net/articles/JsonPath/
func Present(expression string) cute.AssertBody {
	return func(body []byte) error {
		return present(body, expression)
	}
}

// NotEmpty is a function to asserts that value is present
// value can't be nil or 0
// About expression - https://goessner.net/articles/JsonPath/
func NotEmpty(expression string) cute.AssertBody {
	return func(body []byte) error {
		return notEmpty(body, expression)
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
	obj, err := oj.Parse(js)
	if err != nil {
		return nil, fmt.Errorf("could not parse json in ,GetValueFromJSON error: '%s'", err)
	}

	jsonPath, err := jp.ParseString(expression)
	if err != nil {
		return nil, fmt.Errorf("could not parse path in ,GetValueFromJSON error: '%s'", err)
	}

	res := jsonPath.Get(obj)

	if len(res) == 0 {
		return nil, fmt.Errorf("could not find element by path %v in JSON", expression)
	}

	return res[0], nil
}
