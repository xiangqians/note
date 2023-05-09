// json util
// @author xiangqian
// @date 13:26 2023/04/02
package json

import (
	encoding_json "encoding/json"
)

// Serialize 使用 Marshal 序列化
func Serialize(v any) (string, error) {
	buf, err := encoding_json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

// Deserialize 使用 Unmarshal 反序列化
func Deserialize(text string, v any) error {
	return encoding_json.Unmarshal([]byte(text), v)
}
