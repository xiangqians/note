// local
// @author xiangqian
// @date 22:59 2023/02/14
package app

import (
	"log"
	"time"
)

// 时区
func local() {
	// GoLang time 包默认是UTC
	time.Local = time.UTC

	// 修改为上海时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Println(err)
		return
	}
	time.Local = loc

	log.Printf("local: %s\n", time.Local)
}
