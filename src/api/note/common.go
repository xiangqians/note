// note
// @author xiangqian
// @date 17:50 2023/02/04
package note

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"note/src/api/common"
	api_common_context "note/src/api/common/context"
	"note/src/api/common/db"
	"note/src/api/common/session"
	"note/src/typ"
	"note/src/util/json"
	util_os "note/src/util/os"
	"note/src/util/str"
	"note/src/util/time"
	"os"
	"sort"
	"strings"
)

const NoteSessionKey = "note"

// DeserializeHist 反序列化历史记录
func DeserializeHist(hist string) ([]typ.Note, error) {
	if hist == "" {
		return nil, nil
	}

	// hists
	hists := make([]typ.Note, 0, 1) // len 0, cap ?
	err := json.Deserialize(hist, &hists)
	if err != nil {
		return nil, err
	}

	// sort
	Sort(&hists)

	return hists, nil
}

// SerializeHist 序列化历史记录
func SerializeHist(hists []typ.Note) (string, error) {
	return json.Serialize(hists)
}

// Sort 对notes进行排序
func Sort(notes *[]typ.Note) {
	sort.Slice(*notes, func(i, j int) bool {
		return (*notes)[i].UpdTime > (*notes)[j].UpdTime
	})
}

func RedirectToList(context *gin.Context, pid int64, err any) {
	resp := typ.Resp[any]{
		Msg: str.ConvTypeToStr(err),
	}

	// 记录查询参数
	note, err := session.Get[typ.Note](context, "note", false)
	if err == nil {
		api_common_context.Redirect(context, fmt.Sprintf("/note/list?pid=%d&deleted=%d", pid, note.Del), resp)
		return
	}

	api_common_context.Redirect(context, fmt.Sprintf("/note/list?pid=%d", pid), resp)
}

// DbList 查询列表
func DbList(context *gin.Context, note typ.Note) ([]typ.Note, int64, error) {
	// sql
	sql, args := DbQrySql(note,
		"ORDER BY n.`type`, n.`name`, (CASE WHEN n.`upd_time` > n.`add_time` THEN n.`upd_time` ELSE n.`add_time` END) DESC ", "LIMIT 10000")

	// qry
	notes, count, err := db.Qry[[]typ.Note](context, sql, args...)
	if err != nil || count == 0 {
		notes = nil
	}

	// init path
	if note.QryPath > 0 && notes != nil {
		for i, l := 0, len(notes); i < l; i++ {
			InitPath(&notes[i])
		}
	}

	return notes, count, err
}

// DbQry 查询
// id: 主键id
// qryPath: 查询路径，0-不查询，1-查询，2-查询并包含自身的
// del: 删除标识
func DbQry(context *gin.Context, id int64, qryPath int8, del byte) (typ.Note, int64, error) {
	// sql
	sql, args := DbQrySql(typ.Note{
		Abs: typ.Abs{
			Id:  id,
			Del: del,
		},
		Pid:     -1,
		QryPath: qryPath,
	}, "LIMIT 1")

	// qry
	note, count, err := db.Qry[typ.Note](context, sql, args...)
	if qryPath > 0 && err == nil && count > 0 {
		InitPath(&note)
	}

	return note, count, err
}

// InitPath 初始化Note path & pathLink
func InitPath(note *typ.Note) {
	path := (*note).Path
	if path == "" {
		return
	}

	pathArr := strings.Split(path, "/")
	l := len(pathArr)
	pathLinkArr := make([]string, 0, l) // len 0, cap ?
	for i := 0; i < l; i++ {
		e := pathArr[i]
		if e == "" {
			pathLinkArr = append(pathLinkArr, "")
			continue
		}

		eArr := strings.Split(e, ":")
		pathArr[i] = eArr[1]
		pathLink := fmt.Sprintf("<a href=\"/note/list?pid=%s&t=%d\">%s</a>\n", eArr[0], time.NowUnix(), eArr[1])
		pathLinkArr = append(pathLinkArr, pathLink)
	}
	(*note).Path = strings.Join(pathArr, "/")
	(*note).PathLink = strings.Join(pathLinkArr, "/")
}

// DbQrySql 查询sql
// note: 查询实体类
func DbQrySql(note typ.Note, last ...string) (string, []any) {
	args := make([]any, 0, 1)
	sql := "SELECT n.`id`, n.`pid`, n.`name`, n.`type`, n.`size`, n.`hist`, n.`hist_size`, n.`del`, n.`add_time`, n.`upd_time` "

	// path sql
	qryPathSql := "(CASE WHEN pn10.`id` IS NULL THEN '' ELSE '/' || pn10.`id` || ':' ||pn10.`name` END) " +
		"|| (CASE WHEN pn9.`id` IS NULL THEN '' ELSE '/' || pn9.`id` || ':' || pn9.`name` END) " +
		"|| (CASE WHEN pn8.`id` IS NULL THEN '' ELSE '/' || pn8.`id` || ':' || pn8.`name` END) " +
		"|| (CASE WHEN pn7.`id` IS NULL THEN '' ELSE '/' || pn7.`id` || ':' || pn7.`name` END) " +
		"|| (CASE WHEN pn6.`id` IS NULL THEN '' ELSE '/' || pn6.`id` || ':' || pn6.`name` END) " +
		"|| (CASE WHEN pn5.`id` IS NULL THEN '' ELSE '/' || pn5.`id` || ':' || pn5.`name` END) " +
		"|| (CASE WHEN pn4.`id` IS NULL THEN '' ELSE '/' || pn4.`id` || ':' || pn4.`name` END) " +
		"|| (CASE WHEN pn3.`id` IS NULL THEN '' ELSE '/' || pn3.`id` || ':' || pn3.`name` END) " +
		"|| (CASE WHEN pn2.`id` IS NULL THEN '' ELSE '/' || pn2.`id` || ':' || pn2.`name` END) " +
		"|| (CASE WHEN pn1.`id` IS NULL THEN '' ELSE '/' || pn1.`id` || ':' || pn1.`name` END) "
	switch note.QryPath {
	// 查询路径
	case 1:
		sql += fmt.Sprintf(", CASE WHEN n.`pid` = 0 THEN  '/' ELSE (%s) END AS 'path' ", qryPathSql)

	// 查询包含自身的
	case 2:
		sql += fmt.Sprintf(", (%s || (CASE WHEN n.`id` IS NULL THEN '' ELSE '/' || n.`id` || ':' || n.`name` END)) AS 'path' ", qryPathSql)

	default:
		//
	}

	sql += "FROM `note` n "
	if note.QryPath == 1 || note.QryPath == 2 {
		sql += "LEFT JOIN `note` pn1 ON pn1.`type` = 'd' AND pn1.id = n.pid " +
			"LEFT JOIN `note` pn2 ON pn2.`type` = 'd' AND pn2.id = pn1.pid " +
			"LEFT JOIN `note` pn3 ON pn3.`type` = 'd' AND pn3.id = pn2.pid " +
			"LEFT JOIN `note` pn4 ON pn4.`type` = 'd' AND pn4.id = pn3.pid " +
			"LEFT JOIN `note` pn5 ON pn5.`type` = 'd' AND pn5.id = pn4.pid " +
			"LEFT JOIN `note` pn6 ON pn6.`type` = 'd' AND pn6.id = pn5.pid " +
			"LEFT JOIN `note` pn7 ON pn7.`type` = 'd' AND pn7.id = pn6.pid " +
			"LEFT JOIN `note` pn8 ON pn8.`type` = 'd' AND pn8.id = pn7.pid " +
			"LEFT JOIN `note` pn9 ON pn9.`type` = 'd' AND pn9.id = pn8.pid " +
			"LEFT JOIN `note` pn10 ON pn10.`type` = 'd' AND pn10.id = pn9.pid "
	}

	// contains sub
	if note.ContainsSub != 0 && note.Pid > 0 {

		// 递归查询所有子节点
		//WITH RECURSIVE tmp AS (
		//	SELECT n.* FROM note n WHERE n.`del` = 0 AND n.`id` = 1
		//	UNION ALL
		//	SELECT n.* FROM tmp t JOIN note n ON n.`pid` = t.`id`
		//)
		//SELECT * FROM tmp

		sql += "JOIN ( " +
			"WITH RECURSIVE tmp(id, pid) AS ( " +
			"SELECT n.id, n.pid FROM note n WHERE n.`del` = 0 AND n.`id` = ? " +
			"UNION ALL " +
			"SELECT n.id, n.pid FROM tmp t JOIN note n ON n.`pid` = t.`id`" +
			") " +
			"SELECT id FROM tmp " +
			") r ON r.id = n.pid "
		args = append(args, note.Pid)
	}

	// del, 只查询 0/1 状态
	if note.Del != 0 {
		sql += "WHERE n.`del` = 1 "
	} else {
		sql += "WHERE n.`del` = 0 "
	}

	// id
	if note.Id > 0 {
		sql += "AND n.`id` = ? "
		args = append(args, note.Id)
	}

	// pid
	if note.ContainsSub == 0 && note.Pid >= 0 {
		sql += "AND n.`pid` = ? "
		args = append(args, note.Pid)
	}

	// name
	if note.Name != "" {
		sql += "AND n.`name` LIKE '%' || ? || '%' "
		args = append(args, note.Name)
	}

	// type
	if note.Type != "" {
		sql += "AND n.`type` = ? "
		args = append(args, note.Type)
	}

	sql += "GROUP BY n.id "

	// last
	if last != nil && len(last) > 0 {
		for _, e := range last {
			sql += e
		}
	}

	return sql, args
}

// ReadHist 读取笔记历史信息
func ReadHist(context *gin.Context, note typ.Note) ([]byte, error) {
	return read(context, note, true)
}

// Read 读取笔记信息
func Read(context *gin.Context, note typ.Note) ([]byte, error) {
	return read(context, note, false)
}

func read(context *gin.Context, note typ.Note, hist bool) ([]byte, error) {
	// file path
	var path string
	var err error
	if hist {
		path, err = HistPath(context, note)
	} else {
		path, err = Path(context, note)
	}
	if err != nil {
		return nil, err
	}

	// open file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// read file
	buf, err := io.ReadAll(file)
	return buf, err
}

// DelHistNote 删除历史笔记
func DelHistNote(context *gin.Context, note typ.Note) (string, error) {
	// path
	path, err := HistPath(context, note)
	if err != nil {
		return path, err
	}

	// del
	return path, util_os.DelFile(path)
}

// HistPath 获取笔记历史记录物理路径
func HistPath(context *gin.Context, note typ.Note) (string, error) {
	// dir
	dataDir := common.DataDir(context)
	noteDir := fmt.Sprintf("%s%s%s%s%s%s%s", dataDir,
		util_os.FileSeparator(), "hist",
		util_os.FileSeparator(), "note",
		util_os.FileSeparator(), note.Type)
	if !util_os.IsExist(noteDir) {
		err := util_os.MkDir(noteDir)
		if err != nil {
			return "", err
		}
	}

	// file name
	time := note.UpdTime
	name := fmt.Sprintf("%d_%d", note.Id, time)

	// path
	return fmt.Sprintf("%s%s%s", noteDir, util_os.FileSeparator(), name), nil
}

// DelNote 删除笔记
func DelNote(context *gin.Context, note typ.Note) (string, error) {
	// path
	path, err := Path(context, note)
	if err != nil {
		return path, err
	}

	// del
	return path, util_os.DelFile(path)
}

// ClearNote 清空笔记
func ClearNote(context *gin.Context, note typ.Note) (string, error) {
	// path
	path, err := Path(context, note)
	if err != nil {
		return path, err
	}

	// exist ?
	if !util_os.IsExist(path) {
		return path, nil
	}

	// open
	file, err := os.OpenFile(path,
		os.O_WRONLY|os.O_TRUNC, // 只写（O_WRONLY） & 清空文件（O_TRUNC）
		0666)

	// close
	file.Close()

	return path, err
}

// Path 获取文件物理路径
func Path(context *gin.Context, note typ.Note) (string, error) {
	// dir
	dataDir := common.DataDir(context)
	noteDir := fmt.Sprintf("%s%s%s%s%s", dataDir,
		util_os.FileSeparator(), "note",
		util_os.FileSeparator(), note.Type)
	if !util_os.IsExist(noteDir) {
		err := util_os.MkDir(noteDir)
		if err != nil {
			return "", err
		}
	}

	// file name
	name := fmt.Sprintf("%d", note.Id)

	// path
	return fmt.Sprintf("%s%s%s", noteDir, util_os.FileSeparator(), name), nil
}
