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
	util_os "note/src/util/os"
	util_time "note/src/util/time"
	"strconv"
)

// Handle 注册路由和相应的处理器函数
func Handle(i18nFs, staticFs, templateFs embed.FS, mux *http.ServeMux) {
	// 处理静态资源
	handleStatic(staticFs, mux)

	// 处理模板
	handleTemplate(templateFs, mux)
}

// 处理静态资源
func handleStatic(staticFs embed.FS, mux *http.ServeMux) {
	// 嵌入的静态文件
	fs, err := fs.Sub(staticFs, "embed/static")
	if err != nil {
		panic(err)
	}

	// 创建一个文件服务器，将静态文件从嵌入的FS中提供
	fileServer := http.FileServer(http.FS(fs))

	// 注册路由到处理器
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
}

// 处理模板
func handleTemplate(templateFs embed.FS, mux *http.ServeMux) {
	// 模板函数
	templateFuncMap := template.FuncMap{
		// 获取i18n文件中key对应的value
		"Localize": func(key string) string {
			return key
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

	handle := func(pattern string, handler func(r *http.Request) (templateName string, response model.Response)) {
		// 注册路由和相应的处理器函数
		mux.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
			log.Println("--> ", request.URL.Path)

			// 解析嵌入的父模板
			parentTmpl, err := template.New("index.html").Funcs(templateFuncMap).ParseFS(templateFs, "embed/template/index.html",
				"embed/template/common/foot1.html",
				"embed/template/common/foot2.html",
				"embed/template/common/table.html",
				"embed/template/common/variable.html")
			if err != nil {
				panic(err)
			}

			err = parentTmpl.Execute(writer, map[string]string{"title": "Golang Embed 测试"})
			if err != nil {
				panic(err)
			}
		})
	}

	handle("/", index.Index)
}
