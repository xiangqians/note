// @author xiangqian
// @date 20:22 2023/06/10
package main

import (
	"fmt"
	"log"
	"net/http"
	"note/src/handler"
	_ "note/src/log"
	"note/src/model"
)

func main() {
	// 创建了一个 http.ServeMux 对象，用于注册和管理路由和处理器函数
	mux := http.NewServeMux()

	// 注册路由和相应的处理器函数
	handler.Handle(mux)

	port := model.Ini.Server.Port
	
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	log.Printf("Server started on port %d\n", port)
	panic(server.ListenAndServe())
}
