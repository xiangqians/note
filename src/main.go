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
)

//go:embed embed/i18n
var i18nFs embed.FS

//go:embed embed/static
var staticFs embed.FS

//go:embed embed/template
var templateFs embed.FS

func main() {
	// 创建了一个 http.ServeMux 对象，用于注册和管理路由和处理器函数
	mux := http.NewServeMux()

	// 嵌入资源
	handler.I18nFs = i18nFs
	handler.StaticFs = staticFs
	handler.TemplateFs = templateFs

	// 注册路由和相应的处理器函数
	handler.Handle(mux)

	// 服务监听端口
	port := model.Ini.Server.Port

	// 配置服务
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	// 监听
	log.Printf("Server started on port %d\n", port)
	panic(server.ListenAndServe())
}
