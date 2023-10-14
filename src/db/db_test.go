// db test
// @author xiangqian
// @date 20:47 2023/06/12
package db

import (
	"gorm.io/gorm"
	"log"
	"note/src/typ"
	"testing"
)

var db *gorm.DB

func init() {
	dsn := "C:\\Users\\xiangqian\\Desktop\\tmp\\note\\data\\database.db"
	log.Println("db_test init", dsn)
	var err error
	db, err = Db(dsn)
	if err != nil {
		panic(err)
	}
}

func TestDb(t *testing.T) {
	log.Println(db)
}

func TestRaw1(t *testing.T) {
	result, err := Raw[int](db, "select 10+10")
	if err != nil {
		panic(err)
	}
	log.Println(result)
}

func TestRaw2(t *testing.T) {
	user, err := Raw[typ.User](db, "SELECT `id`, `name`, `nickname`, `passwd`, `rem`, `try`, `del`, `add_time`, `upd_time` FROM `user` LIMIT 1")
	if err != nil {
		panic(err)
	}
	log.Println(user)
}

func TestRaw3(t *testing.T) {
	users, err := Raw[[]typ.User](db, "SELECT `id`, `name`, `nickname`, `passwd`, `rem`, `try`, `del`, `add_time`, `upd_time` FROM `user` LIMIT 10")
	if err != nil {
		panic(err)
	}
	log.Println(users)
}

func TestPage(t *testing.T) {
	page, err := Page[typ.User](db, 1, 10, "SELECT `id`, `name`, `nickname`, `passwd`, `rem`, `try`, `del`, `add_time`, `upd_time` FROM `user`")
	if err != nil {
		panic(err)
	}
	log.Println(page)
	if page.Total > 0 {
		log.Println(len(page.Data))
	}
}
