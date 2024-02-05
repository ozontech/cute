package utils

import (
	"bytes"
	"encoding/json"
)

// ToJSON returns string Json representation of any object that can be marshaled.
func ToJSON(v interface{}) (string, error) {
	j, err := json.Marshal(v)

	return string(j), err
}

// PrettyJSON make indent to json byte array. Returns prettified json as []byte or error if is it impossible
func PrettyJSON(b []byte) ([]byte, error) {
	var out bytes.Buffer

	err := json.Indent(&out, b, "", "    ")

	return out.Bytes(), err
}
