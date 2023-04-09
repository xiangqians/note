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
	typ_page "note/src/typ/page"
	typ_resp "note/src/typ/resp"
	util_os "note/src/util/os"
	util_str "note/src/util/str"
	util_time "note/src/util/time"
	"os"
	"strings"
	"time"
)

// Del 删除文件
func Del(context *gin.Context) {
	redirect := func(pid int64, err any) {
		resp := typ_resp.Resp[any]{Msg: util_str.TypeToStr(err)}
		common.Redirect(context, fmt.Sprintf("/note/list?pid=%d", pid), resp)
	}

	// id
	id, err := common.Param[int64](context, "id")
	if err != nil {
		redirect(0, err)
		return
	}

	// note
	note, _, err := DbQry(context, id, false)
	pid := note.Pid
	if err != nil {
		redirect(pid, err)
		return
	}

	// 如果是目录则校验目录下是否有子文件
	if typ_ft.ExtNameOf(note.Type) == typ_ft.FtD {
		var count int64
		count, err = DbCount(context, id)
		if err != nil {
			redirect(pid, err)
			return
		}

		if count != 0 {
			redirect(pid, errors.New(i18n.MustGetMessage("i18n.cannotDelNonEmptyDir")))
			return
		}
	}

	// delete
	_, err = common.DbDel(context, "UPDATE `note` SET `del` = 1, `upd_time` = ? WHERE `id` = ?", util_time.NowUnix(), id)

	// redirect
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
		var note typ_api.Note
		var count int64
		note, count, err = DbQry(context, dstId, false)
		if err != nil || count == 0 || typ_ft.FtD != typ_ft.ExtNameOf(note.Type) { // 拷贝到目标类型必须是目录
			redirect(dstId, err)
			return
		}
	}

	// update
	_, err = common.DbUpd(context, "UPDATE `note` SET `pid` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `pid` <> ?",
		dstId,
		util_time.NowUnix(),
		srcId,
		dstId)

	// redirect
	redirect(dstId, err)
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

func Get(context *gin.Context) {
	// id
	id, err := common.Param[int64](context, "id")
	if err != nil {
		log.Println(err)
		return
	}

	// note
	note, count, err := DbQry(context, id, false)
	if err != nil || count == 0 {
		log.Println(err)
		return
	}

	// 排除目录
	if typ_ft.FtD == typ_ft.ExtNameOf(note.Type) {
		return
	}

	// path
	path, err := Path(context, note)
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

	context.File(path)

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
	err = common.VerifyName(note.Name)
	if err != nil {
		redirect(pid, err)
		return
	}

	// update
	_, err = common.DbUpd(context, "UPDATE `note` SET `name` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `name` <> ?", note.Name, util_time.NowUnix(), note.Id, note.Name)

	redirect(pid, err)
	return
}

// Upload 上传文件
func Upload(context *gin.Context) {
	// method
	method := context.Request.Method

	// redirect
	redirect := func(id int64, pid int64, err any) {
		resp := typ_resp.Resp[any]{
			Msg: util_str.TypeToStr(err),
		}
		switch method {
		case http.MethodPost:
			common.Redirect(context, fmt.Sprintf("/note/list?pid=%d", pid), resp)

		case http.MethodPut:
			common.Redirect(context, fmt.Sprintf("/note/%d/view", id), resp)
		}
	}

	var id int64
	var pid int64
	var err error

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

	// file name
	fn := fh.Filename
	err = common.VerifyName(fn)
	if err != nil {
		redirect(id, pid, err)
		return
	}

	// file type
	// 校验文件类型，只支持上传 html/pdf/zip
	contentType := fh.Header.Get("Content-Type")
	ft := typ_ft.ContentTypeOf(contentType)
	if ft == typ_ft.FtUnk || !(ft == typ_ft.FtHtml || ft == typ_ft.FtPdf || ft == typ_ft.FtZip) {
		redirect(id, pid, fmt.Sprintf("%s: %s", errors.New(i18n.MustGetMessage("i18n.fileTypeUnsupported")), contentType))
		return
	}

	// size
	fs := fh.Size

	// 原笔记信息
	var oldNote typ_api.Note

	// 校验 id 或 pid
	switch method {
	case http.MethodPost:
		// 校验 pid 是否存在
		if pid != 0 {
			var note typ_api.Note
			var count int64
			note, count, err = DbQry(context, pid, false)
			if err != nil || count == 0 || typ_ft.ExtNameOf(note.Type) != typ_ft.FtD { // 父节点必须是目录
				redirect(id, pid, err)
				return
			}
		}

	case http.MethodPut:
		// 校验 id 是否存在
		var count int64
		oldNote, count, err = DbQry(context, id, false)
		if err != nil || count == 0 {
			redirect(id, pid, err)
			return
		}
	}

	// 操作数据库
	switch method {
	case http.MethodPost:
		id, err = common.DbAdd(context, "INSERT INTO `note` (`pid`, `name`, `type`, `size`, `add_time`) VALUES (?, ?, ?, ?, ?)",
			pid, fn, ft, fs, util_time.NowUnix())

	case http.MethodPut:
		_, err = common.DbUpd(context, "UPDATE `note` SET `name` = ?, `type` = ?, `size` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ?",
			fn, ft, fs, util_time.NowUnix(), id)
	}
	if err != nil {
		redirect(id, pid, err)
		return
	}

	// path
	note := typ_api.Note{}
	note.Id = id
	note.Type = string(ft)
	path, err := Path(context, note)
	if err != nil {
		redirect(id, pid, err)
		return
	}

	// 清空文件
	if method == http.MethodPut && util_os.IsExist(path) {
		var file *os.File
		file, err = os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			redirect(id, pid, err)
			return
		}
		file.Close()
	}

	// 保存文件
	err = context.SaveUploadedFile(fh, path)

	// 保存文件成功时，判断如果重传不是同一个文件类型，则删除之前文件
	if err == nil && method == http.MethodPut && oldNote.Type != note.Type {
		var oldPath string
		oldPath, err = Path(context, oldNote)
		if err == nil {
			util_os.DelFile(oldPath)
		}
	}

	// redirect
	redirect(id, pid, err)
	return
}

// Add 新增文件
func Add(context *gin.Context) {
	// note
	note := typ_api.Note{}
	err := common.ShouldBind(context, &note)
	pid := note.Pid

	// redirect
	redirect := func(err any) {
		resp := typ_resp.Resp[any]{Msg: util_str.TypeToStr(err)}
		common.Redirect(context, fmt.Sprintf("/note/list?pid=%d", pid), resp)
	}

	// note err ?
	if err != nil {
		redirect(err)
		return
	}

	// name
	note.Name = strings.TrimSpace(note.Name)
	err = common.VerifyName(note.Name)
	if err != nil {
		redirect(err)
		return
	}

	// 校验文件类型
	// 只支持添加 目录 和 md文件
	ft := typ_ft.ExtNameOf(strings.TrimSpace(note.Type))
	if !(ft == typ_ft.FtD || ft == typ_ft.FtMd) {
		redirect(fmt.Sprintf("%s: %s", errors.New(i18n.MustGetMessage("i18n.fileTypeUnsupported")), note.Type))
		return
	}
	note.Type = string(ft)

	// add
	id, err := common.DbAdd(context, "INSERT INTO `note` (`pid`, `name`, `type`, `add_time`) VALUES (?, ?, ?, ?)", note.Pid, note.Name, note.Type, util_time.NowUnix())
	if err != nil {
		redirect(err)
		return
	}
	note.Id = id

	// 如果不是目录，则创建物理文件
	if ft != typ_ft.FtD {
		// path
		var path string
		path, err = Path(context, note)
		if err != nil {
			redirect(err)
			return
		}

		// create
		var file *os.File
		file, err = os.Create(path)
		if err != nil {
			redirect(err)
			return
		}
		defer file.Close()
	}

	redirect(err)
	return
}

// List 文件列表页面
func List(context *gin.Context) {
	html := func(pnote typ_api.Note, notes []typ_api.Note, types []string, err error) {
		resp := typ_resp.Resp[map[string]any]{
			Msg: util_str.TypeToStr(err),
			Data: map[string]any{
				"pnote": pnote,
				"notes": notes,
				"types": types,
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

	// type
	t, err := common.Query[string](context, "type")
	t = strings.TrimSpace(t)
	//log.Printf("t = %s\n", t)
	ft := typ_ft.ExtNameOf(t)
	if ft == typ_ft.FtUnk {
		t = ""
	}

	// pnote
	var pnote typ_api.Note
	if pid < 0 {
		pnote.Path = ""
		pnote.PathLink = ""

	} else if pid == 0 {
		pnote.Path = "/"
		pnote.PathLink = "/"

	} else {
		sql, args := DbQrySql(typ_api.Note{Abs: typ_api.Abs{Id: pid}, Pid: -1}, true)
		sql += "LIMIT 1"
		var count int64
		pnote, count, err = common.DbQry[typ_api.Note](context, sql, args...)
		if err != nil || count == 0 {
			html(pnote, nil, nil, err)
			return
		}

		if pnote.Path != "/" {
			pnote.Path += "/"
		}
		pnote.Path += fmt.Sprintf("%d:%s", pnote.Id, pnote.Name)
		InitPath(&pnote)
	}

	// 查询
	notes, err := DbList(context, typ_api.Note{
		Pid:  pid,
		Name: name,
		Type: t,
	})

	types, count, err := common.DbQry[[]string](context, "SELECT DISTINCT(`type`) FROM `note` WHERE `del` = 0")
	if err != nil || count == 0 {
		types = nil
	}

	// html
	html(pnote, notes, types, err)
	return
}

// FileRead 读取文件
func FileRead(context *gin.Context, note typ_api.Note) ([]byte, error) {
	// file path
	path, err := Path(context, note)
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

// HistPath 获取笔记历史记录物理路径
func HistPath(context *gin.Context, note typ_api.Note) (string, error) {
	// dir
	dataDir := common.DataDir(context)
	noteDir := fmt.Sprintf("%s%s%s%s%s%s%s", dataDir,
		util_os.FileSeparator, "hist",
		util_os.FileSeparator, "note",
		util_os.FileSeparator, note.Type)
	if !util_os.IsExist(noteDir) {
		err := util_os.MkDir(noteDir)
		if err != nil {
			return "", err
		}
	}

	// file name
	time := note.AddTime
	if note.AddTime < note.UpdTime {
		time = note.UpdTime
	}
	name := fmt.Sprintf("%d_%d", note.Id, time)

	// path
	return fmt.Sprintf("%s%s%s", noteDir, util_os.FileSeparator, name), nil
}

// Path 获取文件物理路径
func Path(context *gin.Context, note typ_api.Note) (string, error) {
	// dir
	dataDir := common.DataDir(context)
	noteDir := fmt.Sprintf("%s%s%s%s%s", dataDir,
		util_os.FileSeparator, "note",
		util_os.FileSeparator, note.Type)
	if !util_os.IsExist(noteDir) {
		err := util_os.MkDir(noteDir)
		if err != nil {
			return "", err
		}
	}

	// file name
	name := fmt.Sprintf("%d", note.Id)

	// path
	return fmt.Sprintf("%s%s%s", noteDir, util_os.FileSeparator, name), nil
}

func DbCount(context *gin.Context, pid int64) (int64, error) {
	count, _, err := common.DbQry[int64](context, "SELECT COUNT(1) FROM `note` WHERE `pid` = ?", pid)
	return count, err
}

func DbPage(context *gin.Context, note typ_api.Note) (typ_page.Page[typ_api.Note], error) {
	req, _ := common.PageReq(context)
	path := true
	sql, args := DbQrySql(note, path)
	page, err := common.DbPage[typ_api.Note](context, req, sql, args...)
	if path && err == nil {
		data := page.Data
		if data != nil && len(data) > 0 {
			for i, l := 0, len(data); i < l; i++ {
				InitPath(&data[i])
			}
		}
	}
	return page, err
}

func DbList(context *gin.Context, note typ_api.Note) ([]typ_api.Note, error) {
	// 查询
	path := false
	if note.Pid == -1 {
		path = true
	}
	sql, args := DbQrySql(note, path)
	sql += "LIMIT 10000"
	notes, count, err := common.DbQry[[]typ_api.Note](context, sql, args...)
	if err != nil || count == 0 {
		notes = nil
	}

	if path && err == nil && count > 0 {
		for i, l := 0, len(notes); i < l; i++ {
			InitPath(&notes[i])
		}
	}

	return notes, err
}

func DbQry(context *gin.Context, id int64, path bool) (typ_api.Note, int64, error) {
	sql, args := DbQrySql(typ_api.Note{Abs: typ_api.Abs{Id: id}, Pid: -1}, path)
	sql += "LIMIT 1"
	note, count, err := common.DbQry[typ_api.Note](context, sql, args...)
	if path && err == nil && count > 0 {
		InitPath(&note)
	}
	return note, count, err
}

// InitPath 初始化 path & pathLink
func InitPath(note *typ_api.Note) {
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
		pathLink := fmt.Sprintf("<a href=\"/note/list?pid=%s&t=%d\">%s</a>\n", eArr[0], util_time.NowUnix(), eArr[1])
		pathLinkArr = append(pathLinkArr, pathLink)
	}
	(*note).Path = strings.Join(pathArr, "/")
	(*note).PathLink = strings.Join(pathLinkArr, "/")
}

func DbQrySql(note typ_api.Note, path bool) (string, []any) {
	args := make([]any, 0, 1)
	sql := "SELECT n.`id`, n.`pid`, n.`name`, n.`type`, n.`size`, n.`hist`, n.`hist_size`, n.`add_time`, n.`upd_time` "
	if path {
		sql += ", CASE WHEN n.`pid` = 0 THEN  '/' ELSE " +
			"( (CASE WHEN pn10.`id` IS NULL THEN '' ELSE '/' || pn10.`id` || ':' ||pn10.`name` END) " +
			"|| (CASE WHEN pn9.`id` IS NULL THEN '' ELSE '/' || pn9.`id` || ':' || pn9.`name` END) " +
			"|| (CASE WHEN pn8.`id` IS NULL THEN '' ELSE '/' || pn8.`id` || ':' || pn8.`name` END) " +
			"|| (CASE WHEN pn7.`id` IS NULL THEN '' ELSE '/' || pn7.`id` || ':' || pn7.`name` END) " +
			"|| (CASE WHEN pn6.`id` IS NULL THEN '' ELSE '/' || pn6.`id` || ':' || pn6.`name` END) " +
			"|| (CASE WHEN pn5.`id` IS NULL THEN '' ELSE '/' || pn5.`id` || ':' || pn5.`name` END) " +
			"|| (CASE WHEN pn4.`id` IS NULL THEN '' ELSE '/' || pn4.`id` || ':' || pn4.`name` END) " +
			"|| (CASE WHEN pn3.`id` IS NULL THEN '' ELSE '/' || pn3.`id` || ':' || pn3.`name` END) " +
			"|| (CASE WHEN pn2.`id` IS NULL THEN '' ELSE '/' || pn2.`id` || ':' || pn2.`name` END) " +
			"|| (CASE WHEN pn1.`id` IS NULL THEN '' ELSE '/' || pn1.`id` || ':' || pn1.`name` END)) " +
			"END AS 'path' "
	}
	sql += "FROM `note` n "
	if path {
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

	// del
	sql += "WHERE n.`del` = ? "
	args = append(args, note.Del)

	// id
	if note.Id > 0 {
		sql += "AND n.`id` = ? "
		args = append(args, note.Id)
	}

	// pid
	if note.Pid >= 0 {
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
	sql += "ORDER BY (CASE WHEN n.`upd_time` > n.`add_time` THEN n.`upd_time` ELSE n.`add_time` END) DESC "

	return sql, args
}
