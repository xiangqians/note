// main
// @author xiangqian
// @date 20:22 2023/06/10
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/api"
	"note/src/config"
)

func main() {
	// 初始化配置
	config.Init()

	// gin模式
	gin.SetMode(gin.DebugMode)

	// default Engine
	engine := gin.Default()

	// 初始化Engine配置
	config.InitEngine(engine)

	// 初始API
	api.Init(engine)

	// run
	addr := fmt.Sprintf(":%d", config.GetArg().Port)
	engine.Run(addr)
}
