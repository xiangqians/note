// local
// @author xiangqian
// @date 22:59 2023/02/14
package app

import (
	"log"
	"time"
)

// 时区
func loc() {
	loc, err := time.LoadLocation(appArg.Loc)
	if err != nil {
		log.Println(appArg.Loc, err)

		// GoLang time 包默认是UTC
		loc = time.UTC
	}

	// set loc
	time.Local = loc

	log.Printf("set loc: %s\n", loc)
}
