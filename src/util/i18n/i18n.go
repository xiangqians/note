// @author xiangqian
// @date 15:53 2023/11/19
package i18n

import (
	"embed"
	util_json "note/src/util/json"
)

type Language string

const (
	ZH Language = "zh"
	EN Language = "en"
)

var zhMessageMap map[string]string
var enMessageMap map[string]string

func Init(fs embed.FS) {
	// zh
	bytes, err := fs.ReadFile("embed/i18n/zh.json")
	if err != nil {
		panic(err)
	}
	err = util_json.Deserialize(bytes, &zhMessageMap)
	if err != nil {
		panic(err)
	}

	// en
	bytes, err = fs.ReadFile("embed/i18n/en.json")
	if err != nil {
		panic(err)
	}
	err = util_json.Deserialize(bytes, &enMessageMap)
	if err != nil {
		panic(err)
	}
}

func GetMessage(name string, language Language) string {
	switch language {
	case ZH:
		return zhMessageMap[name]

	case EN:
		return enMessageMap[name]

	default:
		return ""
	}
}
