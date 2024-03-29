// @author xiangqian
// @date 22:33 2023/11/07
package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	pkg_template "html/template"
	"io/fs"
	"log"
	"net/http"
	"note/src/embed"
	"note/src/handler/common"
	"note/src/handler/index"
	"note/src/handler/note"
	"note/src/handler/system"
	"note/src/model"
	"note/src/session"
	"note/src/util/i18n"
	util_json "note/src/util/json"
	util_os "note/src/util/os"
	util_time "note/src/util/time"
	"os"
	"strings"
)

// 模式，dev、test、prod
var mode = os.Getenv("NOTE_MODE")

// 项目目录
var projectDir, _ = os.Getwd()

// Handle 注册路由到处理器
func Handle(router *mux.Router) {
	// 处理静态资源
	handleStatic(router)

	// 处理模板
	handleTemplate(router)
}

// 处理模板
func handleTemplate(router *mux.Router) {
	// 模板函数
	templateFuncMap := pkg_template.FuncMap{
		// 判断一个字符串是否包含另一个字符串
		"contains": func(s, substr string) bool {
			return strings.Contains(s, substr)
		},

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

		"Serialize": func(v any) string {
			json, err := util_json.Serialize(v, false)
			if err != nil {
				log.Println(v, err)
			}
			return json
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

	handlerFunc := func(handler func(request *http.Request, writer http.ResponseWriter, session *session.Session) (name string, response model.Response)) http.HandlerFunc {
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

			triggerCalculateFolderSizeTask := func() {
				if request.Method == http.MethodPost && strings.HasPrefix(path, contextPath+"/note") {
					note.TriggerCalculateFolderSizeTask()
				}
			}

			name, response := handler(request, writer, session)
			if name == "" {
				triggerCalculateFolderSizeTask()
				return
			}

			// 判断是否进行重定向
			if strings.HasPrefix(name, "redirect:") {
				// 保存消息到会话
				msg := response.Msg
				if msg != "" {
					session.SetMsg(msg)
				}

				// 301 永久重定向
				url := contextPath + name[len("redirect:"):]
				if strings.Contains(url, "?") {
					url += fmt.Sprintf("&t=%d", util_time.NowUnix())
				} else {
					url += fmt.Sprintf("?t=%d", util_time.NowUnix())
				}
				http.Redirect(writer, request, url, http.StatusMovedPermanently)
				triggerCalculateFolderSizeTask()
				return
			}

			var template *pkg_template.Template
			var err error
			if mode == "dev" {
				template, err = pkg_template.New(name).Funcs(templateFuncMap).ParseFiles(fmt.Sprintf("%s/src/embed/template/%s.html", projectDir, name),
					fmt.Sprintf("%s/src/embed/template/common/foot1.html", projectDir),
					fmt.Sprintf("%s/src/embed/template/common/foot2.html", projectDir),
					fmt.Sprintf("%s/src/embed/template/common/table.html", projectDir),
					fmt.Sprintf("%s/src/embed/template/common/variable.html", projectDir),
					fmt.Sprintf("%s/src/embed/template/common/float.html", projectDir))
			} else {
				template, err = pkg_template.New(name).Funcs(templateFuncMap).ParseFS(embed.TemplateFs,
					fmt.Sprintf("embed/template/%s.html", name), // 主模板文件
					"embed/template/common/foot1.html",          // 嵌套模板文件
					"embed/template/common/foot2.html",          // 嵌套模板文件
					"embed/template/common/table.html",          // 嵌套模板文件
					"embed/template/common/variable.html",       // 嵌套模板文件
					"embed/template/common/float.html")          // 嵌套模板文件
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
			triggerCalculateFolderSizeTask()
		}
	}

	imageHandlerFunc := func(handler func(request *http.Request, writer http.ResponseWriter, session *session.Session, table string) (name string, response model.Response)) http.HandlerFunc {
		return handlerFunc(func(request *http.Request, writer http.ResponseWriter, session *session.Session) (name string, response model.Response) {
			return handler(request, writer, session, common.TableImage)
		})
	}

	audioHandlerFunc := func(handler func(request *http.Request, writer http.ResponseWriter, session *session.Session, table string) (name string, response model.Response)) http.HandlerFunc {
		return handlerFunc(func(request *http.Request, writer http.ResponseWriter, session *session.Session) (name string, response model.Response) {
			return handler(request, writer, session, common.TableAudio)
		})
	}

	videoHandlerFunc := func(handler func(request *http.Request, writer http.ResponseWriter, session *session.Session, table string) (name string, response model.Response)) http.HandlerFunc {
		return handlerFunc(func(request *http.Request, writer http.ResponseWriter, session *session.Session) (name string, response model.Response) {
			return handler(request, writer, session, common.TableVideo)
		})
	}

	noteHandlerFunc := func(handler func(request *http.Request, writer http.ResponseWriter, session *session.Session, table string) (name string, response model.Response)) http.HandlerFunc {
		return handlerFunc(func(request *http.Request, writer http.ResponseWriter, session *session.Session) (name string, response model.Response) {
			return handler(request, writer, session, common.TableNote)
		})
	}

	router.NotFoundHandler = handlerFunc(func(request *http.Request, writer http.ResponseWriter, session *session.Session) (name string, response model.Response) {
		return common.NotFound(request, writer, session, nil)
	})

	// system
	router.HandleFunc(contextPath+"/signin", handlerFunc(system.SignIn))
	router.HandleFunc(contextPath+"/signout", handlerFunc(system.SignOut))
	router.HandleFunc(contextPath+"/setting", handlerFunc(system.Setting))

	// index
	router.HandleFunc(contextPath+"/", handlerFunc(index.Index))

	// image
	router.HandleFunc(contextPath+"/image", imageHandlerFunc(common.List))
	router.HandleFunc(contextPath+"/image/upload", imageHandlerFunc(common.Upload)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/image/{id:[0-9]+}/reupload", imageHandlerFunc(common.ReUpload)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/image/{id:[0-9]+}/rename", imageHandlerFunc(common.Rename)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/image/{id:[0-9]+}/del", imageHandlerFunc(common.Del)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/image/{id:[0-9]+}/restore", imageHandlerFunc(common.Restore)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/image/{id:[0-9]+}/permlydel", imageHandlerFunc(common.PermlyDel)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/image/{id:[0-9]+}", imageHandlerFunc(common.Get))
	router.HandleFunc(contextPath+"/image/{id:[0-9]+}/view", imageHandlerFunc(common.View))

	// audio
	router.HandleFunc(contextPath+"/audio", audioHandlerFunc(common.List))
	router.HandleFunc(contextPath+"/audio/upload", audioHandlerFunc(common.Upload)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/audio/{id:[0-9]+}/reupload", audioHandlerFunc(common.ReUpload)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/audio/{id:[0-9]+}/rename", audioHandlerFunc(common.Rename)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/audio/{id:[0-9]+}/del", audioHandlerFunc(common.Del)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/audio/{id:[0-9]+}/restore", audioHandlerFunc(common.Restore)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/audio/{id:[0-9]+}/permlydel", audioHandlerFunc(common.PermlyDel)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/audio/{id:[0-9]+}", audioHandlerFunc(common.Get))
	router.HandleFunc(contextPath+"/audio/{id:[0-9]+}/view", audioHandlerFunc(common.View))

	// video
	router.HandleFunc(contextPath+"/video", videoHandlerFunc(common.List))
	router.HandleFunc(contextPath+"/video/upload", videoHandlerFunc(common.Upload)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/video/{id:[0-9]+}/reupload", videoHandlerFunc(common.ReUpload)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/video/{id:[0-9]+}/rename", videoHandlerFunc(common.Rename)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/video/{id:[0-9]+}/del", videoHandlerFunc(common.Del)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/video/{id:[0-9]+}/restore", videoHandlerFunc(common.Restore)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/video/{id:[0-9]+}/permlydel", videoHandlerFunc(common.PermlyDel)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/video/{id:[0-9]+}", videoHandlerFunc(common.Get))
	router.HandleFunc(contextPath+"/video/{id:[0-9]+}/view", videoHandlerFunc(common.View))

	// note
	router.HandleFunc(contextPath+"/note", noteHandlerFunc(common.List))
	router.HandleFunc(contextPath+"/note/{pid:[0-9]+}/list", noteHandlerFunc(common.List))
	router.HandleFunc(contextPath+"/note/addfolder", handlerFunc(note.AddFolder)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/note/addmdfile", handlerFunc(note.AddMdFile)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/note/upload", noteHandlerFunc(common.Upload)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/note/{id:[0-9]+}/reupload", noteHandlerFunc(common.ReUpload)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/note/{id:[0-9]+}/rename", noteHandlerFunc(common.Rename)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/note/paste", handlerFunc(note.Paste)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/note/{id:[0-9]+}/del", noteHandlerFunc(common.Del)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/note/{id:[0-9]+}/restore", noteHandlerFunc(common.Restore)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/note/{id:[0-9]+}/permlydel", noteHandlerFunc(common.PermlyDel)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/note/{id:[0-9]+}", noteHandlerFunc(common.Get)).Methods(http.MethodGet)
	router.HandleFunc(contextPath+"/note/{id:[0-9]+}", handlerFunc(note.Upd)).Methods(http.MethodPost)
	router.HandleFunc(contextPath+"/note/{id:[0-9]+}/view", noteHandlerFunc(common.View))
}

// 处理静态资源
func handleStatic(router *mux.Router) {
	var handler http.Handler
	if mode == "dev" {
		// 创建一个文件处理器（文件服务器），本地目录提供静态文件
		handler = http.FileServer(http.Dir(fmt.Sprintf("%s/src/embed/static", projectDir)))
	} else {
		// 嵌入的静态文件
		ioFs, err := fs.Sub(embed.StaticFs, "embed/static")
		if err != nil {
			panic(err)
		}
		// 创建一个文件处理器（文件服务器），FS提供静态文件
		handler = http.FileServer(http.FS(ioFs))
	}

	// 上下文路径
	contextPath := model.Ini.Server.ContextPath

	// 注册路由到处理器
	router.PathPrefix(contextPath + "/static/").Handler(http.StripPrefix(contextPath+"/static/", handler))
}
