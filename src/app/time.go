// time
// @author xiangqian
// @date 22:59 2023/02/14
package app

import (
	"log"
	"time"
)

func Local() {
	// GoLang time 包默认是UTC
	time.Local = time.UTC

	// 修改为北京时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	time.Local = loc

	log.Printf("loc: %s\n", time.Local)
}
