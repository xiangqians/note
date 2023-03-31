// db test
// @author xiangqian
// @date 18:33 2023/03/30
package test

import (
	"fmt"
	_db "note/src/db"
	"note/src/typ/api"
	typ_page "note/src/typ/page"
	"testing"
)

func TestDb(t *testing.T) {

	db := _db.Get("C:\\Users\\xiangqian\\Desktop\\tmp\\note\\data\\database.db")
	defer db.Close()

	err := db.Open()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = db.Begin()
	if err != nil {
		fmt.Println(err)
		return
	}

	//var user api.User
	//var user []api.User
	rows, err := db.Qry("SELECT `id`, `name`, `nickname`, `rem`, `add_time`, `upd_time` FROM `user` LIMIT 1")
	//rows, err := db.Qry("SELECT `add_time` FROM `user` union all SELECT `id` FROM `user` union all SELECT `upd_time` FROM `user`")
	//users, _, _ := _db.RowsMapper[api.User](rows)
	users, _, _ := _db.RowsMapper[[]api.User](rows, err)
	//users, _, _ := _db.RowsMapper[[]int64](rows)
	//users, _, _ := _db.RowsMapper[int64](rows)
	//users, _, _ := _db.RowsMapper[map[string]any](rows)
	//fmt.Println("i", i)
	//user = i.(api.User)
	fmt.Println("users", users)

	page, err := _db.Page[api.User](db, typ_page.PageReq{Current: 2, Size: 10}, "SELECT `id`, `name`, `nickname`, `rem`, `add_time`, `upd_time` FROM `user`")
	fmt.Println("page", page)

	db.Commit()
}
