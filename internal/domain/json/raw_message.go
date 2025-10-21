package json

import (
	jsonstd "encoding/json"
	"reflect"
)

type RawMessage jsonstd.RawMessage

func (m *RawMessage) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	if len(*m) == 0 {
		return []byte("null"), nil
	}
	return *m, nil
}

func (m *RawMessage) UnmarshalJSON(data []byte) error {
	if m == nil {
		return &jsonstd.InvalidUnmarshalError{Type: reflect.TypeOf(m)}
	}
	*m = append((*m)[0:0], data...)
	return nil
}
