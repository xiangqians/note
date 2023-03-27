// file
// @author xiangqian
// @date 17:50 2023/02/04
package file

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
	"note/src/typ"
	"note/src/util"
	"os"
	"strings"
	"time"
)

// List 文件列表页面
func List(pContext *gin.Context) {
	html := func(pfile typ.File, files []typ.File, err error) {
		resp := typ.Resp[map[string]any]{
			Msg: util.TypeAsStr(err),
			Data: map[string]any{
				"pfile": pfile,
				"files": files,
			},
		}
		common.HtmlOkNew(pContext, "file/list.html", resp)
	}

	// id
	id, err := common.Query[int64](pContext, "id")
	//log.Printf("id = %d\n", id)

	// name
	name, err := common.Query[string](pContext, "name")
	name = strings.TrimSpace(name)
	//log.Printf("name = %s\n", name)

	// pf
	var pfile typ.File
	if id < 0 {
		pfile.Path = ""
		pfile.PathLink = ""

	} else if id == 0 {
		pfile.Path = "/"
		pfile.PathLink = "/"

	} else {
		sql := "SELECT f1.id, f1.pid, f1.`name`, f1.`type`, f1.`size`, f1.add_time, f1.upd_time, " +
			"( (CASE WHEN f10.`id` IS NULL THEN '' ELSE '/' ||f10.`id` || ':' ||f10.`name` END) " +
			"|| (CASE WHEN f9.`id` IS NULL THEN '' ELSE '/' || f9.`id` || ':' || f9.`name` END) " +
			"|| (CASE WHEN f8.`id` IS NULL THEN '' ELSE '/' || f8.`id` || ':' || f8.`name` END) " +
			"|| (CASE WHEN f7.`id` IS NULL THEN '' ELSE '/' || f7.`id` || ':' || f7.`name` END) " +
			"|| (CASE WHEN f6.`id` IS NULL THEN '' ELSE '/' || f6.`id` || ':' || f6.`name` END) " +
			"|| (CASE WHEN f5.`id` IS NULL THEN '' ELSE '/' || f5.`id` || ':' || f5.`name` END) " +
			"|| (CASE WHEN f4.`id` IS NULL THEN '' ELSE '/' || f4.`id` || ':' || f4.`name` END) " +
			"|| (CASE WHEN f3.`id` IS NULL THEN '' ELSE '/' || f3.`id` || ':' || f3.`name` END) " +
			"|| (CASE WHEN f2.`id` IS NULL THEN '' ELSE '/' || f2.`id` || ':' || f2.`name` END) " +
			"|| (CASE WHEN f1.`id` IS NULL THEN '' ELSE '/' || f1.`id` || ':' || f1.`name` END))  AS 'path' " +
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
		pfile, _, err = common.DbQry[typ.File](pContext, sql, id)
		if err != nil {
			html(pfile, nil, err)
			return
		}

		pathArr := strings.Split(pfile.Path, "/")
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

			pathLink := fmt.Sprintf("<a href=\"/?id=%s\">%s</a>\n", vArr[0], vArr[1])
			pathLinkArr = append(pathLinkArr, pathLink)
		}
		pfile.Path = strings.Join(pathArr, "/")
		pfile.PathLink = strings.Join(pathLinkArr, "/")
	}

	// 查询
	args := make([]any, 0, 2)
	var files []typ.File = nil
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
	sql += "ORDER BY f.`type`, (CASE WHEN f.`upd_time` > f.`add_time` THEN f.`upd_time` ELSE f.`add_time` END) DESC "
	if id < 0 {
		sql += "LIMIT 10000"
	}
	files, count, err = common.DbQry[[]typ.File](pContext, sql, args...)
	if err != nil || count == 0 {
		files = nil
	}

	html(pfile, files, err)
	return
}

// Add 新增文件
func Add(pContext *gin.Context) {
	redirect := func(id int64, err any) {
		resp := typ.Resp[any]{Msg: util.TypeAsStr(err)}
		common.RedirectNew(pContext, fmt.Sprintf("/file/list?id=%d", id), resp)
	}

	// file
	f := typ.File{}
	err := common.ShouldBind(pContext, &f)
	pid := f.Pid
	if err != nil {
		redirect(pid, err)
		return
	}

	// name
	f.Name = strings.TrimSpace(f.Name)
	err = util.VerifyFileName(f.Name)
	if err != nil {
		redirect(pid, err)
		return
	}

	// 校验文件类型
	// 只支持添加 目录 和 md文件
	fType := typ.FileTypeOf(strings.TrimSpace(f.Type))
	if !(fType == typ.FileTypeD || fType == typ.FileTypeMd) {
		redirect(pid, fmt.Sprintf("%s, %s", errors.New(i18n.MustGetMessage("i18n.fileTypeUnsupported")), f.Type))
		return
	}
	f.Type = string(fType)

	// add
	id, err := common.DbAdd(pContext, "INSERT INTO `file` (`pid`, `name`, `type`, `add_time`) VALUES (?, ?, ?, ?)", f.Pid, f.Name, f.Type, time.Now().Unix())
	if err != nil {
		redirect(pid, err)
		return
	}
	f.Id = id

	// 如果不是目录，则创建物理文件
	if fType != typ.FileTypeD {
		// file path
		fp, err := Path(pContext, f)
		if err != nil {
			redirect(pid, err)
			return
		}

		// create
		pFile, err := os.Create(fp)
		if err != nil {
			redirect(pid, err)
			return
		}
		defer pFile.Close()
	}

	redirect(pid, err)
	return
}

// Upload 上传文件
func Upload(pContext *gin.Context) {
	method := pContext.Request.Method
	redirect := func(id int64, pid int64, msg any) {
		switch method {
		case http.MethodPost:
			common.Redirect(pContext, fmt.Sprintf("/file/list?id=%d", pid), nil, msg)

		case http.MethodPut:
			common.Redirect(pContext, fmt.Sprintf("/file/%d/edit", id), nil, msg)
		}
	}

	var id int64
	var pid int64
	var err error

	// 获取 id 或 pid
	switch method {
	// 上传文件，需要pid
	case http.MethodPost:
		pid, err = common.PostForm[int64](pContext, "pid")

	// 重传文件，需要id
	case http.MethodPut:
		id, err = common.PostForm[int64](pContext, "id")
	}
	if err != nil {
		redirect(id, pid, err)
		return
	}

	// fh
	fh, err := pContext.FormFile("file")
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
	// 校验文件类型，只支持上传 pdf 和 zip
	ftStr := ""
	index := strings.LastIndex(fn, ".")
	if index > 0 {
		ftStr = fn[index+1:]
	}
	ft := typ.FileTypeOf(strings.TrimSpace(ftStr))
	if ft == typ.FileTypeUnk || !(ft == typ.FileTypeHtml || ft == typ.FileTypePdf || ft == typ.FileTypeZip) {
		redirect(id, pid, fmt.Sprintf("%s, %s", errors.New(i18n.MustGetMessage("i18n.fileTypeUnsupported")), ftStr))
		return
	}

	// size
	fs := fh.Size

	// 校验 id 或 pid
	switch method {
	case http.MethodPost:
		// 校验 pid 是否存在
		if pid != 0 {
			f, count, err := common.DbQry[typ.File](pContext, "SELECT f.id, f.pid, f.`name`, f.`type`, f.`size`, f.`add_time`, f.`upd_time` FROM `file` f WHERE f.`del` = 0 AND f.`id` = ?", pid)
			if err != nil || count == 0 || typ.FileTypeOf(f.Type) != typ.FileTypeD {
				redirect(id, pid, err)
				return
			}
		}

	case http.MethodPut:
		// 校验 id 是否存在
		f, count, err := common.DbQry[typ.File](pContext, "SELECT f.id, f.pid, f.`name`, f.`type`, f.`size`, f.`add_time`, f.`upd_time` FROM `file` f WHERE f.`del` = 0 AND f.`id` = ?", id)
		if err != nil || count == 0 {
			redirect(id, pid, err)
			return
		}

		if ft != typ.FileTypeOf(f.Type) {
			redirect(id, pid, "重传必须是相同文件类型")
			return
		}
	}

	// 操作数据库
	switch method {
	case http.MethodPost:
		id, err = common.DbAdd(pContext, "INSERT INTO `file` (`pid`, `name`, `type`, `size`, `add_time`) VALUES (?, ?, ?, ?, ?)",
			pid, fn, ft, fs, time.Now().Unix())

	case http.MethodPut:
		_, err = common.DbUpd(pContext, "UPDATE `file` SET `name` = ?, `type` = ?, `size` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ?",
			fn, ft, fs, time.Now().Unix(), id)
	}
	if err != nil {
		redirect(id, pid, err)
		return
	}

	// path
	f := typ.File{}
	f.Id = id
	f.Type = string(ft)
	fp, err := Path(pContext, f)
	if err != nil {
		redirect(id, pid, err)
		return
	}

	// 清空文件
	if method == http.MethodPut && util.IsExistOfPath(fp) {
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
	err = pContext.SaveUploadedFile(fh, fp)

	redirect(id, pid, err)
	return
}

// UpdName 文件重命名
func UpdName(pContext *gin.Context) {
	redirect := func(pid int64, msg any) {
		common.Redirect(pContext, fmt.Sprintf("/file/list?id=%d", pid), nil, msg)
	}

	// file
	f := typ.File{}
	err := common.ShouldBind(pContext, &f)
	pid := f.Pid
	if err != nil {
		redirect(pid, err)
		return
	}

	// name
	f.Name = strings.TrimSpace(f.Name)
	err = util.VerifyFileName(f.Name)
	if err != nil {
		redirect(pid, err)
		return
	}

	//fType, count, err := common.DbQry[string](pContext, "SELECT `type` FROM `file` WHERE `del` = 0 AND `id` = ?", f.Id)
	//if count > 0 {
	//	name := f.Name
	//	ft := typ.FileTypeOf(fType)
	//	if ft != typ.FileTypeD && ft != typ.FileTypeUnk && !strings.HasSuffix(name, string(ft)) {
	//		name = fmt.Sprintf("%s.%s", name, string(ft))
	//	}
	//
	//	// update
	//	_, err = common.DbUpd(pContext, "UPDATE `file` SET `name` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `name` <> ?", name, time.Now().Unix(), f.Id, name)
	//}

	// update
	_, err = common.DbUpd(pContext, "UPDATE `file` SET `name` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `name` <> ?", f.Name, time.Now().Unix(), f.Id, f.Name)

	redirect(pid, err)
	return
}

// Cut 剪切文件
func Cut(pContext *gin.Context) {
	redirect := func(id int64, msg any) {
		common.Redirect(pContext, fmt.Sprintf("/?id=%d", id), nil, msg)
	}

	// dst id
	dstId, err := common.Param[int64](pContext, "dstId")
	if err != nil {
		redirect(dstId, err)
		return
	}

	// src id
	srcId, err := common.Param[int64](pContext, "srcId")
	if err != nil {
		redirect(dstId, err)
		return
	}

	// dst
	if dstId != 0 {
		f, _, err := common.DbQry[typ.File](pContext, "SELECT f.id, f.pid, f.`name`, f.`type`, f.`size`, f.`add_time`, f.`upd_time` FROM `file` f WHERE f.`del` = 0 AND f.`id` = ?", dstId)
		if err != nil || typ.FileTypeD != typ.FileTypeOf(f.Type) {
			redirect(dstId, err)
			return
		}
	}

	// update
	_, err = common.DbUpd(pContext, "UPDATE `file` SET `pid` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `pid` <> ?",
		dstId,
		time.Now().Unix(),
		srcId,
		dstId)

	redirect(dstId, err)
	return
}

// Del 删除文件
func Del(pContext *gin.Context) {
	redirect := func(id int64, msg any) {
		common.Redirect(pContext, fmt.Sprintf("/?id=%d", id), nil, msg)
	}

	// id
	id, err := common.Param[int64](pContext, "id")
	if err != nil {
		redirect(0, err)
		return
	}

	// pid
	pid, _, err := common.DbQry[int64](pContext, "SELECT f.pid FROM `file` f WHERE f.del = 0 AND f.id = ?", id)
	if err != nil {
		redirect(pid, err)
		return
	}

	// update
	_, err = common.DbDel(pContext, "UPDATE `file` SET `del` = 1, `upd_time` = ? WHERE `id` = ?", time.Now().Unix(), id)

	redirect(pid, err)
	return
}

func Get(pContext *gin.Context) {
	// id
	id, err := common.Param[int64](pContext, "id")
	if err != nil {
		log.Println(err)
		return
	}

	// f
	f, count, err := common.DbQry[typ.File](pContext, "SELECT f.`id`, f.`pid`, f.`name`, f.`type`, f.`size`, f.`add_time`, f.`upd_time` FROM `file` f WHERE f.`del` = 0 AND f.`id` = ?", id)
	if err != nil || count == 0 {
		log.Println(err)
		return
	}

	// 排除目录
	if typ.FileTypeD == typ.FileTypeOf(f.Type) {
		return
	}

	// path
	fPath, err := Path(pContext, f)
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
	writer := pContext.Writer
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
	err = util.IOCopy(pFile, pContext.Writer, 0)
	if err != nil {
		log.Println(err)
		return
	}
	*/

	pContext.File(fPath)

	return
}

// View 查看文件页面
func View(pContext *gin.Context) {
	// id
	id, err := common.Param[int64](pContext, "id")
	if err != nil {
		FileDefaultViewPage(pContext, typ.File{}, err)
		return
	}

	// query
	f, count, err := common.DbQry[typ.File](pContext, "SELECT f.`id`, f.`pid`, f.`name`, f.`type`, f.`size`, f.`add_time`, f.`upd_time` FROM `file` f WHERE f.`del` = 0 AND f.`id` = ?", id)
	if err != nil || count == 0 {
		FileDefaultViewPage(pContext, f, err)
		return
	}

	// type
	switch typ.FileTypeOf(f.Type) {
	// markdown
	case typ.FileTypeMd:
		FileMdViewPage(pContext, f)

	// html
	case typ.FileTypeHtml:
		FileHtmlViewPage(pContext, f)

	// pdf
	case typ.FileTypePdf:
		FilePdfViewPage(pContext, f)

	// default
	default:
		FileDefaultViewPage(pContext, f, err)
	}
}

// FileDefaultViewPage 默认查看文件
func FileDefaultViewPage(pContext *gin.Context, f typ.File, err error) {
	common.HtmlOk(pContext, "file/default/view.html", gin.H{"f": f}, err)
}

// FileMdViewPage 查看md文件
// https://github.com/russross/blackfriday
// https://pkg.go.dev/github.com/russross/blackfriday/v2
func FileMdViewPage(pContext *gin.Context, f typ.File) {
	html := func(html string, msg any) {
		common.HtmlOk(pContext, "file/md/view.html", gin.H{"f": f, "html": html}, msg)
	}

	// read
	buf, err := FileRead(pContext, f)
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

// FileHtmlViewPage 查看html文件
func FileHtmlViewPage(pContext *gin.Context, f typ.File) {
	html := func(html string, msg any) {
		common.HtmlOk(pContext, "file/html/view.html", gin.H{"f": f, "html": html}, msg)
	}

	// read
	buf, err := FileRead(pContext, f)
	if err != nil {
		html("", err)
		return
	}

	html(string(buf), nil)
}

// FilePdfViewPage 查看pdf文件
func FilePdfViewPage(pContext *gin.Context, f typ.File) {
	v, _ := common.Query[string](pContext, "v")
	v = strings.TrimSpace(v)
	switch v {
	case "1.0":
		// v1.0
	case "2.0":
		// v2.0
	default:
		v = "2.0"
	}
	common.HtmlOk(pContext, fmt.Sprintf("file/pdf/view_v%s.html", v), gin.H{"f": f}, nil)
}

// Edit 文件修改页
func Edit(pContext *gin.Context) {
	// id
	id, err := common.Param[int64](pContext, "id")
	if err != nil {
		FileDefaultEditPage(pContext, typ.File{}, err)
		return
	}

	// query
	f, count, err := common.DbQry[typ.File](pContext, "SELECT f.id, f.pid, f.`name`, f.`type`, f.`size`, f.add_time, f.upd_time FROM `file` f WHERE f.id = ?", id)
	if err != nil || count == 0 {
		FileDefaultEditPage(pContext, f, err)
		return
	}

	// type
	switch typ.FileTypeOf(f.Type) {
	// markdown
	case typ.FileTypeMd:
		FileMdEditPage(pContext, f)

	// default
	default:
		FileDefaultEditPage(pContext, f, err)
	}
}

// FileDefaultEditPage 默认文件修改页
func FileDefaultEditPage(pContext *gin.Context, f typ.File, err error) {
	common.HtmlOk(pContext, "file/default/edit.html", gin.H{"f": f}, err)
}

// FileMdEditPage md文件修改页
func FileMdEditPage(pContext *gin.Context, f typ.File) {
	html := func(content string, msg any) {
		common.HtmlOk(pContext, "file/md/edit.html", gin.H{"f": f, "content": content}, msg)
	}

	// read
	buf, err := FileRead(pContext, f)
	content := ""
	if err == nil {
		content = string(buf)
	}

	html(content, err)
}

// UpdContent 修改文件内容
func UpdContent(pContext *gin.Context) {
	json := func(err error) {
		if err != nil {
			pContext.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}

		pContext.JSON(http.StatusOK, nil)
	}

	// id
	id, err := common.PostForm[int64](pContext, "id")
	if err != nil {
		json(err)
		return
	}
	//log.Println("id", id)

	// f
	f, count, err := common.DbQry[typ.File](pContext, "SELECT f.id, f.pid, f.`name`, f.`type`, f.`size`, f.add_time, f.upd_time FROM `file` f WHERE f.del = 0 AND f.id = ?", id)
	if count == 0 || typ.FileTypeOf(f.Type) != typ.FileTypeMd {
		json(nil)
		return
	}

	// content
	content, err := common.PostForm[string](pContext, "content")
	if err != nil {
		json(err)
		return
	}
	//log.Println("content", content)

	// os file
	fPath, err := Path(pContext, f)
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
	_, err = common.DbUpd(pContext, "UPDATE `file` SET `size` = ?, `upd_time` = ? WHERE id = ?", size, time.Now().Unix(), id)
	if err != nil {
		json(err)
		return
	}

	json(nil)
	return
}

// FileRead 读取文件
func FileRead(pContext *gin.Context, f typ.File) ([]byte, error) {
	// file path
	fPath, err := Path(pContext, f)
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
func Path(pContext *gin.Context, f typ.File) (string, error) {
	// dir
	dataDir := common.DataDir(pContext)
	fDir := fmt.Sprintf("%s%s%s%s%s", dataDir, util.FileSeparator, "file", util.FileSeparator, f.Type)
	if !util.IsExistOfPath(fDir) {
		err := util.Mkdir(fDir)
		if err != nil {
			return "", err
		}
	}

	// path
	return fmt.Sprintf("%s%s%d", fDir, util.FileSeparator, f.Id), nil
}
