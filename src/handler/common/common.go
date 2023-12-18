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

func NotFound(err any) (string, model.Response) {
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

// RedirectView 重定向到详情页
func RedirectView(table string, id int64, err any) (string, model.Response) {
	return fmt.Sprintf("redirect:/%s/%d/view", table, id), model.Response{Msg: util_string.String(err)}
}

func getPNoteSql() string {
	return "SELECT" +
		"    (CASE WHEN p10.`id` IS NULL THEN '' ELSE '/' || p10.`id` END)" +
		" || (CASE WHEN p9.`id` IS NULL THEN '' ELSE '/' || p9.`id`END)" +
		" || (CASE WHEN p8.`id` IS NULL THEN '' ELSE '/' || p8.`id`END)" +
		" || (CASE WHEN p7.`id` IS NULL THEN '' ELSE '/' || p7.`id`END)" +
		" || (CASE WHEN p6.`id` IS NULL THEN '' ELSE '/' || p6.`id`END)" +
		" || (CASE WHEN p5.`id` IS NULL THEN '' ELSE '/' || p5.`id`END)" +
		" || (CASE WHEN p4.`id` IS NULL THEN '' ELSE '/' || p4.`id`END)" +
		" || (CASE WHEN p3.`id` IS NULL THEN '' ELSE '/' || p3.`id`END)" +
		" || (CASE WHEN p2.`id` IS NULL THEN '' ELSE '/' || p2.`id`END)" +
		" || (CASE WHEN p1.`id` IS NULL THEN '' ELSE '/' || p1.`id`END) AS 'ids_str'," +
		"    (CASE WHEN p10.`id` IS NULL THEN '' ELSE '/' || p10.`name` END)" +
		" || (CASE WHEN p9.`id` IS NULL THEN '' ELSE '/' || p9.`name` END)" +
		" || (CASE WHEN p8.`id` IS NULL THEN '' ELSE '/' || p8.`name` END)" +
		" || (CASE WHEN p7.`id` IS NULL THEN '' ELSE '/' || p7.`name` END)" +
		" || (CASE WHEN p6.`id` IS NULL THEN '' ELSE '/' || p6.`name` END)" +
		" || (CASE WHEN p5.`id` IS NULL THEN '' ELSE '/' || p5.`name` END)" +
		" || (CASE WHEN p4.`id` IS NULL THEN '' ELSE '/' || p4.`name` END)" +
		" || (CASE WHEN p3.`id` IS NULL THEN '' ELSE '/' || p3.`name` END)" +
		" || (CASE WHEN p2.`id` IS NULL THEN '' ELSE '/' || p2.`name` END)" +
		" || (CASE WHEN p1.`id` IS NULL THEN '' ELSE '/' || p1.`name` END) AS 'names_str'" +
		" FROM `note` p1" +
		" LEFT JOIN `note` p2 ON p2.`id` = p1.`pid`" +
		" LEFT JOIN `note` p3 ON p3.`id` = p2.`pid`" +
		" LEFT JOIN `note` p4 ON p4.`id` = p3.`pid`" +
		" LEFT JOIN `note` p5 ON p5.`id` = p4.`pid`" +
		" LEFT JOIN `note` p6 ON p6.`id` = p5.`pid`" +
		" LEFT JOIN `note` p7 ON p7.`id` = p6.`pid`" +
		" LEFT JOIN `note` p8 ON p8.`id` = p7.`pid`" +
		" LEFT JOIN `note` p9 ON p9.`id` = p8.`pid`" +
		" LEFT JOIN `note` p10 ON p10.`id` = p9.`pid`" +
		" WHERE p1.`id` = ?"
}
