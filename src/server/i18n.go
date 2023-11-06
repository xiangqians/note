// i18n
// @author xiangqian
// @date 23:16 2023/07/18
package server

import (
	"encoding/json"
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"note/src/context"
	"note/src/session"
	"strings"
)

const (
	ZH = "zh"
	EN = "en"
)

// 初始化i18n
// https://github.com/gin-contrib/i18n
func initI18n(engine *gin.Engine) {
	engine.Use(i18n.Localize(i18n.WithBundle(&i18n.BundleCfg{
		RootPath:         "./res/i18n",
		AcceptLanguage:   []language.Tag{language.Chinese, language.English},
		DefaultLanguage:  language.Chinese,
		UnmarshalFunc:    json.Unmarshal,
		FormatBundleFile: "json",
	}), i18n.WithGetLngHandle(
		func(ctx *gin.Context, defaultLanguage string) string {
			name := "language"
			// 从url中获取language
			language, _ := context.Query[string](ctx, name)
			language = strings.ToLower(language)
			switch language {
			case ZH:
				language = ZH
			case EN:
				language = EN
			default:
				language = ""
			}

			// 从请求头获取 Accept-Language
			if language == "" {
				// 从请求头获取 Accept-Language
				acceptLanguage, _ := context.Header[string](ctx, "Accept-Language")
				// en,zh-CN;q=0.9,zh;q=0.8
				if strings.HasPrefix(acceptLanguage, ZH) {
					language = ZH
				} else if strings.HasPrefix(acceptLanguage, EN) {
					language = EN
				}
			}

			// 从session中获取language
			if language == "" {
				sessionLanguage, err := session.Get[string](ctx, name, false)
				if err == nil {
					switch sessionLanguage {
					case ZH:
						return ZH
					case EN:
						return EN
					}
				}
			}

			// 如果language未指定，则使用默认language
			if language == "" {
				language = defaultLanguage
			}

			// 存储language到session
			session.Set(ctx, name, language)

			return language
		},
	)))
}
