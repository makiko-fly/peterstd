package json

import "github.com/json-iterator/go"

type RawMessage = jsoniter.RawMessage

func Marshal(v interface{}) ([]byte, error) {
	return jsoniter.Marshal(v)
}

func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return jsoniter.MarshalIndent(v, prefix, indent)
}

func MarshalToString(v interface{}) (string, error) {
	return jsoniter.MarshalToString(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return jsoniter.Unmarshal(data, v)
}

// Must 系列函数，忽略错误

func MustMarshal(v interface{}) []byte {
	b, _ := jsoniter.Marshal(v)
	return b
}

func MustMarshalIndent(v interface{}, prefix, indent string) []byte {
	b, _ := jsoniter.MarshalIndent(v, prefix, indent)
	return b
}

func MustMarshalToString(v interface{}) string {
	s, _ := jsoniter.MarshalToString(v)
	return s
}

func JSONCopy(dst, src interface{}) error {
	bytes, err := jsoniter.Marshal(src)
	if err != nil {
		return err
	}
	return jsoniter.Unmarshal(bytes, dst)
}
