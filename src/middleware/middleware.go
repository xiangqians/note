// Middleware
// @author xiangqian
// @date 21:46 2022/12/23
package middleware

import (
	"encoding/json"
	"github.com/gin-contrib/i18n"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"log"
	"net/http"
	"note/src/arg"
	src_context "note/src/context"
	"note/src/session"
	"note/src/typ"
	"note/src/util/crypto/bcrypt"
	"strings"
)

// Init 初始化中间件
func Init(engine *gin.Engine) {
	sessionMiddleware(engine)
	i18nMiddleware(engine)
	staticMiddleware(engine)
	permMiddleware(engine)
}

// 权限中间件
func permMiddleware(engine *gin.Engine) {
	// 未授权拦截
	engine.Use(func(context *gin.Context) {
		// 请求路径
		reqPath := context.Request.URL.Path

		// 静态资源放行
		if strings.HasPrefix(reqPath, arg.Arg.Path+"/static/") {
			context.Next()
			return
		}

		// 是否已登录
		signIn := false
		user, err := session.GetUser(context)
		if err == nil && user.Id > 0 {
			signIn = true
		}

		// 用户登录和注册放行
		if reqPath == arg.Arg.Path+"/user/signin" || // 登录页
			reqPath == arg.Arg.Path+"/user/signup" || // 注册页
			(context.Request.Method == http.MethodPost && (reqPath == arg.Arg.Path+"/user/signin0" || reqPath == arg.Arg.Path+"/user/signup0")) { // 登录接口和注册接口
			// 如果已登录则重定向到首页
			if signIn {
				// 重定向到首页
				src_context.Redirect(context, arg.Arg.Path+"/")
				// 中止调用链
				context.Abort()
			} else
			// 如果未登录，放行登录或注册
			{
				context.Next()
			}
			return
		}

		// 未登录
		if !signIn {
			// 重定向到登录页
			//context.Request.URL.Path = arg.Arg.Path+"/user/signin"
			//engine.HandleContext(context)
			// OR
			src_context.Redirect(context, arg.Arg.Path+"/user/signin")
			// 中止调用链
			context.Abort()
			return
		}
	})
}

// 静态资源处理中间件
// https://github.com/gin-contrib/static
func staticMiddleware(engine *gin.Engine) {
	engine.Use(static.Serve(arg.Arg.Path+"/static", static.LocalFile("./res/static", false)))
}

// i18n中间件
// https://github.com/gin-contrib/i18n
func i18nMiddleware(engine *gin.Engine) {
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

// session中间件
func sessionMiddleware(engine *gin.Engine) {
	// 密钥
	passwd := "$2a$10$NkWzRTyz1ZNnNfjLmxreaeZ31DCiwCEWJlXJAVDkG8fD9Ble2mg4K"
	hash, err := bcrypt.Generate(passwd)
	if err != nil {
		log.Println(err)
		hash = passwd
	}
	keyPairs := []byte(hash)[:32]

	// session存储引擎支持：基于内存、redis、mysql等
	// 1、创建基于cookie的存储引擎
	//store := cookie.NewStore(keyPairs)
	// 2、创建基于mem（内存）的存储引擎，其实就是一个 map[interface]interface 对象
	//store := memstore.NewStore(keyPairs)
	store := session.NewStore(keyPairs)

	// store配置
	store.Options(sessions.Options{
		//Secure: true,
		//SameSite: http.SameSiteNoneMode,
		Path:   "/",
		MaxAge: 60 * 60 * 12, // 12h，设置session过期时间，seconds
	})

	// 设置session中间件
	engine.Use(sessions.Sessions("NoteSessionId", // session & cookie 名称
		store))
}
