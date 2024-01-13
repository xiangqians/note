// @author xiangqian
// @date 23:21 2023/10/23
package common

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"note/src/db"
	"note/src/model"
	"note/src/session"
	util_string "note/src/util/string"
	"sort"
	"strconv"
	"strings"
)

func List(request *http.Request, writer http.ResponseWriter, session *session.Session, table string) (string, model.Response) {
	// 当前页
	current, _ := strconv.ParseInt(request.URL.Query().Get("current"), 10, 64)
	if current <= 0 {
		current = 1
	}

	// 页数量
	size64, _ := strconv.ParseInt(request.URL.Query().Get("size"), 10, 64)
	size := uint8(size64)
	if size <= 0 {
		if table == TableNote {
			size = 50
		} else {
			size = 10
		}
	}
	if size > 100 {
		size = 100
	}

	// 检索条件
	search := strings.TrimSpace(request.URL.Query().Get("search"))

	page := model.Page{
		Current: current,
		Size:    size,
		Search:  search,
		Total:   0,
	}

	columns := []column{
		{
			tables:     []string{TableImage, TableAudio, TableVideo, TableNote},
			name:       " id:",
			index:      -1,
			statement:  "t.`id` = ?",
			parseValue: parseInt64Value,
			value:      nil,
		},
		{
			tables:     nil,
			name:       " c:", // contain & child，是否包含子目录
			index:      -1,
			statement:  "",
			parseValue: parseBoolValue,
			value:      bool(false),
		},
		{
			tables:     []string{TableImage, TableAudio, TableVideo, TableNote},
			name:       " name:",
			index:      -1,
			statement:  "t.`name` LIKE '%' || ? || '%'", // sqlite在模糊查询时大小写不敏感
			parseValue: parseStrValue,
			value:      nil,
		},
		{
			tables:     []string{TableImage, TableAudio, TableVideo, TableNote},
			name:       " type:",
			index:      -1,
			statement:  "t.`type` = ?",
			parseValue: parseStrValue,
			value:      nil,
		},
		{
			tables:     []string{TableImage, TableAudio, TableVideo, TableNote},
			name:       " del:",
			index:      -1,
			statement:  "t.`del` = ?",
			parseValue: parseDelValue,
			value:      byte(0),
		},
	}

	length := len(columns)

	if search != "" {
		search = fmt.Sprintf(" %s", search)

		// 获取字段在检索文本search中的位置
		for i := 0; i < length; i++ {
			index := strings.Index(search, columns[i].name)
			(&columns[i]).index = index
		}

		// 字段切片，根据索引（index）升序排序
		sort.Slice(columns, func(i, j int) bool {
			return columns[i].index < columns[j].index
		})

		// 获取字段值
		for i := 0; i < length; i++ {
			index := columns[i].index
			if index == -1 {
				continue
			}

			var nextIndex int = -1
			for i := i + 1; i < length; i++ {
				nextIndex = columns[i].index
				if nextIndex != -1 {
					break
				}
			}

			if nextIndex == -1 {
				nextIndex = len(search)
			}

			strValue := strings.TrimSpace(search[index+len(columns[i].name) : nextIndex])
			value := columns[i].parseValue(strValue)
			if value != nil {
				(&columns[i]).value = value
			}
		}
	}

	// note
	if table == TableNote {
		return noteList(request, session, page, columns)
	}

	// html模板
	html := func(page model.Page, err any) (string, model.Response) {
		return fmt.Sprintf("%s/list", table),
			model.Response{Msg: session.GetMsg() + util_string.String(err), Data: map[string]any{
				"table": table,
				"page":  page,
			}}
	}

	db := db.Get()

	// len 0, cap ?
	statements := make([]string, 0, length)
	values := make([]any, 0, length)

	for i := 0; i < length; i++ {
		if columns[i].value != nil {
			tables := columns[i].tables
			for j := 0; j < len(tables); j++ {
				if tables[j] == table {
					statements = append(statements, columns[i].statement)
					values = append(values, columns[i].value)
					break
				}
			}
		}
	}

	sql := "SELECT t.`id`, t.`name`, t.`type`, t.`size`, t.`del`, t.`add_time`, t.`upd_time`"
	sql += fmt.Sprintf(" FROM `%s` t", table)
	sql += fmt.Sprintf(" WHERE %s ", strings.Join(statements, " AND "))
	sql += " ORDER BY (CASE WHEN t.`upd_time` > t.`add_time` THEN t.`upd_time` ELSE t.`add_time` END) DESC"
	result, err := db.Page(sql, current, size, values...)
	if err != nil {
		return html(page, err)
	}

	page.Total = result.Count()
	(&page).InitIndexes()

	var data any
	switch table {
	case TableImage, TableAudio, TableVideo:
		var abss []model.Abs
		err = result.Scan(&abss)
		data = abss
	}

	page.Data = data
	return html(page, err)
}

func noteList(request *http.Request, session *session.Session, page model.Page, columns []column) (string, model.Response) {
	// html模板
	html := func(page model.Page, pNote model.PNote, err any) (string, model.Response) {
		return fmt.Sprintf("%s/list", TableNote),
			model.Response{Msg: session.GetMsg() + util_string.String(err), Data: map[string]any{
				"table": TableNote,
				"page":  page,
				"pNote": pNote,
			}}
	}

	var pNote model.PNote

	var pid int64
	if strings.HasSuffix(request.URL.Path, "/list") {
		vars := mux.Vars(request)
		var err error
		pid, err = strconv.ParseInt(vars["pid"], 10, 64)
		if err != nil || pid < 0 {
			return html(page, pNote, err)
		}
	}

	db := db.Get()

	if pid > 0 {
		result, err := db.Get(getPNoteSql(), pid)
		if err != nil {
			return html(page, pNote, err)
		}

		err = result.Scan(&pNote)
		if err != nil {
			return html(page, pNote, err)
		}

		if pNote.IdsStr != "" {
			pNote.Ids = strings.Split(pNote.IdsStr, "/")[1:]
			pNote.Names = strings.Split(pNote.NamesStr, "/")[1:]
		}
	}

	pNote.Id = pid

	// contain & child，是否包含子目录
	c := false

	length := len(columns)

	// len 0, cap ?
	statements := make([]string, 0, length)
	values := make([]any, 0, length)

	for i := 0; i < length; i++ {
		if columns[i].value == nil {
			continue
		}

		if columns[i].name == " c:" {
			c = columns[i].value.(bool)
			continue
		}

		tables := columns[i].tables
		for j := 0; j < len(tables); j++ {
			if tables[j] == TableNote {
				statements = append(statements, columns[i].statement)
				values = append(values, columns[i].value)
				break
			}
		}
	}

	if !c {
		statements = append(statements, "t.`pid` = ?")
		values = append(values, pid)
	}

	pNote.C = c

	var sql string

	// 递归查询当前目录和子目录所有文件
	if c && pid > 0 {
		sql = fmt.Sprintf("WITH RECURSIVE `tmp`(`id`, `pid`, `name`, `type`, `size`, `del`, `add_time`, `upd_time`) AS ("+
			" SELECT t.`id`, t.`pid`, t.`name`, t.`type`, t.`size`, t.`del`, t.`add_time`, t.`upd_time`"+
			" FROM `note` t"+
			" WHERE t.`pid` = %d"+ // 起点条件
			" UNION ALL"+
			" SELECT t.`id`, t.`pid`, t.`name`, t.`type`, t.`size`, t.`del`, t.`add_time`, t.`upd_time`"+
			" FROM `note` t"+
			" INNER JOIN `tmp` ON t.pid = `tmp`.id)"+ // 关联递归查询结果
			" SELECT t.`id`, t.`pid`, t.`name`, t.`type`, t.`size`, t.`del`, t.`add_time`, t.`upd_time`", pid)
	} else
	// 查询当前目录文件 | 所有文件
	{
		sql = "SELECT t.`id`, t.`pid`, t.`name`, t.`type`, t.`size`, t.`del`, t.`add_time`, t.`upd_time`"
	}

	if c {
		sql += ", (CASE WHEN p10.`id` IS NULL THEN '' ELSE '/' || p10.`id` END)" +
			" || (CASE WHEN p9.`id` IS NULL THEN '' ELSE '/' || p9.`id`END)" +
			" || (CASE WHEN p8.`id` IS NULL THEN '' ELSE '/' || p8.`id`END)" +
			" || (CASE WHEN p7.`id` IS NULL THEN '' ELSE '/' || p7.`id`END)" +
			" || (CASE WHEN p6.`id` IS NULL THEN '' ELSE '/' || p6.`id`END)" +
			" || (CASE WHEN p5.`id` IS NULL THEN '' ELSE '/' || p5.`id`END)" +
			" || (CASE WHEN p4.`id` IS NULL THEN '' ELSE '/' || p4.`id`END)" +
			" || (CASE WHEN p3.`id` IS NULL THEN '' ELSE '/' || p3.`id`END)" +
			" || (CASE WHEN p2.`id` IS NULL THEN '' ELSE '/' || p2.`id`END)" +
			" || (CASE WHEN p1.`id` IS NULL THEN '' ELSE '/' || p1.`id`END) AS 'pids_str'," +
			"    (CASE WHEN p10.`id` IS NULL THEN '' ELSE '/' || p10.`name` END)" +
			" || (CASE WHEN p9.`id` IS NULL THEN '' ELSE '/' || p9.`name` END)" +
			" || (CASE WHEN p8.`id` IS NULL THEN '' ELSE '/' || p8.`name` END)" +
			" || (CASE WHEN p7.`id` IS NULL THEN '' ELSE '/' || p7.`name` END)" +
			" || (CASE WHEN p6.`id` IS NULL THEN '' ELSE '/' || p6.`name` END)" +
			" || (CASE WHEN p5.`id` IS NULL THEN '' ELSE '/' || p5.`name` END)" +
			" || (CASE WHEN p4.`id` IS NULL THEN '' ELSE '/' || p4.`name` END)" +
			" || (CASE WHEN p3.`id` IS NULL THEN '' ELSE '/' || p3.`name` END)" +
			" || (CASE WHEN p2.`id` IS NULL THEN '' ELSE '/' || p2.`name` END)" +
			" || (CASE WHEN p1.`id` IS NULL THEN '' ELSE '/' || p1.`name` END) AS 'pnames_str'"
	}

	// 递归查询当前目录和子目录所有文件
	if c && pid > 0 {
		sql += " FROM tmp t"
	} else
	// 查询当前目录文件 | 所有文件
	{
		sql += " FROM `note` t"
	}

	if c {
		sql += "" +
			" LEFT JOIN `note` p1 ON p1.`id` = t.`pid`" +
			" LEFT JOIN `note` p2 ON p2.`id` = p1.`pid`" +
			" LEFT JOIN `note` p3 ON p3.`id` = p2.`pid`" +
			" LEFT JOIN `note` p4 ON p4.`id` = p3.`pid`" +
			" LEFT JOIN `note` p5 ON p5.`id` = p4.`pid`" +
			" LEFT JOIN `note` p6 ON p6.`id` = p5.`pid`" +
			" LEFT JOIN `note` p7 ON p7.`id` = p6.`pid`" +
			" LEFT JOIN `note` p8 ON p8.`id` = p7.`pid`" +
			" LEFT JOIN `note` p9 ON p9.`id` = p8.`pid`" +
			" LEFT JOIN `note` p10 ON p10.`id` = p9.`pid`"
	}
	sql += fmt.Sprintf(" WHERE %s ", strings.Join(statements, " AND "))
	sql += " GROUP BY t.id"
	sql += " ORDER BY t.`type`, t.`name`, (CASE WHEN t.`upd_time` > t.`add_time` THEN t.`upd_time` ELSE t.`add_time` END) DESC"

	result, err := db.Page(sql, page.Current, page.Size, values...)
	if err != nil {
		return html(page, pNote, err)
	}

	page.Total = result.Count()
	(&page).InitIndexes()

	var notes []model.Note
	err = result.Scan(&notes)
	page.Data = notes
	if err == nil && c {
		length = len(notes)
		for i := 0; i < length; i++ {
			if notes[i].PidsStr != "" {
				notes[i].Pids = strings.Split(notes[i].PidsStr, "/")[1:]
				notes[i].Pnames = strings.Split(notes[i].PnamesStr, "/")[1:]
			}
		}
	}
	return html(page, pNote, err)
}

func parseStrValue(strValue string) any {
	return strValue
}

func parseInt64Value(strValue string) any {
	intValue, _ := strconv.ParseInt(strValue, 10, 64)
	return intValue
}

func parseBoolValue(strValue string) any {
	return strings.EqualFold(strValue, "true")
}

func parseDelValue(strValue string) any {
	del := parseBoolValue(strValue).(bool)
	if del {
		return byte(1)
	}
	return byte(0)
}

// 字段
type column struct {
	tables     []string         // 表
	name       string           // 字段名集
	index      int              // 字段在检索字符串（search）中索引
	statement  string           // 语句
	parseValue func(string) any // 解析值
	value      any              // 值（默认值）
}
