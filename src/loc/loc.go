// local
// @author xiangqian
// @date 22:59 2023/02/14
package loc

import (
	"log"
	"note/src/arg"
	"time"
)

// Init 初始化时区
func Init() {
	arg := arg.Get()
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
