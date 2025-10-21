package json

func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
