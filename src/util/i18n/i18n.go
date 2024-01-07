// @author xiangqian
// @date 15:53 2023/11/19
package i18n

import (
	"fmt"
	"note/src/embed"
	util_json "note/src/util/json"
	"os"
)

const (
	ZH string = "zh"
	EN string = "en"
)

// 模式，dev、test、prod
var mode = os.Getenv("NOTE_MODE")

// 项目目录
var projectDir, _ = os.Getwd()

var zhMessageMap map[string]string
var enMessageMap map[string]string

func GetMessage(name, language string) string {
	if mode == "dev" {
		// zh
		bytes, err := os.ReadFile(fmt.Sprintf("%s/src/embed/i18n/zh.json", projectDir))
		if err != nil {
			panic(err)
		}
		err = util_json.Deserialize(bytes, &zhMessageMap)
		if err != nil {
			panic(err)
		}

		// en
		bytes, err = os.ReadFile(fmt.Sprintf("%s/src/embed/i18n/en.json", projectDir))
		if err != nil {
			panic(err)
		}
		err = util_json.Deserialize(bytes, &enMessageMap)
		if err != nil {
			panic(err)
		}

	} else {
		// zh
		if zhMessageMap == nil {
			bytes, err := embed.I18nFs.ReadFile("embed/i18n/zh.json")
			if err != nil {
				panic(err)
			}
			err = util_json.Deserialize(bytes, &zhMessageMap)
			if err != nil {
				panic(err)
			}
		}

		// en
		if enMessageMap == nil {
			bytes, err := embed.I18nFs.ReadFile("embed/i18n/en.json")
			if err != nil {
				panic(err)
			}
			err = util_json.Deserialize(bytes, &enMessageMap)
			if err != nil {
				panic(err)
			}
		}
	}

	switch language {
	case ZH:
		return zhMessageMap[name]

	case EN:
		return enMessageMap[name]

	default:
		panic(fmt.Sprintf("%s[%s] undefined", language, name))
	}
}
