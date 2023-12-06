// @author xiangqian
// @date 19:45 2023/12/04
package common

import (
	"fmt"
	"note/src/model"
	util_string "note/src/util/string"
)

// 数据目录
var dataDir = model.Ini.Data.Dir

// html 404
func html404(err any) (string, model.Response) {
	return "404", model.Response{Msg: util_string.String(err)}
}

// 重定向到列表
func redirectList(table string, err any) (string, model.Response) {
	return fmt.Sprintf("redirect:/%s", table), model.Response{Msg: util_string.String(err)}
}

// 重定向到详情页
func redirectView(table string, id int64, err any) (string, model.Response) {
	return fmt.Sprintf("redirect:/%s/%d/view", table, id), model.Response{Msg: util_string.String(err)}
}
