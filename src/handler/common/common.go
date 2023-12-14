// @author xiangqian
// @date 19:45 2023/12/04
package common

import (
	"fmt"
	"net/url"
	"note/src/model"
	util_string "note/src/util/string"
	"strings"
)

// 数据目录
var dataDir = model.Ini.Data.Dir

// html 404
func html404(err any) (string, model.Response) {
	return "404", model.Response{Msg: util_string.String(err)}
}

// RedirectList 重定向到列表
func RedirectList(table string, paramMap map[string]any, err any) (string, model.Response) {
	name := fmt.Sprintf("redirect:/%s", table)
	if paramMap != nil {
		// len 0, cap ?
		params := make([]string, 0, len(paramMap))
		for key, value := range paramMap {
			params = append(params, fmt.Sprintf("%s=%s", key, url.QueryEscape(util_string.String(value))))
		}
		name += "?" + strings.Join(params, "&")
	}
	return name, model.Response{Msg: util_string.String(err)}
}

// 重定向到笔记列表
func redirectNoteList(pid int64, err any) (string, model.Response) {
	return "redirect:/note?search=pid%3A%20" + fmt.Sprintf("%d", pid), model.Response{Msg: util_string.String(err)}
}

// 重定向到详情页
func redirectView(table string, id int64, err any) (string, model.Response) {
	return fmt.Sprintf("redirect:/%s/%d/view", table, id), model.Response{Msg: util_string.String(err)}
}
