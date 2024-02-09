package cute

import "encoding/json"

type jsonMarshaler struct {
}

func (j jsonMarshaler) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (j jsonMarshaler) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
