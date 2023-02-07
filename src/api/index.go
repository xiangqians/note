// index
// @author xiangqian
// @date 17:21 2023/02/04
package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"note/src/typ"
	"time"
)

func IndexPage(pContext *gin.Context) {
	html := func(pf typ.File, fs []typ.File, err error) {
		Html(pContext, "index.html", gin.H{"pf": pf, "fs": fs}, err)
	}

	id, err := Query[int64](pContext, "id")
	log.Printf("id = %d\n", id)

	// pf
	var pf typ.File
	if id == 0 {
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

	// 查询目录下的所有目录和文件
	fs, count, err := DbQry[[]typ.File](pContext, "SELECT f.id, f.pid, f.`name`, f.`type`, f.`size`, f.add_time, f.upd_time FROM `file` f WHERE f.del = 0 AND f.pid = ?", id)
	if err != nil {
		html(pf, nil, err)
		return
	}

	if count == 0 {
		fs = nil
	}

	html(pf, fs, nil)
	return
}

// FileAdd 新增文件
func FileAdd(pContext *gin.Context) {
	redirect := func(id int64, msg any) {
		Redirect(pContext, fmt.Sprintf("/?id=%d", id), nil, msg)
	}

	// file
	f := typ.File{}
	err := ShouldBind(pContext, &f)
	pid := f.Pid
	if err != nil {
		redirect(pid, err)
		return
	}

	// add
	_, err = DbAdd(pContext, "INSERT INTO `file` (`pid`, `name`, `type`, `add_time`) VALUES (?, ?, ?, ?)", f.Pid, f.Name, f.Type, time.Now().Unix())
	if err != nil {
		redirect(pid, err)
		return
	}

	redirect(pid, nil)
	return
}

// FileRename 文件重命名
func FileRename(pContext *gin.Context) {
	redirect := func(id int64, msg any) {
		Redirect(pContext, fmt.Sprintf("/?id=%d", id), nil, msg)
	}

	// file
	f := typ.File{}
	err := ShouldBind(pContext, &f)
	pid := f.Pid
	if err != nil {
		redirect(pid, err)
		return
	}

	// update
	_, err = DbUpd(pContext, "UPDATE `file` SET `name` = ?, `upd_time` = ? WHERE id = ?", f.Name, time.Now().Unix(), f.Id)
	if err != nil {
		redirect(pid, err)
		return
	}

	redirect(pid, nil)
	return
}

// FileDel 删除文件
func FileDel(pContext *gin.Context) {
	redirect := func(id int64, msg any) {
		Redirect(pContext, fmt.Sprintf("/?id=%d", id), nil, msg)
	}

	id, _ := Param[int64](pContext, "id")
	pid, _, _ := DbQry[int64](pContext, "SELECT f.pid FROM `file` f WHERE f.del = 0 AND f.id = ?", id)
	_, err := DbDel(pContext, "UPDATE `file` SET del = 1, `upd_time` = ? WHERE id = ?", time.Now().Unix(), id)
	if err != nil {
		redirect(pid, err)
		return
	}

	redirect(pid, nil)
	return
}
