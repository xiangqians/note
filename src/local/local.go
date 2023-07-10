// local
// @author xiangqian
// @date 19:49 2023/07/10
package local

import (
	"log"
	"note/src/arg"
	"time"
)

// Init 初始化时区
func Init() {
	loc, err := time.LoadLocation(arg.Arg.Loc)
	if err != nil {
		log.Println(arg.Arg.Loc, err)

		// GoLang time 包默认是UTC
		loc = time.UTC
	}

	// set loc
	time.Local = loc

	log.Printf("loc: %s\n", loc)
}
