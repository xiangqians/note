// @author xiangqian
// @date 20:47 2023/06/12
package db

import (
	"log"
	"note/src/model"
	util_json "note/src/util/json"
	util_time "note/src/util/time"
	"sync"
	"testing"
	"time"
)

func TestVersion(t *testing.T) {
	db := Get()

	var sql string
	switch model.Ini.Db.Driver {
	case "sqlite", "sqlite3":
		sql = "SELECT SQLITE_VERSION()"
	default:
		panic("未知数据库类型")
	}

	result, err := db.Get(sql)
	if err != nil {
		panic(err)
	}

	var version string
	result.Scan(&version)
	log.Println(version)
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
			time.Sleep(1 * time.Second)

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

func TestAdd(t *testing.T) {
	db := Get()
	rowsAffected, insertId, err := db.Add("INSERT INTO `image` (`name`, `type`, `size`, `add_time`) VALUES (?, ?, ?, ?)", "test", "png", 10232, util_time.NowUnix())
	if err != nil {
		panic(err)
	}
	log.Println("rowsAffected", rowsAffected)
	log.Println("insertId", insertId)
}

func TestUpd(t *testing.T) {
	db := Get()
	rowsAffected, err := db.Upd("UPDATE `image` SET `name` = ? Where `id` = ?", "6", 6)
	if err != nil {
		panic(err)
	}
	log.Println("rowsAffected", rowsAffected)
}

func TestGetField(t *testing.T) {
	db := Get()

	// count
	result, err := db.Get("SELECT COUNT(1) FROM `image`")
	if err != nil {
		panic(err)
	}
	var count int64
	result.Scan(&count)
	log.Println("count", count)

	// name
	result, err = db.Get("SELECT `name` FROM `image` LIMIT 1")
	if err != nil {
		panic(err)
	}
	var name int64
	result.Scan(&name)
	log.Println("name", name)
}

func TestGetStruct(t *testing.T) {
	db := Get()
	result, err := db.Get("SELECT `id`, `name`, `type`, `size`, `del`, `add_time`, `upd_time` FROM `image` LIMIT 1")
	if err != nil {
		panic(err)
	}

	//time.Sleep(1 * time.Second)

	var image model.Image
	result.Scan(&image)

	json, err := util_json.Serialize(image, true)
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
	result, err := db.Get("SELECT `id`, `name`, `type`, `size`, `del`, `add_time`, `upd_time` FROM `image` LIMIT 10")
	if err != nil {
		panic(err)
	}

	var images []model.Image
	result.Scan(&images)

	json, err := util_json.Serialize(images, true)
	if err != nil {
		panic(err)
	}
	log.Println("\n", json)
}

func TestPage(t *testing.T) {
	page := model.Page{
		Current: 1,
		Size:    2,
	}

	db := Get()
	result, err := db.Page("SELECT `id`, `name`, `type`, `size`, `del`, `add_time`, `upd_time` FROM `image`", page.Current, page.Size)
	if err != nil {
		panic(err)
	}

	page.Total = result.Count()
	(&page).InitIndexes()

	var images []model.Image
	result.Scan(&images)
	page.Data = images

	json, err := util_json.Serialize(page, true)
	if err != nil {
		panic(err)
	}
	log.Println("\n", json)
}
