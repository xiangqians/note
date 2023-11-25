// @author xiangqian
// @date 22:33 2023/11/07
package handler

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"note/src/handler/index"
	"note/src/model"
	"note/src/session"
	"note/src/util/i18n"
	util_os "note/src/util/os"
	util_time "note/src/util/time"
	"strconv"
	"strings"
)

// Handle 注册路由和相应的处理器函数
func Handle(staticFs, templateFs embed.FS, mux *http.ServeMux) {
	// 处理静态资源
	handleStatic(staticFs, mux)

	// 处理模板
	handleTemplate(templateFs, mux)
}

// 处理模板
func handleTemplate(templateFs embed.FS, mux *http.ServeMux) {
	// 当前语言
	currentLanguage := i18n.ZH

	// 模板函数
	templateFuncMap := template.FuncMap{
		// i18n国际化
		"Localize": func(name string) string {
			return i18n.GetMessage(name, currentLanguage)
		},

		// 两数相加
		"Add": func(arg1 any, arg2 any) int64 {
			i1, _ := strconv.ParseInt(fmt.Sprintf("%v", arg1), 10, 64)
			i2, _ := strconv.ParseInt(fmt.Sprintf("%v", arg2), 10, 64)
			return i1 + i2
		},

		"NowUnix": func() int64 {
			return util_time.NowUnix()
		},
		"HumanizUnix": func(unix int64) string {
			return util_time.HumanizUnix(unix)
		},
		"HumanizFileSize": func(size int64) string {
			return util_os.HumanizFileSize(size)
		},
	}

	handle := func(pattern string, handler func(request *http.Request, session *session.Session) (name string, response model.Response)) {
		// 注册路由和相应的处理器函数
		mux.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
			// 获取会话
			session := session.GetSession(writer, request)

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

			// 设置当前语言（多线程会存在问题，但，笔记应用没有什么并发量，可忽略）
			currentLanguage = language

			// 处理器
			name, response := handler(request, session)

			// 判断是否进行重定向
			if strings.HasPrefix(name, "redirect:") {
				// 301 永久重定向
				http.Redirect(writer, request, name[len("redirect:"):], http.StatusMovedPermanently)
				return
			}

			// 解析模板
			template, err := template.New(fmt.Sprintf("%s.html", name)).Funcs(templateFuncMap).ParseFS(templateFs,
				fmt.Sprintf("embed/template/%s.html", name), // 主模板文件
				"embed/template/common/foot1.html",          // 嵌套模板文件
				"embed/template/common/foot2.html",          // 嵌套模板文件
				"embed/template/common/table.html",          // 嵌套模板文件
				"embed/template/common/variable.html") // 嵌套模板文件
			if err != nil {
				panic(err)
			}

			// 执行模板
			err = template.Execute(writer, map[string]any{
				"contextPath": model.Ini.Server.ContextPath, // 上下文路径
				"path":        request.URL.Path,             // 请求路径
				"uri":         request.RequestURI,           // 请求uri地址
				"user":        model.User{},                 // 当前登录用户信息
				"response":    response,                     // 响应数据
			})
			if err != nil {
				panic(err)
			}
		})
	}

	// user

	// index
	handle("/", index.Index)
}

// 处理静态资源
func handleStatic(embedFs embed.FS, mux *http.ServeMux) {
	// 嵌入的静态文件
	ioFs, err := fs.Sub(embedFs, "embed/static")
	if err != nil {
		panic(err)
	}

	// 创建一个文件服务器处理器，将静态文件从嵌入的FS中提供
	handler := http.FileServer(http.FS(ioFs))

	// 注册路由到处理器
	mux.Handle("/static/", http.StripPrefix("/static/", handler))
}
