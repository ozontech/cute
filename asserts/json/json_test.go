package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type jsonTest struct {
	caseName   string
	data       string
	expression string
	expect     interface{}
	IsNilErr   bool
}

func TestDiff(t *testing.T) {
	testCases := []struct {
		name          string
		originalJSON  string
		bodyJSON      string
		expectedError string
	}{
		{
			name:          "SameJSON",
			originalJSON:  `{"key1": "value1", "key2": "value2"}`,
			bodyJSON:      `{"key1": "value1", "key2": "value2"}`,
			expectedError: "", // No error expected, JSONs are the same
		},
		{
			name:          "DifferentValueJSON",
			originalJSON:  `{"key1": "value1", "key2": "value2"}`,
			bodyJSON:      `{"key1": "value1", "key2": "value3"}`,
			expectedError: "JSON is not the same",
		},
		{
			name:          "MissingKeyJSON",
			originalJSON:  `{"key1": "value1", "key2": "value2"}`,
			bodyJSON:      `{"key1": "value1"}`,
			expectedError: "JSON is not the same",
		},
		{
			name:          "ExtraKeyJSON",
			originalJSON:  `{"key1": "value1"}`,
			bodyJSON:      `{"key1": "value1", "key2": "value2"}`,
			expectedError: "JSON is not the same",
		},
		{
			name:          "EmptyJSON",
			originalJSON:  `{}`,
			bodyJSON:      `{}`,
			expectedError: "", // No error expected, empty JSONs are the same
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Call the Diff function with the test input
			err := Diff(testCase.originalJSON)([]byte(testCase.bodyJSON))

			// Check if the error message matches the expected result
			if testCase.expectedError == "" {
				assert.NoError(t, err) // No error expected
			} else {
				assert.Error(t, err) // Error expected
				assert.Contains(t, err.Error(), testCase.expectedError)
			}
		})
	}
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
			expression: "$.o[0]['1']",
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
			expression: "$.o[0]['1']",
			IsNilErr:   true,
		},
		{
			caseName:   "check not correct path",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.not_correct",
		},
		{
			caseName:   "empty integer",
			data:       `{"o":0}`,
			expression: "$.o",
			IsNilErr:   true,
		},
		{
			caseName:   "empty object",
			data:       `{"o":null}`,
			expression: "$.o",
			IsNilErr:   true,
		},
		{
			caseName:   "empty string",
			data:       `{"o":null, "b":""}`,
			expression: "$.b",
			IsNilErr:   true,
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

func TestNotEmpty(t *testing.T) {
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
			expression: "$.o[0]['1']",
			IsNilErr:   true,
		},
		{
			caseName:   "check not correct path",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.not_correct",
		},
		{
			caseName:   "empty integer",
			data:       `{"o":0}`,
			expression: "$.o",
		},
		{
			caseName:   "empty object",
			data:       `{"o":null}`,
			expression: "$.o",
		},
		{
			caseName:   "empty string",
			data:       `{"o":null, "b":""}`,
			expression: "$.b",
		},
	}

	for _, test := range tests {
		err := NotEmpty(test.expression)([]byte(test.data))

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

func TestLengthGreaterThan(t *testing.T) {
	tests := []jsonTest{
		{
			caseName:   "correct check array",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.o",
			expect:     2,
			IsNilErr:   true,
		},
		{
			caseName:   "not correct check array",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.o",
			expect:     4,
		},
		{
			caseName:   "not correct check array when equal",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.o",
			expect:     3,
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
		err := LengthGreaterThan(test.expression, test.expect.(int))([]byte(test.data))

		if test.IsNilErr {
			require.NoError(t, err, "failed test %v", test.caseName)
		} else {
			require.Error(t, err, "failed test %v", test.caseName)
		}
	}
}

func TestLengthGreaterOrEqualThan(t *testing.T) {
	tests := []jsonTest{
		{
			caseName:   "correct check array",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.o",
			expect:     2,
			IsNilErr:   true,
		},
		{
			caseName:   "correct check array when equal",
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
			caseName:   "correct check string when equal",
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
			caseName:   "correct check map",
			data:       `{"o":[{"1":"a"}, {"2":"b"}, {"3":"c"}]}`,
			expression: "$.o",
			expect:     1,
			IsNilErr:   true,
		},
		{
			caseName:   "correct check map when equal",
			data:       `{"o":[{"1":"a"}, {"2":"b"}, {"3":"c"}]}`,
			expression: "$.o",
			expect:     3,
			IsNilErr:   true,
		},
		{
			caseName:   "not correct check map",
			data:       `{"o":[{"1":"a"}, {"2":"b"}, {"3":"c"}]}`,
			expression: "$.o",
			expect:     5,
		},
		{
			caseName:   "check not correct path",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.not_correct",
			expect:     0,
		},
	}

	for _, test := range tests {
		err := LengthGreaterOrEqualThan(test.expression, test.expect.(int))([]byte(test.data))

		if test.IsNilErr {
			require.NoError(t, err, "failed test %v", test.caseName)
		} else {
			require.Error(t, err, "failed test %v", test.caseName)
		}
	}
}

func TestLengthLessThan(t *testing.T) {
	tests := []jsonTest{
		{
			caseName:   "correct check array",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.o",
			expect:     4,
			IsNilErr:   true,
		},
		{
			caseName:   "not correct check array",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.o",
			expect:     3,
		},
		{
			caseName:   "correct check string",
			data:       `{"o":"123456"}`,
			expression: "$.o",
			expect:     7,
			IsNilErr:   true,
		},
		{
			caseName:   "not correct check string",
			data:       `{"o":"123456"}`,
			expression: "$.o",
			expect:     6,
		},
		{
			caseName:   "correct check map",
			data:       `{"o":[{"1":"a"}, {"2":"b"}, {"3":"c"}]}`,
			expression: "$.o",
			expect:     4,
			IsNilErr:   true,
		},
		{
			caseName:   "not correct check map",
			data:       `{"o":[{"1":"a"}, {"2":"b"}, {"3":"c"}]}`,
			expression: "$.o",
			expect:     3,
		},
		{
			caseName:   "check not correct path",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.not_correct",
			expect:     0,
		},
	}

	for _, test := range tests {
		err := LengthLessThan(test.expression, test.expect.(int))([]byte(test.data))

		if test.IsNilErr {
			require.NoError(t, err, "failed test %v", test.caseName)
		} else {
			require.Error(t, err, "failed test %v", test.caseName)
		}
	}
}

func TestLengthLessOrEqualThan(t *testing.T) {
	tests := []jsonTest{
		{
			caseName:   "correct check array",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.o",
			expect:     4,
			IsNilErr:   true,
		},
		{
			caseName:   "correct check array when equal",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.o",
			expect:     3,
			IsNilErr:   true,
		},
		{
			caseName:   "not correct check array",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.o",
			expect:     2,
		},
		{
			caseName:   "correct check string",
			data:       `{"o":"123456"}`,
			expression: "$.o",
			expect:     7,
			IsNilErr:   true,
		},
		{
			caseName:   "correct check string when equal",
			data:       `{"o":"123456"}`,
			expression: "$.o",
			expect:     6,
			IsNilErr:   true,
		},
		{
			caseName:   "not correct check string",
			data:       `{"o":"123456"}`,
			expression: "$.o",
			expect:     5,
		},
		{
			caseName:   "correct check map",
			data:       `{"o":[{"1":"a"}, {"2":"b"}, {"3":"c"}]}`,
			expression: "$.o",
			expect:     4,
			IsNilErr:   true,
		},
		{
			caseName:   "correct check map when equal",
			data:       `{"o":[{"1":"a"}, {"2":"b"}, {"3":"c"}]}`,
			expression: "$.o",
			expect:     3,
			IsNilErr:   true,
		},
		{
			caseName:   "not correct check map",
			data:       `{"o":[{"1":"a"}, {"2":"b"}, {"3":"c"}]}`,
			expression: "$.o",
			expect:     2,
		},
		{
			caseName:   "check not correct path",
			data:       `{"o":["a", "b", "c"]}`,
			expression: "$.not_correct",
			expect:     0,
		},
	}

	for _, test := range tests {
		err := LengthLessOrEqualThan(test.expression, test.expect.(int))([]byte(test.data))

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
			caseName:   "valid array",
			data:       `{"arr": ["one","two"]}`,
			expression: "$.arr",
			expect:     []string{"one", "two"},
			IsNilErr:   true,
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
		{
			caseName:   "check 186135434",
			data:       `{"a":186135434, "b":{"bs":"sb"}}`,
			expression: "$.a",
			expect:     186135434,
			IsNilErr:   true,
		},
		{
			caseName:   "check float",
			data:       `{"a":1.0000001, "b":{"bs":"sb"}}`,
			expression: "$.a",
			expect:     1.0000001,
			IsNilErr:   true,
		},
		{
			caseName:   "check float 2",
			data:       `{"a":999.0000001, "b":{"bs":"sb"}}`,
			expression: "$.a",
			expect:     999.0000001,
			IsNilErr:   true,
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

func TestGetValueFromJSON(t *testing.T) {
	testCases := []struct {
		name          string
		inputJSON     string
		expression    string
		expectedValue []interface{}
		expectedError string
	}{
		{
			name:          "ValidExpressionObject",
			inputJSON:     `{"key1": "value1", "key2": {"key3": "value3"}}`,
			expression:    "key2.key3",
			expectedValue: []interface{}{"value3"},
			expectedError: "", // No error expected
		},
		{
			name:          "ValidExpressionArray",
			inputJSON:     `{"key1": "value1", "key2": [1, 2, 3]}`,
			expression:    "key2[1]",
			expectedValue: []interface{}{int64(2)},
			expectedError: "", // No error expected
		},
		{
			name:          "ValidExpressionMap",
			inputJSON:     `{"key1": "value1", "key2": {"subkey1": "subvalue1"}}`,
			expression:    "key2",
			expectedValue: []interface{}{map[string]interface{}{"subkey1": "subvalue1"}},
			expectedError: "", // No error expected
		},
		{
			name:          "InvalidJSON",
			inputJSON:     `invalid json`,
			expression:    "key1",
			expectedValue: nil,
			expectedError: "could not parse json",
		},
		{
			name:          "InvalidExpression",
			inputJSON:     `{"key1": "value1"}`,
			expression:    "key2",
			expectedValue: nil,
			expectedError: "could not find element by path key2 in JSON",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Call the GetValueFromJSON function with the test input
			value, err := GetValueFromJSON([]byte(testCase.inputJSON), testCase.expression)

			// Check if the error message matches the expected result
			if testCase.expectedError == "" {
				assert.NoError(t, err) // No error expected
			} else {
				assert.Error(t, err) // Error expected
				assert.Contains(t, err.Error(), testCase.expectedError)
			}

			// Check if the returned value is an array and matches the expected result
			assert.IsType(t, []interface{}{}, value)
			assert.Equal(t, testCase.expectedValue, value)
		})
	}
}
