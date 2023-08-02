// db test
// @author xiangqian
// @date 20:47 2023/06/12
package db

import (
	"log"
	"note/src/typ"
	"testing"
)

func TestDb(t *testing.T) {
	dsn := "C:\\Users\\xiangqian\\Desktop\\tmp\\note\\data\\database.db"
	db, err := Db(dsn)
	if err != nil {
		panic(err)
	}

	var result int64
	db.Raw("select 1+10").Take(&result)
	log.Println(result)

	var user typ.User
	db.Raw("SELECT `id`, `name`, `nickname`, `passwd`, `rem`, `try`, `add_time`, `upd_time` FROM `user` LIMIT 1").Scan(&user)
	log.Println(user)

	// len 0, cap ?
	users := make([]typ.User, 0, 1)
	db.Raw("SELECT `id`, `name`, `nickname`, `passwd`, `rem`, `try`, `add_time`, `upd_time` FROM `user`").Scan(&users)
	log.Println(len(users), users)
}
