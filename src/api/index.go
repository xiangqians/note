// index
// @author xiangqian
// @date 17:21 2023/02/04
package api

import (
	"github.com/gin-gonic/gin"
	"note/src/typ"
	"strings"
)

// IndexPage index页面
func IndexPage(pContext *gin.Context) {
	html := func(pf typ.File, fs []typ.File, err error) {
		Html(pContext, "index.html", gin.H{"pf": pf, "fs": fs}, err)
	}

	// id
	id, err := Query[int64](pContext, "id")
	//log.Printf("id = %d\n", id)

	// name
	name, err := Query[string](pContext, "name")
	name = strings.TrimSpace(name)
	//log.Printf("name = %s\n", name)

	// pf
	var pf typ.File
	if id < 0 {
		pf.Path = ""
	} else if id == 0 {
		pf.Path = "/"
	} else {
		sql := "SELECT f1.id, f1.pid, f1.`name`, f1.`type`, f1.`size`, f1.add_time, f1.upd_time, " +
			"((CASE WHEN f10.`name` IS NULL THEN '' ELSE '/' || f10.`name` END) " +
			"|| (CASE WHEN f9.`name` IS NULL THEN '' ELSE '/' || f9.`name` END) " +
			"|| (CASE WHEN f8.`name` IS NULL THEN '' ELSE '/' || f8.`name` END) " +
			"|| (CASE WHEN f7.`name` IS NULL THEN '' ELSE '/' || f7.`name` END) " +
			"|| (CASE WHEN f6.`name` IS NULL THEN '' ELSE '/' || f6.`name` END) " +
			"|| (CASE WHEN f5.`name` IS NULL THEN '' ELSE '/' || f5.`name` END) " +
			"|| (CASE WHEN f4.`name` IS NULL THEN '' ELSE '/' || f4.`name` END) " +
			"|| (CASE WHEN f3.`name` IS NULL THEN '' ELSE '/' || f3.`name` END) " +
			"|| (CASE WHEN f2.`name` IS NULL THEN '' ELSE '/' || f2.`name` END) " +
			"|| (CASE WHEN f1.`name` IS NULL THEN '' ELSE '/' || f1.`name` END))  AS 'path' " +
			"FROM `file` f1 " +
			"LEFT JOIN `file` f2 ON f2.del = 0 AND f2.`type` = 'd' AND f2.id = f1.pid " +
			"LEFT JOIN `file` f3 ON f3.del = 0 AND f3.`type` = 'd' AND f3.id = f2.pid " +
			"LEFT JOIN `file` f4 ON f4.del = 0 AND f4.`type` = 'd' AND f4.id = f3.pid " +
			"LEFT JOIN `file` f5 ON f5.del = 0 AND f5.`type` = 'd' AND f5.id = f4.pid " +
			"LEFT JOIN `file` f6 ON f6.del = 0 AND f6.`type` = 'd' AND f6.id = f5.pid " +
			"LEFT JOIN `file` f7 ON f7.del = 0 AND f7.`type` = 'd' AND f7.id = f6.pid " +
			"LEFT JOIN `file` f8 ON f8.del = 0 AND f8.`type` = 'd' AND f8.id = f7.pid " +
			"LEFT JOIN `file` f9 ON f9.del = 0 AND f9.`type` = 'd' AND f9.id = f8.pid " +
			"LEFT JOIN `file` f10 ON f10.del = 0 AND f10.`type` = 'd' AND f10.id = f9.pid " +
			"WHERE f1.del = 0 AND f1.`type` = 'd' AND f1.id = ? " +
			"GROUP BY f1.id"
		pf, _, err = DbQry[typ.File](pContext, sql, id)
		if err != nil {
			html(pf, nil, err)
			return
		}
	}

	// 查询
	args := make([]any, 0, 2)
	var fs []typ.File = nil
	var count int64
	sql := "SELECT f.`id`, f.`pid`, f.`name`, f.`type`, f.`size`, f.`add_time`, f.`upd_time` FROM `file` f WHERE f.`del` = 0 "
	if id >= 0 {
		sql += "AND f.`pid` = ? "
		args = append(args, id)
	}
	if name != "" {
		sql += "AND f.`name` LIKE '%' || ? || '%' "
		args = append(args, name)
	}
	sql += "ORDER BY f.`type` "
	if id < 0 {
		sql += "LIMIT 10000"
	}
	fs, count, err = DbQry[[]typ.File](pContext, sql, args...)
	if err != nil || count == 0 {
		fs = nil
	}

	html(pf, fs, err)
	return
}
