// i18n
// @author xiangqian
// @date 23:16 2023/07/18
package i18n

import (
	"encoding/json"
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"note/src/session"
	"note/src/typ"
	"strings"
)

// Init 初始化i18n
// https://github.com/gin-contrib/i18n
func Init(engine *gin.Engine) {
	engine.Use(i18n.Localize(i18n.WithBundle(&i18n.BundleCfg{
		RootPath:         "./res/i18n",
		AcceptLanguage:   []language.Tag{language.Chinese, language.English},
		DefaultLanguage:  language.Chinese,
		UnmarshalFunc:    json.Unmarshal,
		FormatBundleFile: "json",
	}), i18n.WithGetLngHandle(
		func(context *gin.Context, defaultLang string) string {
			// 从url中获取lang
			lang := strings.ToLower(strings.TrimSpace(context.Query("lang")))
			if lang != "" && !(lang == typ.Zh || lang == typ.En) {
				lang = ""
			}

			// 从session中获取lang
			sessionLang, err := session.Get[string](context, "lang", false)
			if err != nil {
				sessionLang = ""
			}
			if lang == "" {
				lang = sessionLang
			}

			// 从请求头获取 Accept-Language
			if lang == "" {
				// 从请求头获取 Accept-Language
				acceptLanguage := context.GetHeader("Accept-Language")
				// en,zh-CN;q=0.9,zh;q=0.8
				if strings.HasPrefix(acceptLanguage, typ.Zh) {
					lang = typ.Zh
				} else if strings.HasPrefix(acceptLanguage, typ.En) {
					lang = typ.En
				}
			}

			// 如果lang未指定，则使用默认lang
			if lang == "" {
				lang = defaultLang
			}

			// 存储lang到session
			if sessionLang != lang {
				session.Set(context, "lang", lang)
			}

			return lang
		},
	)))
}
