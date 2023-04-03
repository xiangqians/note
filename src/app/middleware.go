// Middleware
// @author xiangqian
// @date 21:46 2022/12/23
package app

import (
	"encoding/json"
	"github.com/gin-contrib/i18n"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"net/http"
	api_common "note/src/api/common"
	typ_locale "note/src/typ/locale"
	typ_resp "note/src/typ/resp"
	"strings"
)

// 权限中间件
func permMiddleware(engine *gin.Engine) {
	// 未授权拦截
	engine.Use(func(context *gin.Context) {
		reqPath := context.Request.URL.Path
		reqMethod := context.Request.Method

		// 静态资源放行
		if strings.HasPrefix(reqPath, "/static") {
			context.Next()
			return
		}

		// isLogin
		isLogin := false
		_, err := api_common.GetSessionUser(context)
		if err == nil {
			isLogin = true
		}

		// 用户注册和登录放行
		if reqPath == "/user/reg" || reqPath == "/user/login" ||
			(reqMethod == http.MethodPost && (reqPath == "/user" || reqPath == "/user/login0")) {
			if isLogin {
				api_common.Redirect(context, "/", typ_resp.Resp[any]{})
				context.Abort()
			} else {
				context.Next()
			}
			return
		}

		if !isLogin {
			// 重定向
			//context.Request.URL.Path = "/user/login"
			//engine.HandleContext(context)
			//
			api_common.Redirect(context, "/user/login", typ_resp.Resp[any]{})

			// 中止调用链
			context.Abort()
			return
		}
	})
}

// 静态资源处理中间件
func staticMiddleware(engine *gin.Engine) {
	// 静态资源处理
	// https://github.com/gin-contrib/static
	engine.Use(static.Serve("/static", static.LocalFile("./static", false)))
}

// i18n中间件
func i18nMiddleware(engine *gin.Engine) {
	// apply i18n middleware
	// https://github.com/gin-contrib/i18n
	engine.Use(i18n.Localize(i18n.WithBundle(&i18n.BundleCfg{
		RootPath:         "./i18n",
		AcceptLanguage:   []language.Tag{language.Chinese, language.English},
		DefaultLanguage:  language.Chinese,
		UnmarshalFunc:    json.Unmarshal,
		FormatBundleFile: "json",
	}), i18n.WithGetLngHandle(
		func(context *gin.Context, defaultLang string) string {
			// 从url中获取lang
			lang := strings.ToLower(strings.TrimSpace(context.Query("lang")))
			if lang != "" && !(lang == typ_locale.Zh || lang == typ_locale.En) {
				lang = ""
			}

			// 从session中获取lang
			session := api_common.Session(context)
			sessionLang := ""
			if v, r := session.Get("lang").(string); r {
				sessionLang = v
			}
			if lang == "" {
				lang = sessionLang
			}

			if lang == "" {
				// 从请求头获取 Accept-Language
				acceptLanguage := context.GetHeader("Accept-Language")
				// en,zh-CN;q=0.9,zh;q=0.8
				if strings.HasPrefix(acceptLanguage, typ_locale.Zh) {
					lang = typ_locale.Zh
				} else if strings.HasPrefix(acceptLanguage, typ_locale.En) {
					lang = typ_locale.En
				}
			}

			if lang == "" {
				lang = defaultLang
			}

			if sessionLang != lang {
				session.Set("lang", lang)
				session.Save()
			}
			return lang
		},
	)))
}

// session中间件
func sessionMiddleware(engine *gin.Engine) {
	// 密钥
	keyPairs := []byte("123456")

	// 创建基于cookie的存储引擎
	//store := cookie.NewStore(keyPairs)
	// 创建基于mem（内存）的存储引擎，其实就是一个 map[interface]interface 对象
	store := memstore.NewStore(keyPairs)

	// store配置
	store.Options(sessions.Options{
		//Secure: true,
		//SameSite: http.SameSiteNoneMode,
		Path:   "/",
		MaxAge: 60 * 60 * 12, // 12h，设置session过期时间，seconds
	})

	// 设置session中间件
	// session中间件基于内存（其他存储引擎支持：redis、mysql等）实现
	engine.Use(sessions.Sessions("NoteSessionId", // session & cookie 名称
		store))
}
