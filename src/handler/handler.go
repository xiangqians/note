// @author xiangqian
// @date 22:33 2023/11/07
package handler

import (
	"fmt"
	"net/http"
)

// Handle 注册路由和相应的处理器函数
func Handle(mux *http.ServeMux) {
	// 静态文件服务器
	fileServer := http.FileServer(http.Dir("./res/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// 注册路由和相应的处理器函数
	mux.HandleFunc("/", homeHandler)
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
