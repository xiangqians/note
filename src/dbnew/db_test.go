// @author xiangqian
// @date 20:47 2023/06/12
package dbnew

import (
	"log"
	"note/src/model"
	util_json "note/src/util/json"
	"testing"
	"time"
)

var dbConnPool DbConnPool

func init() {
	dbConnPool = &GormDbConnPool{}
}

func GetDb() Db {
	dsn := "C:\\Users\\xiangqian\\Desktop\\tmp\\note\\data\\database.db"
	db, err := dbConnPool.Get(dsn)
	if err != nil {
		panic(err)
	}
	return db
}

func TestDb(t *testing.T) {
	flag := false
	for i := 0; i < 10; i++ {
		go func(i int) {
			db := GetDb()
			log.Println(i, db)
			if i == 9 {
				flag = true
			}
		}(i)
	}

	for {
		if flag {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func TestAdd(t *testing.T) {
	db := GetDb()
	rowsAffected, insertId, err := db.Add("INSERT INTO `user` (`name`, `nickname`, `passwd`, `rem`) VALUES (?, ?, ?, ?)", "test", "测试", "passwd", "备注")
	if err != nil {
		panic(err)
	}
	log.Println("rowsAffected", rowsAffected)
	log.Println("insertId", insertId)
}

func TestUpd(t *testing.T) {
	db := GetDb()
	rowsAffected, err := db.Upd("UPDATE `user` SET `nickname` = ? Where `name` = ?", "测试2", "test")
	if err != nil {
		panic(err)
	}
	log.Println("rowsAffected", rowsAffected)
}

func TestGet1(t *testing.T) {
	db := GetDb()
	result, err := db.Get("SELECT 10+10")
	if err != nil {
		panic(err)
	}
	var i int8
	result.Scan(&i)
	log.Println(i)
}

func TestGet2(t *testing.T) {
	db := GetDb()

	// count
	result, err := db.Get("SELECT COUNT(1) FROM `user`")
	if err != nil {
		panic(err)
	}
	var count int64
	result.Scan(&count)
	log.Println("count", count)

	// name
	result, err = db.Get("SELECT `name` FROM `user` LIMIT 1")
	if err != nil {
		panic(err)
	}
	var name int64
	result.Scan(&name)
	log.Println("name", name)
}

func TestGet3(t *testing.T) {
	db := GetDb()
	result, err := db.Get("SELECT `id`, `name`, `nickname`, `passwd`, `rem`, `try`, `del`, `add_time`, `upd_time` FROM `user` LIMIT 1")
	if err != nil {
		panic(err)
	}
	var user model.User
	result.Scan(&user)

	json, err := util_json.Serialize(user, true)
	if err != nil {
		panic(err)
	}
	log.Println("\n", json)
}

func TestGet4(t *testing.T) {
	db := GetDb()
	result, err := db.Get("SELECT `id`, `name`, `nickname`, `passwd`, `rem`, `try`, `del`, `add_time`, `upd_time` FROM `user` LIMIT 10")
	if err != nil {
		panic(err)
	}
	var users []model.User
	result.Scan(&users)

	json, err := util_json.Serialize(users, true)
	if err != nil {
		panic(err)
	}
	log.Println("\n", json)
}

func TestPage(t *testing.T) {
	db := GetDb()
	result, err := db.Page("SELECT `id`, `name`, `nickname`, `passwd`, `rem`, `try`, `del`, `add_time`, `upd_time` FROM `user`", 1, 2)
	if err != nil {
		panic(err)
	}
	var users []model.User
	result.Scan(&users)
	log.Println("count", result.Count())

	json, err := util_json.Serialize(users, true)
	if err != nil {
		panic(err)
	}
	log.Println("\n", json)
}
