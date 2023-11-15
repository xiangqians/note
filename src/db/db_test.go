// @author xiangqian
// @date 20:47 2023/06/12
package db

import (
	"log"
	"note/src/structure"
	util_json "note/src/util/json"
	"sync"
	"testing"
	"time"
)

func TestDb(t *testing.T) {
	var waitGroup sync.WaitGroup
	for i := 0; i < 10; i++ {
		// 添加任务数
		waitGroup.Add(1)
		go func(i int) {
			db := Get()
			log.Println(i, db)
			// 完成任务
			waitGroup.Done()
		}(i)
	}
	// 阻塞等待所有任务完成
	waitGroup.Wait()
}

func TestVersion(t *testing.T) {
	db := Get()

	var sql string
	switch structure.Ini.Db.Driver {
	case "sqlite", "sqlite3":
		sql = "SELECT SQLITE_VERSION()"
	}

	result, err := db.Get(sql)
	if err != nil {
		panic(err)
	}

	var version string
	result.Scan(&version)
	log.Println(version)
}

func TestAdd(t *testing.T) {
	db := Get()
	rowsAffected, insertId, err := db.Add("INSERT INTO `user` (`name`, `nickname`, `passwd`, `rem`) VALUES (?, ?, ?, ?)", "test", "测试", "passwd", "备注")
	if err != nil {
		panic(err)
	}
	log.Println("rowsAffected", rowsAffected)
	log.Println("insertId", insertId)
}

func TestUpd(t *testing.T) {
	db := Get()
	rowsAffected, err := db.Upd("UPDATE `user` SET `nickname` = ? Where `name` = ?", "测试2", "test")
	if err != nil {
		panic(err)
	}
	log.Println("rowsAffected", rowsAffected)
}

func TestGetSum(t *testing.T) {
	var waitGroup sync.WaitGroup
	for i := 0; i < 10; i++ {
		// 添加任务数
		waitGroup.Add(1)
		go func() {
			db := Get()
			result, err := db.Get("SELECT 10+10")
			if err != nil {
				panic(err)
			}
			var i int
			result.Scan(&i)
			log.Println(i)
			// 完成任务
			waitGroup.Done()
		}()
	}
	// 阻塞等待所有任务完成
	waitGroup.Wait()
}

func TestGetField(t *testing.T) {
	db := Get()

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

func TestGetStruct(t *testing.T) {
	db := Get()
	result, err := db.Get("SELECT `id`, `name`, `nickname`, `passwd`, `rem`, `try`, `del`, `add_time`, `upd_time` FROM `user` LIMIT 1")
	if err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Second)
	var user structure.User
	result.Scan(&user)

	json, err := util_json.Serialize(user, true)
	if err != nil {
		panic(err)
	}
	log.Println("\n", json)
}

func TestConcurrentGetStruct(t *testing.T) {
	var waitGroup sync.WaitGroup
	for i := 0; i < 10; i++ {
		// 添加任务数
		waitGroup.Add(1)
		go func() {
			TestGetStruct(t)
			// 完成任务
			waitGroup.Done()
		}()
	}
	// 阻塞等待所有任务完成
	waitGroup.Wait()
	log.Println(Stats())
}

func TestGetStructSlice(t *testing.T) {
	db := Get()
	result, err := db.Get("SELECT `id`, `name`, `nickname`, `passwd`, `rem`, `try`, `del`, `add_time`, `upd_time` FROM `user` LIMIT 10")
	if err != nil {
		panic(err)
	}
	var users []structure.User
	result.Scan(&users)

	json, err := util_json.Serialize(users, true)
	if err != nil {
		panic(err)
	}
	log.Println("\n", json)
}

func TestPage(t *testing.T) {
	page := structure.Page{
		Current: 1,
		Size:    2,
	}

	db := Get()
	result, err := db.Page("SELECT `id`, `name`, `nickname`, `passwd`, `rem`, `try`, `del`, `add_time`, `upd_time` FROM `user`", page.Current, page.Size)
	if err != nil {
		panic(err)
	}

	page.Total = result.Count()
	page.InitIndexes()

	var users []structure.User
	result.Scan(&users)
	page.Data = users

	json, err := util_json.Serialize(page, true)
	if err != nil {
		panic(err)
	}
	log.Println("\n", json)
}
