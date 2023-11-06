// @author xiangqian
// @date 20:22 2023/06/10
package main

import (
	// 初始化日志记录器
	_ "note/src/log"

	// 初始化系统
	_ "note/src/sys"

	"note/src/server"
)

func main() {
	// 启动应用
	server.Run()
	//log.Println("夫天地者，万物之逆旅也；光阴者，百代之过客也。而浮生若梦，为欢几何？古人秉烛夜游，良有以也。况阳春召我以烟景，大块假我以文章。会桃花之芳园，序天伦之乐事。群季俊秀，皆为惠连；吾人咏歌，独惭康乐。幽赏未已，高谈转清。开琼筵以坐花，飞羽觞而醉月。不有佳咏，何伸雅怀？如诗不成，罚依金谷酒数。")
}
