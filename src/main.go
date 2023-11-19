// @author xiangqian
// @date 20:22 2023/06/10
package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"note/src/handler"
	_ "note/src/log"
	"note/src/model"
	util_i18n "note/src/util/i18n"
)

// 嵌入资源
// embed.FS

//go:embed embed/i18n
var i18nFs embed.FS

//go:embed embed/static
var staticFs embed.FS

//go:embed embed/template
var templateFs embed.FS

func main() {
	// 初始化i18n
	util_i18n.Init(i18nFs)

	// 创建了一个 http.ServeMux 对象，用于注册和管理路由和处理器函数
	mux := http.NewServeMux()

	// 注册路由和相应的处理器函数
	handler.Handle(staticFs, templateFs, mux)

	// 服务监听端口
	port := model.Ini.Server.Port

	// 配置服务
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
		//ReadTimeout:    60 * time.Second, // 设置HTTP服务器读取超时。服务器在读取请求主体时的超时时间，当超过指定的时间后，如果请求主体还未完全读取，服务器将关闭连接
		//WriteTimeout:   60 * time.Second, // 设置HTTP服务器写入超时。服务器在写入响应主体时的超时时间，当超过指定的时间后，如果响应主体还未完全写入，服务器将关闭连接
		MaxHeaderBytes: 1 << 20, // 设置接收的HTTP请求头的最大字节数。默认情况下，MaxHeaderBytes 的值为 1 << 20，即 1MB（1024 * 1024 Byte）。这意味着如果请求头的大小超过了 1MB，服务器将返回一个错误响应 http.ErrHeaderTooLong
	}

	// 监听
	log.Printf("Server started on port %d\n", port)
	panic(server.ListenAndServe())
}
