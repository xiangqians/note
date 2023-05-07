// json util
// @author xiangqian
// @date 13:26 2023/04/02
package json

import (
	_json "encoding/json"
)

// Serialize 使用 Marshal 序列化
func Serialize(v any) (string, error) {
	buf, err := _json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

// Deserialize 使用 Unmarshal 反序列化
func Deserialize(text string, v any) error {
	return _json.Unmarshal([]byte(text), v)
}
