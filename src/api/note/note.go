// note
// @author xiangqian
// @date 17:50 2023/02/04
package note

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday/v2"
	"io"
	"log"
	"net/http"
	"note/src/api/common"
	typ_api "note/src/typ/api"
	typ_ft "note/src/typ/ft"
	typ_resp "note/src/typ/resp"
	util_os "note/src/util/os"
	util_str "note/src/util/str"
	"os"
	"strings"
	"time"
)

// List 文件列表页面
func List(context *gin.Context) {
	html := func(pnote typ_api.Note, notes []typ_api.Note, err error) {
		resp := typ_resp.Resp[map[string]any]{
			Msg: util_str.TypeToStr(err),
			Data: map[string]any{
				"pnote": pnote,
				"notes": notes,
			},
		}
		common.HtmlOk(context, "note/list.html", resp)
	}

	// id
	pid, err := common.Query[int64](context, "pid")
	//log.Printf("id = %d\n", id)

	// name
	name, err := common.Query[string](context, "name")
	name = strings.TrimSpace(name)
	//log.Printf("name = %s\n", name)

	// pnote
	var pnote typ_api.Note
	if pid < 0 {
		pnote.Path = ""

	} else if pid == 0 {
		pnote.Path = "/"

	} else {
		sql := "SELECT n1.id, n1.pid, n1.`name`, n1.`type`, n1.`size`, n1.add_time, n1.upd_time, " +
			"( (CASE WHEN n10.`id` IS NULL THEN '' ELSE '/' ||n10.`id` || ':' ||n10.`name` END) " +
			"|| (CASE WHEN n9.`id` IS NULL THEN '' ELSE '/' || n9.`id` || ':' || n9.`name` END) " +
			"|| (CASE WHEN n8.`id` IS NULL THEN '' ELSE '/' || n8.`id` || ':' || n8.`name` END) " +
			"|| (CASE WHEN n7.`id` IS NULL THEN '' ELSE '/' || n7.`id` || ':' || n7.`name` END) " +
			"|| (CASE WHEN n6.`id` IS NULL THEN '' ELSE '/' || n6.`id` || ':' || n6.`name` END) " +
			"|| (CASE WHEN n5.`id` IS NULL THEN '' ELSE '/' || n5.`id` || ':' || n5.`name` END) " +
			"|| (CASE WHEN n4.`id` IS NULL THEN '' ELSE '/' || n4.`id` || ':' || n4.`name` END) " +
			"|| (CASE WHEN n3.`id` IS NULL THEN '' ELSE '/' || n3.`id` || ':' || n3.`name` END) " +
			"|| (CASE WHEN n2.`id` IS NULL THEN '' ELSE '/' || n2.`id` || ':' || n2.`name` END) " +
			"|| (CASE WHEN n1.`id` IS NULL THEN '' ELSE '/' || n1.`id` || ':' || n1.`name` END))  AS 'path' " +
			"FROM `note` n1 " +
			"LEFT JOIN `note` n2 ON n2.del = 0 AND n2.`type` = 'd' AND n2.id = n1.pid " +
			"LEFT JOIN `note` n3 ON n3.del = 0 AND n3.`type` = 'd' AND n3.id = n2.pid " +
			"LEFT JOIN `note` n4 ON n4.del = 0 AND n4.`type` = 'd' AND n4.id = n3.pid " +
			"LEFT JOIN `note` n5 ON n5.del = 0 AND n5.`type` = 'd' AND n5.id = n4.pid " +
			"LEFT JOIN `note` n6 ON n6.del = 0 AND n6.`type` = 'd' AND n6.id = n5.pid " +
			"LEFT JOIN `note` n7 ON n7.del = 0 AND n7.`type` = 'd' AND n7.id = n6.pid " +
			"LEFT JOIN `note` n8 ON n8.del = 0 AND n8.`type` = 'd' AND n8.id = n7.pid " +
			"LEFT JOIN `note` n9 ON n9.del = 0 AND n9.`type` = 'd' AND n9.id = n8.pid " +
			"LEFT JOIN `note` n10 ON n10.del = 0 AND n10.`type` = 'd' AND n10.id = n9.pid " +
			"WHERE n1.del = 0 AND n1.`type` = 'd' AND n1.id = ? " +
			"GROUP BY n1.id"
		pnote, _, err = common.DbQry[typ_api.Note](context, sql, pid)
		if err != nil {
			html(pnote, nil, err)
			return
		}

		pathArr := strings.Split(pnote.Path, "/")
		l := len(pathArr)
		pathLinkArr := make([]string, 0, l) // len 0, cap ?
		for i := 0; i < l; i++ {
			v := pathArr[i]
			if v == "" {
				pathLinkArr = append(pathLinkArr, "")
				continue
			}

			vArr := strings.Split(v, ":")
			pathArr[i] = vArr[1]

			pathLink := fmt.Sprintf("<a href=\"/note/list?pid=%s\">%s</a>\n", vArr[0], vArr[1])
			pathLinkArr = append(pathLinkArr, pathLink)
		}
		//pnote.Path = strings.Join(pathArr, "/")
		pnote.Path = strings.Join(pathLinkArr, "/")
	}

	// 查询
	args := make([]any, 0, 2)
	var notes []typ_api.Note = nil
	var count int64
	sql := "SELECT `id`, `pid`, `name`, `type`, `size`, `add_time`, `upd_time` FROM `note` WHERE `del` = 0 "
	// id
	if pid >= 0 {
		sql += "AND `pid` = ? "
		args = append(args, pid)
	}
	// name
	if name != "" {
		sql += "AND `name` LIKE '%' || ? || '%' "
		args = append(args, name)
	}
	sql += "ORDER BY `type`, `name`, (CASE WHEN `upd_time` > `add_time` THEN `upd_time` ELSE `add_time` END) DESC "
	if pid < 0 {
		sql += "LIMIT 10000"
	}
	notes, count, err = common.DbQry[[]typ_api.Note](context, sql, args...)
	if err != nil || count == 0 {
		notes = nil
	}

	html(pnote, notes, err)
	return
}

// Add 新增文件
func Add(context *gin.Context) {
	redirect := func(pid int64, err any) {
		resp := typ_resp.Resp[any]{Msg: util_str.TypeToStr(err)}
		common.Redirect(context, fmt.Sprintf("/note/list?pid=%d", pid), resp)
	}

	// note
	note := typ_api.Note{}
	err := common.ShouldBind(context, &note)
	pid := note.Pid
	if err != nil {
		redirect(pid, err)
		return
	}

	// name
	note.Name = strings.TrimSpace(note.Name)
	err = util_os.VerifyFileName(note.Name)
	if err != nil {
		redirect(pid, err)
		return
	}

	// 校验文件类型
	// 只支持添加 目录 和 md文件
	fType := typ_ft.ExtNameOf(strings.TrimSpace(note.Type))
	if !(fType == typ_ft.FtD || fType == typ_ft.FtMd) {
		redirect(pid, fmt.Sprintf("%s, %s", errors.New(i18n.MustGetMessage("i18n.fileTypeUnsupported")), note.Type))
		return
	}
	note.Type = string(fType)

	// add
	id, err := common.DbAdd(context, "INSERT INTO `note` (`pid`, `name`, `type`, `add_time`) VALUES (?, ?, ?, ?)", note.Pid, note.Name, note.Type, time.Now().Unix())
	if err != nil {
		redirect(pid, err)
		return
	}
	note.Id = id

	// 如果不是目录，则创建物理文件
	if fType != typ_ft.FtD {
		// path
		path, err := Path(context, note)
		if err != nil {
			redirect(pid, err)
			return
		}

		// create
		file, err := os.Create(path)
		if err != nil {
			redirect(pid, err)
			return
		}
		defer file.Close()
	}

	redirect(pid, err)
	return
}

// Upload 上传文件
func Upload(context *gin.Context) {
	var id int64
	var pid int64
	var err error

	// method
	method := context.Request.Method

	// 有id则是put
	id, err = common.PostForm[int64](context, "id")
	if err == nil && id > 0 {
		method = http.MethodPut
	}

	// redirect
	redirect := func(id int64, pid int64, err any) {
		resp := typ_resp.Resp[any]{
			Msg: util_str.TypeToStr(err),
		}
		switch method {
		case http.MethodPost:
			common.Redirect(context, fmt.Sprintf("/note/list?pid=%d", pid), resp)

		case http.MethodPut:
			common.Redirect(context, fmt.Sprintf("/note/%d/edit", id), resp)
		}
	}

	// 获取 id 或 pid
	switch method {
	// 上传文件，需要pid
	case http.MethodPost:
		pid, err = common.PostForm[int64](context, "pid")

	// 重传文件，需要id
	case http.MethodPut:
		id, err = common.PostForm[int64](context, "id")
	}
	if err != nil {
		redirect(id, pid, err)
		return
	}

	// FileHeader
	fh, err := context.FormFile("file")
	if err != nil || fh == nil {
		redirect(id, pid, err)
		return
	}

	//log.Printf("Filename: %v\n", fh.Filename)
	//log.Printf("Size: %v\n", fh.Size)
	//log.Printf("Header: %v\n", fh.Header)

	// name
	fn := fh.Filename

	// type
	// 校验文件类型，只支持上传 html/pdf/zip
	ftStr := ""
	index := strings.LastIndex(fn, ".")
	if index > 0 {
		ftStr = fn[index+1:]
	}
	ft := typ_ft.ExtNameOf(strings.TrimSpace(ftStr))
	if ft == typ_ft.FtUnk || !(ft == typ_ft.FtHtml || ft == typ_ft.FtPdf || ft == typ_ft.FtZip) {
		redirect(id, pid, fmt.Sprintf("%s, %s", errors.New(i18n.MustGetMessage("i18n.fileTypeUnsupported")), ftStr))
		return
	}

	fn = fn[:index]

	// size
	fs := fh.Size

	// 校验 id 或 pid
	switch method {
	case http.MethodPost:
		// 校验 pid 是否存在
		if pid != 0 {
			note, count, err := DbQry(context, pid)
			if err != nil || count == 0 || typ_ft.ExtNameOf(note.Type) != typ_ft.FtD { // 父节点必须是目录
				redirect(id, pid, err)
				return
			}
		}

	case http.MethodPut:
		// 校验 id 是否存在
		note, count, err := DbQry(context, id)
		if err != nil || count == 0 {
			redirect(id, pid, err)
			return
		}

		if ft != typ_ft.ExtNameOf(note.Type) {
			redirect(id, pid, "重传必须是相同文件类型")
			return
		}
	}

	// 操作数据库
	switch method {
	case http.MethodPost:
		id, err = common.DbAdd(context, "INSERT INTO `note` (`pid`, `name`, `type`, `size`, `add_time`) VALUES (?, ?, ?, ?, ?)",
			pid, fn, ft, fs, time.Now().Unix())

	case http.MethodPut:
		_, err = common.DbUpd(context, "UPDATE `note` SET `name` = ?, `type` = ?, `size` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ?",
			fn, ft, fs, time.Now().Unix(), id)
	}
	if err != nil {
		redirect(id, pid, err)
		return
	}

	// path
	f := typ_api.Note{}
	f.Id = id
	f.Type = string(ft)
	fp, err := Path(context, f)
	if err != nil {
		redirect(id, pid, err)
		return
	}

	// 清空文件
	if method == http.MethodPut && util_os.IsExist(fp) {
		pFile, err := os.OpenFile(fp,
			os.O_WRONLY|os.O_TRUNC, // 只写（O_WRONLY） & 清空文件（O_TRUNC）
			0666)
		if err != nil {
			redirect(id, pid, err)
			return
		}
		pFile.Close()
	}

	// 保存文件
	err = context.SaveUploadedFile(fh, fp)

	redirect(id, pid, err)
	return
}

// UpdName 文件重命名
func UpdName(context *gin.Context) {
	redirect := func(pid int64, err any) {
		resp := typ_resp.Resp[any]{
			Msg: util_str.TypeToStr(err),
		}
		common.Redirect(context, fmt.Sprintf("/note/list?pid=%d", pid), resp)
	}

	// file
	f := typ_api.Note{}
	err := common.ShouldBind(context, &f)
	pid := f.Pid
	if err != nil {
		redirect(pid, err)
		return
	}

	// name
	f.Name = strings.TrimSpace(f.Name)
	err = util_os.VerifyFileName(f.Name)
	if err != nil {
		redirect(pid, err)
		return
	}

	//fType, count, err := common.DbQry[string](context, "SELECT `type` FROM `note` WHERE `del` = 0 AND `id` = ?", f.Id)
	//if count > 0 {
	//	name := f.Name
	//	ft := typ.FtOf(fType)
	//	if ft != typ.FtD && ft != typ.FtUnk && !strings.HasSuffix(name, string(ft)) {
	//		name = fmt.Sprintf("%s.%s", name, string(ft))
	//	}
	//
	//	// update
	//	_, err = common.DbUpd(context, "UPDATE `note` SET `name` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `name` <> ?", name, time.Now().Unix(), f.Id, name)
	//}

	// update
	_, err = common.DbUpd(context, "UPDATE `note` SET `name` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `name` <> ?", f.Name, time.Now().Unix(), f.Id, f.Name)

	redirect(pid, err)
	return
}

// Cut 剪切文件
func Cut(context *gin.Context) {
	redirect := func(id int64, err any) {
		resp := typ_resp.Resp[any]{
			Msg: util_str.TypeToStr(err),
		}
		common.Redirect(context, fmt.Sprintf("/note/list?pid=%d", id), resp)
	}

	// dst id
	dstId, err := common.Param[int64](context, "dstId")
	if err != nil {
		redirect(dstId, err)
		return
	}

	// src id
	srcId, err := common.Param[int64](context, "srcId")
	if err != nil {
		redirect(dstId, err)
		return
	}

	// dst
	if dstId != 0 {
		f, _, err := common.DbQry[typ_api.Note](context, "SELECT id, pid, `name`, `type`, `size`, `add_time`, `upd_time` FROM `note` WHERE `del` = 0 AND `id` = ?", dstId)
		if err != nil || typ_ft.FtD != typ_ft.ExtNameOf(f.Type) {
			redirect(dstId, err)
			return
		}
	}

	// update
	_, err = common.DbUpd(context, "UPDATE `note` SET `pid` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `pid` <> ?",
		dstId,
		time.Now().Unix(),
		srcId,
		dstId)

	redirect(dstId, err)
	return
}

// Del 删除文件
func Del(context *gin.Context) {
	redirect := func(id int64, err any) {
		resp := typ_resp.Resp[any]{Msg: util_str.TypeToStr(err)}
		common.Redirect(context, fmt.Sprintf("/note/list?pid=%d", id), resp)
	}

	// id
	id, err := common.Param[int64](context, "id")
	if err != nil {
		redirect(0, err)
		return
	}

	// pid
	pid, _, err := common.DbQry[int64](context, "SELECT `pid` FROM `note` WHERE `del` = 0 AND `id` = ?", id)
	if err != nil {
		redirect(pid, err)
		return
	}

	// Delete not supported
	redirect(pid, "Delete not supported")
	return

	// update
	_, err = common.DbDel(context, "UPDATE `note` SET `del` = 1, `upd_time` = ? WHERE `id` = ?", time.Now().Unix(), id)

	redirect(pid, err)
	return
}

func Get(context *gin.Context) {
	// id
	id, err := common.Param[int64](context, "id")
	if err != nil {
		log.Println(err)
		return
	}

	// f
	f, count, err := common.DbQry[typ_api.Note](context, "SELECT `id`, `pid`, `name`, `type`, `size`, `add_time`, `upd_time` FROM `note` WHERE `del` = 0 AND `id` = ?", id)
	if err != nil || count == 0 {
		log.Println(err)
		return
	}

	// 排除目录
	if typ_ft.FtD == typ_ft.ExtNameOf(f.Type) {
		return
	}

	// path
	fPath, err := Path(context, f)
	if err != nil {
		log.Println(err)
		return
	}

	/**
	// read all
	buf, err := os.ReadFile(fPath)
	if err != nil {
		log.Println(err)
		return
	}
	writer := context.Writer
	writer.Write(buf)
	writer.Flush()
	*/

	/**
	// open
	pFile, err := os.Open(fPath)
	if err != nil {
		log.Println(err)
		return
	}

	// write
	err = util.IOCopy(pFile, context.Writer, 0)
	if err != nil {
		log.Println(err)
		return
	}
	*/

	context.File(fPath)

	return
}

// View 查看文件页面
func View(context *gin.Context) {
	// id
	id, err := common.Param[int64](context, "id")
	if err != nil {
		DefaultView(context, typ_api.Note{}, err)
		return
	}

	// query
	note, count, err := common.DbQry[typ_api.Note](context, "SELECT `id`, `pid`, `name`, `type`, `size`, `add_time`, `upd_time` FROM `note` WHERE `del` = 0 AND `id` = ?", id)
	if err != nil || count == 0 {
		DefaultView(context, note, err)
		return
	}

	// type
	switch typ_ft.ExtNameOf(note.Type) {
	// markdown
	case typ_ft.FtMd:
		MdView(context, note)

	// html
	case typ_ft.FtHtml:
		HtmlView(context, note)

	// pdf
	case typ_ft.FtPdf:
		PdfView(context, note)

	// default
	default:
		DefaultView(context, note, err)
	}
}

// DefaultView 默认查看文件
func DefaultView(context *gin.Context, note typ_api.Note, err error) {
	resp := typ_resp.Resp[typ_api.Note]{
		Msg:  util_str.TypeToStr(err),
		Data: note,
	}
	common.HtmlOk(context, "note/default/view.html", resp)
}

// MdView 查看md文件
// https://github.com/russross/blackfriday
// https://pkg.go.dev/github.com/russross/blackfriday/v2
func MdView(context *gin.Context, note typ_api.Note) {
	html := func(html string, err any) {
		resp := typ_resp.Resp[map[string]any]{
			Msg: util_str.TypeToStr(err),
			Data: map[string]any{
				"note": note,
				"html": html,
			},
		}
		common.HtmlOk(context, "note/md/view.html", resp)
	}

	// read
	buf, err := FileRead(context, note)
	if err != nil {
		html("", err)
		return
	}

	//output := blackfriday.Run(input)
	//output := blackfriday.Run(input, blackfriday.WithNoExtensions())
	//output := blackfriday.Run(input, blackfriday.WithExtensions(blackfriday.CommonExtensions))

	// https://github.com/russross/blackfriday/issues/394
	buf = bytes.Replace(buf, []byte("\r"), nil, -1)
	//output := blackfriday.Run(input, blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.HardLineBreak))
	buf = blackfriday.Run(buf, blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.HardLineBreak|blackfriday.AutoHeadingIDs|blackfriday.Autolink))

	// 安全过滤
	//buf = bluemonday.UGCPolicy().SanitizeBytes(buf)

	html(string(buf), nil)
}

// HtmlView 查看html文件
func HtmlView(context *gin.Context, note typ_api.Note) {
	html := func(html string, err any) {
		resp := typ_resp.Resp[map[string]any]{
			Msg: util_str.TypeToStr(err),
			Data: map[string]any{
				"note": note,
				"html": html,
			},
		}
		common.HtmlOk(context, "note/html/view.html", resp)
	}

	// read
	buf, err := FileRead(context, note)
	if err != nil {
		html("", err)
		return
	}

	html(string(buf), nil)
}

// PdfView 查看pdf文件
func PdfView(context *gin.Context, note typ_api.Note) {
	v, _ := common.Query[string](context, "v")
	v = strings.TrimSpace(v)
	switch v {
	case "1.0":
		// v1.0
	case "2.0":
		// v2.0
	default:
		v = "2.0"
	}

	note.Url = fmt.Sprintf("/note/%v", note.Id)

	resp := typ_resp.Resp[typ_api.Note]{
		Data: note,
	}
	common.HtmlOk(context, fmt.Sprintf("note/pdf/view_v%s.html", v), resp)
}

// Edit 文件修改页
func Edit(context *gin.Context) {
	// id
	id, err := common.Param[int64](context, "id")
	if err != nil {
		FileDefaultEditPage(context, typ_api.Note{}, err)
		return
	}

	// query
	f, count, err := common.DbQry[typ_api.Note](context, "SELECT `id`, `pid`, `name`, `type`, `size`, `add_time`, `upd_time` FROM `note` WHERE `id` = ?", id)
	if err != nil || count == 0 {
		FileDefaultEditPage(context, f, err)
		return
	}

	// type
	switch typ_ft.ExtNameOf(f.Type) {
	// markdown
	case typ_ft.FtMd:
		FileMdEditPage(context, f)

	// default
	default:
		FileDefaultEditPage(context, f, err)
	}
}

// FileDefaultEditPage 默认文件修改页
func FileDefaultEditPage(context *gin.Context, note typ_api.Note, err error) {
	resp := typ_resp.Resp[typ_api.Note]{
		Msg:  util_str.TypeToStr(err),
		Data: note,
	}
	common.HtmlOk(context, "note/default/edit.html", resp)
}

// FileMdEditPage md文件修改页
func FileMdEditPage(context *gin.Context, note typ_api.Note) {
	html := func(content string, err any) {
		resp := typ_resp.Resp[map[string]any]{
			Msg: util_str.TypeToStr(err),
			Data: map[string]any{
				"note":    note,
				"content": content,
			},
		}

		common.HtmlOk(context, "note/md/edit.html", resp)
	}

	// read
	buf, err := FileRead(context, note)
	content := ""
	if err == nil {
		content = string(buf)
	}

	html(content, err)
}

// UpdContent 修改文件内容
func UpdContent(context *gin.Context) {
	json := func(err error) {
		if err != nil {
			common.JsonBadRequest(context, typ_resp.Resp[any]{Msg: util_str.TypeToStr(err)})
			return
		}

		common.JsonOk(context, typ_resp.Resp[any]{})
	}

	// id
	id, err := common.PostForm[int64](context, "id")
	if err != nil {
		json(err)
		return
	}
	//log.Println("id", id)

	// f
	f, count, err := common.DbQry[typ_api.Note](context, "SELECT `id`, `pid`, `name`, `type`, `size`, `add_time`, `upd_time` FROM `note` WHERE `del` = 0 AND `id` = ?", id)
	if count == 0 || typ_ft.ExtNameOf(f.Type) != typ_ft.FtMd {
		json(nil)
		return
	}

	// content
	content, err := common.PostForm[string](context, "content")
	if err != nil {
		json(err)
		return
	}
	//log.Println("content", content)

	// os file
	fPath, err := Path(context, f)
	if err != nil {
		json(err)
		return
	}
	pFile, err := os.OpenFile(fPath,
		os.O_WRONLY|os.O_TRUNC, // 只写（O_WRONLY） & 清空文件（O_TRUNC）
		0666)
	if err != nil {
		json(err)
		return
	}
	defer pFile.Close()

	// write
	pWriter := bufio.NewWriter(pFile)
	pWriter.WriteString(content)
	pWriter.Flush()

	// file info
	fInfo, err := pFile.Stat()
	if err != nil {
		json(err)
		return
	}

	size := fInfo.Size()

	// update
	_, err = common.DbUpd(context, "UPDATE `note` SET `size` = ?, `upd_time` = ? WHERE id = ?", size, time.Now().Unix(), id)
	if err != nil {
		json(err)
		return
	}

	json(nil)
	return
}

// FileRead 读取文件
func FileRead(context *gin.Context, f typ_api.Note) ([]byte, error) {
	// file path
	fPath, err := Path(context, f)
	if err != nil {
		return nil, err
	}

	// open file
	pFile, err := os.Open(fPath)
	if err != nil {
		return nil, err
	}
	defer pFile.Close()

	// read file
	buf, err := io.ReadAll(pFile)
	return buf, err
}

// Path 获取文件物理路径
func Path(context *gin.Context, note typ_api.Note) (string, error) {
	// dir
	dataDir := common.DataDir(context)
	fDir := fmt.Sprintf("%s%s%s%s%s", dataDir, util_os.FileSeparator, "note", util_os.FileSeparator, note.Type)
	if !util_os.IsExist(fDir) {
		err := util_os.MkDir(fDir)
		if err != nil {
			return "", err
		}
	}

	// path
	return fmt.Sprintf("%s%s%d", fDir, util_os.FileSeparator, note.Id), nil
}

func DbQry(context *gin.Context, id int64) (typ_api.Note, int64, error) {
	return common.DbQry[typ_api.Note](context, "SELECT `id`, `pid`, `name`, `type`, `size`, `add_time`, `upd_time` FROM `note` WHERE `del` = 0 AND `id` = ?", id)
}