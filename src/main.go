// main
// @author xiangqian
// @date 20:22 2023/06/10
package main

import (
	"note/src/api"
	"note/src/app"
)

func main() {
	// 创建Engine
	engine := app.New()

	// 初始API
	api.Init(engine)

	// 运行Engine
	app.Run(engine)
}
