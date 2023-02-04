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
	"note/src/api"
	"note/src/typ"
	"strings"
)

func permMiddleware(pEngine *gin.Engine) {
	// 未授权拦截
	pEngine.Use(func(pContext *gin.Context) {
		reqPath := pContext.Request.URL.Path

		// 静态资源放行
		if strings.HasPrefix(reqPath, "/static") {
			pContext.Next()
			return
		}

		// isLogin
		isLogin := false
		_, err := api.SessionUser(pContext)
		if err == nil {
			isLogin = true
		}

		if reqPath == "/user/regpage" || reqPath == "/user/loginpage" ||
			(reqPath == "/user" && pContext.Request.Method == http.MethodPost) || reqPath == "/user/login" {
			if isLogin {
				pContext.Redirect(http.StatusMovedPermanently, "/")
				pContext.Abort()
			} else {
				pContext.Next()
			}
			return
		}

		if !isLogin {
			// 重定向
			//pContext.Request.URL.Path = "/user/loginpage"
			//pEngine.HandleContext(pContext)
			pContext.Redirect(http.StatusMovedPermanently, "/user/loginpage")

			// 中止调用链
			pContext.Abort()
			return
		}
	})
}

func staticMiddleware(pEngine *gin.Engine) {
	// 静态资源处理
	// https://github.com/gin-contrib/static
	pEngine.Use(static.Serve("/static", static.LocalFile("./static", false)))
}

func i18nMiddleware(pEngine *gin.Engine) {
	// apply i18n middleware
	// https://github.com/gin-contrib/i18n
	pEngine.Use(i18n.Localize(i18n.WithBundle(&i18n.BundleCfg{
		RootPath:         "./i18n",
		AcceptLanguage:   []language.Tag{language.Chinese, language.English},
		DefaultLanguage:  language.Chinese,
		UnmarshalFunc:    json.Unmarshal,
		FormatBundleFile: "json",
	}), i18n.WithGetLngHandle(
		func(pContext *gin.Context, defaultLang string) string {
			// 从url中获取lang
			lang := strings.ToLower(strings.TrimSpace(pContext.Query("lang")))
			if lang != "" && !(lang == typ.LocaleZh || lang == typ.LocaleEn) {
				lang = ""
			}

			// 从session中获取lang
			//session := sessions.Default(pContext)
			session := sessions.Default(pContext)
			sessionLang := ""
			if v, r := session.Get("lang").(string); r {
				sessionLang = v
			}
			if lang == "" {
				lang = sessionLang
			}

			if lang == "" {
				// 从请求头获取 Accept-Language
				acceptLanguage := pContext.GetHeader("Accept-Language")
				// en,zh-CN;q=0.9,zh;q=0.8
				if strings.HasPrefix(acceptLanguage, typ.LocaleZh) {
					lang = typ.LocaleZh
				} else if strings.HasPrefix(acceptLanguage, typ.LocaleEn) {
					lang = typ.LocaleEn
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

func sessionMiddleware(pEngine *gin.Engine) {
	// 密钥
	keyPairs := []byte("123456")
	// 创建基于cookie的存储引擎
	//store := cookie.NewStore(keyPairs)
	// 创建基于mem（内存）的存储引擎，其实就是一个 map[interface]interface 对象
	store := memstore.NewStore(keyPairs)

	// store配置
	store.Options(sessions.Options{
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		MaxAge:   60 * 60 * 24, // 设置session过期时间，seconds
	})

	// 设置session中间件
	// session中间件基于内存（其他存储引擎支持：redis、mysql等）实现
	pEngine.Use(sessions.Sessions("NoteSessionId", // session & cookie 名称
		store))
}
