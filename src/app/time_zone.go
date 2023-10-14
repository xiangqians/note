// 时区
// @author xiangqian
// @date 19:49 2023/07/10
package app

import (
	"log"
	"time"
)

// 初始化时区
func initTimeZone() {
	timeZone := arg.TimeZone
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		log.Println(timeZone, err)

		// GoLang time 包默认是UTC
		loc = time.UTC
	}

	// 设置时区
	time.Local = loc

	log.Printf("TimeZone %s\n", loc)
}
