package cute

import "encoding/json"

// JSONMarshaler is marshaler which use for marshal/unmarshal JSON to/from struct
type JSONMarshaler interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

type jsonMarshaler struct {
}

func (j jsonMarshaler) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (j jsonMarshaler) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
