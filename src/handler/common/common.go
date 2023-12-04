// @author xiangqian
// @date 19:45 2023/12/04
package common

import (
	"fmt"
	"note/src/db"
	"note/src/model"
	util_string "note/src/util/string"
)

// 重定向
func redirect(table string, err any) (string, model.Response) {
	return "redirect:/" + table, model.Response{Msg: util_string.String(err)}
}

// 获取永久删除id，以复用
func getPermlyDelId(table string) (int64, error) {
	db := db.Get()
	result, err := db.Get(fmt.Sprintf("SELECT `id` FROM `%s` WHERE `del` = 2 LIMIT 1", table))
	if err != nil {
		return 0, err
	}

	var id int64
	err = result.Scan(&id)
	return id, err
}
