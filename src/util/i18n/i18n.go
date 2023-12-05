// @author xiangqian
// @date 15:53 2023/11/19
package i18n

import (
	"embed"
	"fmt"
	util_json "note/src/util/json"
	"os"
)

const (
	ZH string = "zh"
	EN string = "en"
)

var zhMessageMap map[string]string
var enMessageMap map[string]string

// 模式，dev、test、prod
var mode = os.Getenv("GO_ENV_MODE")

// 项目目录
var projectDir, _ = os.Getwd()

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
