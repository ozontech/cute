package json

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type jsonTest struct {
	caseName   string
	data       string
	expression string
	expect     interface{}
	IsNilErr   bool
}

func TestNotPresent(t *testing.T) {
	tests := []jsonTest{
		{
			caseName:   "correct check array",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.o",
		},
		{
			caseName:   "not present check",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.b",
			IsNilErr:   true,
		},
		{
			caseName:   "not present check ",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.o[0]",
		},
		{
			caseName:   "correct check map",
			data:       `{"o":[{"1":"a"}, {"2":"b"}, {"3":"c"}]}`,
			expression: "$.o[0][1]",
		},
	}

	for _, test := range tests {
		err := NotPresent(test.expression)([]byte(test.data))

		if test.IsNilErr {
			require.NoError(t, err, "failed test %v", test.caseName)
		} else {
			require.Error(t, err, "failed test %v", test.caseName)
		}
	}
}

func TestPresent(t *testing.T) {
	tests := []jsonTest{
		{
			caseName:   "correct check array",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.o",
			IsNilErr:   true,
		},
		{
			caseName:   "not present check",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.b",
		},
		{
			caseName:   "correct present check array",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.o[0]",
			IsNilErr:   true,
		},
		{
			caseName:   "correct check map",
			data:       `{"o":[{"1":"a"}, {"2":"b"}, {"3":"c"}]}`,
			expression: "$.o[0][1]",
			IsNilErr:   true,
		},
		{
			caseName:   "check not correct path",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.not_correct",
		},
	}

	for _, test := range tests {
		err := Present(test.expression)([]byte(test.data))

		if test.IsNilErr {
			require.NoError(t, err, "failed test %v", test.caseName)
		} else {
			require.Error(t, err, "failed test %v", test.caseName)
		}
	}
}

func TestLength(t *testing.T) {
	tests := []jsonTest{
		{
			caseName:   "correct check array",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.o",
			expect:     3,
			IsNilErr:   true,
		},
		{
			caseName:   "not correct check array",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.o",
			expect:     4,
		},
		{
			caseName:   "correct check string",
			data:       `{"o":"123456"}`,
			expression: "$.o",
			expect:     6,
			IsNilErr:   true,
		},
		{
			caseName:   "not correct check string",
			data:       `{"o":"123456"}`,
			expression: "$.o",
			expect:     99,
		},
		{
			caseName:   "check not contain value",
			data:       `{"o":"123456"}`,
			expression: "$.a",
			expect:     1,
		},
		{
			caseName:   "correct check map",
			data:       `{"o":[{"1":"a"}, {"2":"b"}, {"3":"c"}]}`,
			expression: "$.o",
			expect:     3,
			IsNilErr:   true,
		},
	}

	for _, test := range tests {
		err := Length(test.expression, test.expect.(int))([]byte(test.data))

		if test.IsNilErr {
			require.NoError(t, err, "failed test %v", test.caseName)
		} else {
			require.Error(t, err, "failed test %v", test.caseName)
		}
	}
}

func TestGreaterThan(t *testing.T) {
	tests := []jsonTest{
		{
			caseName:   "correct check array",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.o",
			expect:     3,
			IsNilErr:   true,
		},
		{
			caseName:   "not correct check array",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.o",
			expect:     4,
		},
		{
			caseName:   "correct check string",
			data:       `{"o":"123456"}`,
			expression: "$.o",
			expect:     4,
			IsNilErr:   true,
		},
		{
			caseName:   "not correct check string",
			data:       `{"o":"123456"}`,
			expression: "$.o",
			expect:     99,
		},
		{
			caseName:   "correct check map",
			data:       `{"o":[{"1":"a"}, {"2":"b"}, {"3":"c"}]}`,
			expression: "$.o",
			expect:     1,
			IsNilErr:   true,
		},
		{
			caseName:   "check not correct path",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.not_correct",
			expect:     0,
		},
	}

	for _, test := range tests {
		err := GreaterThan(test.expression, test.expect.(int))([]byte(test.data))

		if test.IsNilErr {
			require.NoError(t, err, "failed test %v", test.caseName)
		} else {
			require.Error(t, err, "failed test %v", test.caseName)
		}
	}
}

func TestLessThan(t *testing.T) {
	tests := []jsonTest{
		{
			caseName:   "correct check array",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.o",
			expect:     2,
			IsNilErr:   false,
		},
		{
			caseName:   "not correct check array",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.o",
			expect:     3,
			IsNilErr:   true,
		},
		{
			caseName:   "correct check string",
			data:       `{"o":"123456"}`,
			expression: "$.o",
			expect:     4,
			IsNilErr:   false,
		},
		{
			caseName:   "not correct check string",
			data:       `{"o":"123456"}`,
			expression: "$.o",
			expect:     99,
			IsNilErr:   true,
		},
		{
			caseName:   "correct check map",
			data:       `{"o":[{"1":"a"}, {"2":"b"}, {"3":"c"}]}`,
			expression: "$.o",
			expect:     1,
			IsNilErr:   false,
		},
		{
			caseName:   "check not correct path",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.not_correct",
			expect:     0,
		},
	}

	for _, test := range tests {
		err := LessThan(test.expression, test.expect.(int))([]byte(test.data))

		if test.IsNilErr {
			require.NoError(t, err, "failed test %v", test.caseName)
		} else {
			require.Error(t, err, "failed test %v", test.caseName)
		}
	}
}

func TestEqual(t *testing.T) {
	tests := []jsonTest{
		{
			caseName:   "valid json",
			data:       `{"first": 777, "second": [{"key_1": "some_key", "value": "some_value"}]}`,
			expression: "$.second[0].value",
			expect:     "some_value",
			IsNilErr:   true,
		},
		{
			caseName: "not valid json",
			data:     "{not_valid_json}",
		},
		{
			caseName:   "3rd party key",
			data:       `{"a":"as", "b":{"bs":"sb"}}`,
			expression: "$.l",
			expect:     nil,
		},
		{
			caseName:   "not array",
			data:       `{"a":"as", "b":{"bs":"sb"}}`,
			expression: "$.b[bs]",
		},
		{
			caseName:   "check equal map",
			data:       `{"a":"as", "b":{"bs":"sb"}}`,
			expression: "$.b",
			expect:     map[string]interface{}{"bs": "sb"},
			IsNilErr:   true,
		},
		{
			caseName:   "check equal string",
			data:       `{"a":"as", "b":{"bs":"sb"}}`,
			expression: "$.a",
			expect:     "as",
			IsNilErr:   true,
		},
		{
			caseName:   "check equal not correct string",
			data:       `{"a":"as", "b":{"bs":"sb"}}`,
			expression: "$.a",
			expect:     []byte("not_correct"),
		},
	}

	for _, test := range tests {
		err := Equal(test.expression, test.expect)([]byte(test.data))

		if test.IsNilErr {
			require.NoError(t, err, "failed test %v", test.caseName)
		} else {
			require.Error(t, err, "failed test %v", test.caseName)
		}
	}
}

func TestNotEqual(t *testing.T) {
	tests := []jsonTest{
		{
			caseName:   "valid json",
			data:       `{"first": 777, "second": [{"key_1": "some_key", "value": "some_value"}]}`,
			expression: "$.second[0].value",
			expect:     "some_value",
			IsNilErr:   false,
		},
		{
			caseName: "not valid json",
			data:     "{not_valid_json}",
		},
		{
			caseName:   "3rd party key",
			data:       `{"a":"as", "b":{"bs":"sb"}}`,
			expression: "$.l",
			expect:     nil,
		},
		{
			caseName:   "not array",
			data:       `{"a":"as", "b":{"bs":"sb"}}`,
			expression: "$.b[bs]",
			expect:     "sb",
		},
		{
			caseName:   "check equal map",
			data:       `{"a":"as", "b":{"bs":"sb"}}`,
			expression: "$.b",
			expect:     map[string]interface{}{"bs": "sb"},
			IsNilErr:   false,
		},
		{
			caseName:   "check equal string",
			data:       `{"a":"as", "b":{"bs":"sb"}}`,
			expression: "$.a",
			expect:     "as",
			IsNilErr:   false,
		},
	}

	for _, test := range tests {
		err := NotEqual(test.expression, test.expect)([]byte(test.data))

		if test.IsNilErr {
			require.NoError(t, err, "failed test %v", test.caseName)
		} else {
			require.Error(t, err, "failed test %v", test.caseName)
		}
	}
}
