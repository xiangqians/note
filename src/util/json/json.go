// @author xiangqian
// @date 13:26 2023/04/02
package json

import (
	"encoding/json"
)

// Serialize 使用Marshal序列化
// v      : 要序列号的实例
// indent : 是否缩进
func Serialize(v any, indent bool) (string, error) {
	var bytes []byte
	var err error
	if indent {
		bytes, err = json.MarshalIndent(v,
			"",   // 指定每行输出开头的字符串
			"\t") // 指定每行要缩进的字符串
	} else {
		bytes, err = json.Marshal(v)
	}
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// Deserialize 使用Unmarshal反序列化
func Deserialize(text string, v any) error {
	return json.Unmarshal([]byte(text), v)
}
