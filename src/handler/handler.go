// @author xiangqian
// @date 22:33 2023/11/07
package handler

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

var I18nFs embed.FS
var StaticFs embed.FS
var TemplateFs embed.FS

// Handle 注册路由和相应的处理器函数
func Handle(mux *http.ServeMux) {
	// 处理静态资源
	// 嵌入的静态文件
	staticFs, err := fs.Sub(StaticFs, "embed/static")
	if err != nil {
		panic(err)
	}
	// 创建一个文件服务器，将静态文件从嵌入的FS中提供
	fileServer := http.FileServer(http.FS(staticFs))
	// 注册路由到处理器
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// 注册路由和相应的处理器函数
	//mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/about", aboutHandler)
	mux.HandleFunc("/contact", contactHandler)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page!")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the about page.")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Contact us at example@example.com")
}
