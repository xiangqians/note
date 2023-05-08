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
	l, err := time.LoadLocation(appArg.Loc)
	if err != nil {
		log.Println(appArg.Loc, err)

		// GoLang time 包默认是UTC
		l = time.UTC
	}

	// set loc
	time.Local = l

	log.Printf("set loc: %s\n", l)
}
