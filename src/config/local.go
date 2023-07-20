// config
// @author xiangqian
// @date 19:49 2023/07/10
package config

import (
	"log"
	"time"
)

// 初始化时区
func initLocal() {
	loc, err := time.LoadLocation(arg.Loc)
	if err != nil {
		log.Println(arg.Loc, err)

		// GoLang time 包默认是UTC
		loc = time.UTC
	}

	// set loc
	time.Local = loc

	log.Printf("loc: %s\n", loc)
}
