package utils

import (
	"bytes"
	"encoding/json"
)

func ToJSON(v any) *bytes.Buffer {

	b, _ := json.Marshal(v)

	return bytes.NewBuffer(b)
}
