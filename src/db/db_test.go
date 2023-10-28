// db test
// @author xiangqian
// @date 20:47 2023/06/12
package db

import (
	"gorm.io/gorm"
	"log"
	"note/src/model"
	util_json "note/src/util/json"
	"testing"
)

var db *gorm.DB

func init() {
	dsn := "C:\\Users\\xiangqian\\Desktop\\tmp\\note\\data\\database.db"
	var err error
	db, err = Db(dsn)
	if err != nil {
		panic(err)
	}
}

func TestDb(t *testing.T) {
	log.Println(db)
}

func TestExec1(t *testing.T) {
	rowsAffected, lastInsertId, err := Exec(db,
		"INSERT INTO `user` (`name`, `nickname`, `passwd`, `rem`) VALUES (?, ?, ?, ?)",
		"test", "测试", "passwd", "备注")
	if err != nil {
		panic(err)
	}
	log.Println("rowsAffected", rowsAffected)
	log.Println("lastInsertId", lastInsertId)
}

func TestExec2(t *testing.T) {
	rowsAffected, lastInsertId, err := Exec(db, "UPDATE `user` SET `nickname` = ? Where `name` = ?", "测试2", "test")
	if err != nil {
		panic(err)
	}
	log.Println("rowsAffected", rowsAffected)
	log.Println("lastInsertId", lastInsertId)
}

func TestRaw1(t *testing.T) {
	result, err := Raw[int](db, "SELECT 10+10")
	if err != nil {
		panic(err)
	}
	log.Println(result)
}

func TestRaw2(t *testing.T) {
	// count
	count, err := Raw[int64](db, "SELECT COUNT(1) FROM `user`")
	if err != nil {
		panic(err)
	}
	log.Println(count)

	// name
	name, err := Raw[int64](db, "SELECT `name` FROM `user` LIMIT 1")
	if err != nil {
		panic(err)
	}
	log.Println(name)
}

func TestRaw3(t *testing.T) {
	user, err := Raw[model.User](db, "SELECT `id`, `name`, `nickname`, `passwd`, `rem`, `try`, `del`, `add_time`, `upd_time` FROM `user` LIMIT 1")
	if err != nil {
		panic(err)
	}

	json, err := util_json.Serialize(user, true)
	if err != nil {
		panic(err)
	}
	log.Println("\n", json)
}

func TestRaw4(t *testing.T) {
	users, err := Raw[[]model.User](db, "SELECT `id`, `name`, `nickname`, `passwd`, `rem`, `try`, `del`, `add_time`, `upd_time` FROM `user` LIMIT 10")
	if err != nil {
		panic(err)
	}

	json, err := util_json.Serialize(users, true)
	if err != nil {
		panic(err)
	}
	log.Println("\n", json)
}

func TestPage(t *testing.T) {
	page, err := Page[model.User](db, 2, 2, "SELECT `id`, `name`, `nickname`, `passwd`, `rem`, `try`, `del`, `add_time`, `upd_time` FROM `user`")
	if err != nil {
		panic(err)
	}

	json, err := util_json.Serialize(page, true)
	if err != nil {
		panic(err)
	}
	log.Println("\n", json)
}
