// @author xiangqian
// @date 15:53 2023/11/19
package i18n

import (
	"embed"
	"fmt"
	"note/src/model"
	util_json "note/src/util/json"
	"os"
)

const (
	ZH string = "zh"
	EN string = "en"
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

func GetMessage(name, language string) string {
	if model.GetMode() == model.ModeDev {
		// zh
		bytes, err := os.ReadFile(fmt.Sprintf("%s/src/embed/i18n/zh.json", model.GetProjectDir()))
		if err != nil {
			panic(err)
		}
		err = util_json.Deserialize(bytes, &zhMessageMap)
		if err != nil {
			panic(err)
		}

		// en
		bytes, err = os.ReadFile(fmt.Sprintf("%s/src/embed/i18n/en.json", model.GetProjectDir()))
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
