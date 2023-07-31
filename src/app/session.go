// session
// @author xiangqian
// @date 20:01 2023/03/22
package app

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	util_crypto_bcrypt "note/src/util/crypto/bcrypt"
)

// 初始化session
func initSession(engine *gin.Engine) {
	// 密钥
	passwd := "$2a$10$NkWzRTyz1ZNnNfjLmxreaeZ31DCiwCEWJlXJAVDkG8fD9Ble2mg4K"
	hash, err := util_crypto_bcrypt.Generate(passwd)
	if err != nil {
		log.Println(err)
		hash = passwd
	}
	keyPairs := []byte(hash)[:32]

	// session存储引擎支持：基于内存、redis、mysql等
	// 1、创建基于cookie的存储引擎
	//store := cookie.NewStore(keyPairs)
	// 2、创建基于mem（内存）的存储引擎，其实就是一个 map[interface]interface 对象
	//store := memstore.NewStore(keyPairs)
	store := NewStore(keyPairs)

	// store配置
	store.Options(sessions.Options{
		//Secure: true,
		//SameSite: http.SameSiteNoneMode,
		Path:   "/",
		MaxAge: 60 * 60 * 12, // 设置session过期时间，单位：s，12h
	})

	// 设置session中间件
	engine.Use(sessions.Sessions("NoteSessionId", // session & cookie 名称
		store))
}
