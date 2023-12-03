// @author xiangqian
// @date 22:33 2023/11/07
package handler

import (
	"embed"
	"fmt"
	"github.com/gorilla/mux"
	pkg_template "html/template"
	"io/fs"
	"log"
	"net/http"
	"note/src/handler/image"
	"note/src/handler/index"
	"note/src/handler/system"
	"note/src/model"
	"note/src/session"
	"note/src/util/i18n"
	util_os "note/src/util/os"
	util_time "note/src/util/time"
	"strings"
)

// Handle 注册路由到处理器
func Handle(staticFs, templateFs embed.FS, router *mux.Router) {
	// 处理静态资源
	handleStatic(staticFs, router)

	// 处理模板
	handleTemplate(templateFs, router)
}

// 处理模板
func handleTemplate(templateFs embed.FS, router *mux.Router) {
	// 模板函数
	templateFuncMap := pkg_template.FuncMap{
		// 递增
		"Increment": func(i any) any {
			if i, ok := i.(int); ok {
				return i + 1
			}
			if i, ok := i.(int8); ok {
				return i + 1
			}
			if i, ok := i.(int16); ok {
				return i + 1
			}
			if i, ok := i.(int32); ok {
				return i + 1
			}
			if i, ok := i.(int64); ok {
				return i + 1
			}
			return 0
		},

		// 递减
		"Decrement": func(i any) any {
			if i, ok := i.(int); ok {
				return i - 1
			}
			if i, ok := i.(int8); ok {
				return i - 1
			}
			if i, ok := i.(int16); ok {
				return i - 1
			}
			if i, ok := i.(int32); ok {
				return i - 1
			}
			if i, ok := i.(int64); ok {
				return i - 1
			}
			return 0
		},

		// i18n国际化
		"Localize": func(name, language string) string {
			return i18n.GetMessage(name, language)
		},

		"NowUnix": func() int64 {
			return util_time.NowUnix()
		},

		"FormatUnix": func(unix int64) string {
			if unix <= 0 {
				return "-"
			}
			return util_time.FormatTime(util_time.ParseUnix(unix))
		},

		"HumanizUnix": func(unix int64, language string) string {
			return util_time.HumanizUnix(unix, language)
		},

		"HumanizFileSize": func(size int64) string {
			return util_os.HumanizFileSize(size)
		},
	}

	// 上下文路径
	contextPath := model.Ini.Server.ContextPath

	handlerFunc := func(handler func(request *http.Request, session *session.Session) (name string, response model.Response)) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			// 获取会话
			session := session.GetSession(request, writer)

			// 从url中获取语言
			language := strings.TrimSpace(request.URL.Query().Get("language"))
			if language != i18n.ZH && language != i18n.EN {
				language = ""
			}

			// 从session中获取语言
			sessionLanguage := session.GetLanguage()
			if language == "" {
				language = sessionLanguage
			}

			// 从请求头获取语言
			if language == "" {
				// 从请求头获取 Accept-Language
				// en,zh-CN;q=0.9,zh;q=0.8
				language = request.Header.Get("Accept-Language")
				if strings.HasPrefix(language, i18n.ZH) {
					language = i18n.ZH
				} else if strings.HasPrefix(language, i18n.EN) {
					language = i18n.EN
				}
			}

			// 设置默认语言
			if language == "" {
				language = i18n.ZH
			}

			// 保存语言到会话中
			if language != sessionLanguage {
				err := session.SetLanguage(language)
				if err != nil {
					log.Println(err)
				}
			}

			// 请求路径
			path := request.URL.Path

			// 获取当前会话系统信息
			system := session.GetSystem()

			// 登录路径
			signInPath := contextPath + "/signin"

			// 已登录
			if system.Passwd != "" {
				if path == signInPath {
					// 重定向到首页
					http.Redirect(writer, request, fmt.Sprintf("%s/?t=%d", contextPath, util_time.NowUnix()), http.StatusMovedPermanently) // 301 永久重定向
					return
				}
			} else
			// 未登录
			{
				if path != signInPath {
					// 重定向到登录页
					http.Redirect(writer, request, fmt.Sprintf("%s?t=%d", signInPath, util_time.NowUnix()), http.StatusMovedPermanently) // 301 永久重定向
					return
				}
			}

			name, response := handler(request, session)

			// 判断是否进行重定向
			if strings.HasPrefix(name, "redirect:") {
				// 保存消息到会话
				msg := response.Msg
				if msg != "" {
					session.SetMsg(msg)
				}

				// 301 永久重定向
				http.Redirect(writer, request, fmt.Sprintf("%s%s?t=%d", contextPath, name[len("redirect:"):], util_time.NowUnix()), http.StatusMovedPermanently)
				return
			}

			var template *pkg_template.Template
			var err error
			if model.GetMode() == model.ModeDev {
				template, err = pkg_template.New(name).Funcs(templateFuncMap).ParseFiles(fmt.Sprintf("%s/src/embed/template/%s.html", model.GetProjectDir(), name),
					fmt.Sprintf("%s/src/embed/template/common/foot1.html", model.GetProjectDir()),
					fmt.Sprintf("%s/src/embed/template/common/foot2.html", model.GetProjectDir()),
					fmt.Sprintf("%s/src/embed/template/common/table.html", model.GetProjectDir()))
			} else {
				template, err = pkg_template.New(name).Funcs(templateFuncMap).ParseFS(templateFs,
					fmt.Sprintf("embed/template/%s.html", name), // 主模板文件
					"embed/template/common/foot1.html",          // 嵌套模板文件
					"embed/template/common/foot2.html",          // 嵌套模板文件
					"embed/template/common/table.html")          // 嵌套模板文件
			}
			if err != nil {
				panic(err)
			}

			// 执行模板
			err = template.Execute(writer, map[string]any{
				"contextPath": contextPath,        // 上下文路径
				"path":        request.URL.Path,   // 请求路径
				"uri":         request.RequestURI, // 请求uri地址
				"language":    language,           // 语言
				"system":      system,             // 系统信息
				"response":    response,           // 响应数据
			})
			if err != nil {
				panic(err)
			}
		}
	}

	router.NotFoundHandler = handlerFunc(func(request *http.Request, session *session.Session) (name string, response model.Response) {
		return "404", model.Response{}
	})

	// system
	router.HandleFunc(contextPath+"/signin", handlerFunc(system.SignIn))
	router.HandleFunc(contextPath+"/signout", handlerFunc(system.SignOut))
	router.HandleFunc(contextPath+"/setting", handlerFunc(system.Setting))

	// index
	router.HandleFunc(contextPath+"/", handlerFunc(index.Index))

	// image
	router.HandleFunc(contextPath+"/image", handlerFunc(image.List))
	router.HandleFunc(contextPath+"/image/rename", handlerFunc(image.Rename))
	router.HandleFunc(contextPath+"/image/{id}", handlerFunc(image.Rename1))
}

// 处理静态资源
func handleStatic(embedFs embed.FS, router *mux.Router) {
	var handler http.Handler

	if model.GetMode() == model.ModeDev {
		// 创建一个文件处理器（文件服务器），本地目录提供静态文件
		handler = http.FileServer(http.Dir(fmt.Sprintf("%s/src/embed/static", model.GetProjectDir())))
	} else {
		// 嵌入的静态文件
		ioFs, err := fs.Sub(embedFs, "embed/static")
		if err != nil {
			panic(err)
		}
		// 创建一个文件处理器（文件服务器），FS提供静态文件
		handler = http.FileServer(http.FS(ioFs))
	}

	// 注册路由到处理器
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", handler))
}
