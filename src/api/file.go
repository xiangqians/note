// file
// @author xiangqian
// @date 17:50 2023/02/04
package api

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
	"note/src/typ"
	"note/src/util"
	"os"
	"strings"
	"time"
)

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
	id, err := DbAdd(pContext, "INSERT INTO `file` (`pid`, `name`, `type`, `add_time`) VALUES (?, ?, ?, ?)", f.Pid, f.Name, f.Type, time.Now().Unix())
	if err != nil {
		redirect(pid, err)
		return
	}
	f.Id = id

	// 如果不是目录，则创建物理文件
	if fType != typ.FileTypeD {
		// file path
		fp, err := FilePath(pContext, f)
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

// FileUpload 上传文件
func FileUpload(pContext *gin.Context) {
	method := pContext.Request.Method
	redirect := func(id int64, pid int64, msg any) {
		switch method {
		case http.MethodPost:
			Redirect(pContext, fmt.Sprintf("/?id=%d", pid), nil, msg)

		case http.MethodPut:
			Redirect(pContext, fmt.Sprintf("/file/%d/editpage", id), nil, msg)
		}
	}

	var id int64
	var pid int64
	var err error

	// 获取 id 或 pid
	switch method {
	// 上传文件，需要pid
	case http.MethodPost:
		pid, err = PostForm[int64](pContext, "pid")

	// 重传文件，需要id
	case http.MethodPut:
		id, err = PostForm[int64](pContext, "id")
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
	if ft == typ.FileTypeUnk || !(ft == typ.FileTypePdf || ft == typ.FileTypeZip) {
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
			f, count, err := DbQry[typ.File](pContext, "SELECT f.id, f.pid, f.`name`, f.`type`, f.`size`, f.`add_time`, f.`upd_time` FROM `file` f WHERE f.`del` = 0 AND f.`id` = ?", pid)
			if err != nil || count == 0 || typ.FileTypeOf(f.Type) != typ.FileTypeD {
				redirect(id, pid, err)
				return
			}
		}

	case http.MethodPut:
		// 校验 id 是否存在
		f, count, err := DbQry[typ.File](pContext, "SELECT f.id, f.pid, f.`name`, f.`type`, f.`size`, f.`add_time`, f.`upd_time` FROM `file` f WHERE f.`del` = 0 AND f.`id` = ?", id)
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
		id, err = DbAdd(pContext, "INSERT INTO `file` (`pid`, `name`, `type`, `size`, `add_time`) VALUES (?, ?, ?, ?, ?)",
			pid, fn, ft, fs, time.Now().Unix())

	case http.MethodPut:
		_, err = DbUpd(pContext, "UPDATE `file` SET `name` = ?, `type` = ?, `size` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ?",
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
	fp, err := FilePath(pContext, f)
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

// FileReUpload 重新上传文件
func FileReUpload(pContext *gin.Context) {
	pContext.Request.Method = http.MethodPut
	FileUpload(pContext)
}

// FileUpdName 文件重命名
func FileUpdName(pContext *gin.Context) {
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

	// name
	f.Name = strings.TrimSpace(f.Name)
	err = util.VerifyFileName(f.Name)
	if err != nil {
		redirect(pid, err)
		return
	}

	// update
	_, err = DbUpd(pContext, "UPDATE `file` SET `name` = ?, `upd_time` = ? WHERE `id` = ? AND `name` <> ?", f.Name, time.Now().Unix(), f.Id, f.Name)

	redirect(pid, err)
	return
}

// FileCut 剪切文件
func FileCut(pContext *gin.Context) {
	redirect := func(id int64, msg any) {
		Redirect(pContext, fmt.Sprintf("/?id=%d", id), nil, msg)
	}

	// dst id
	dstId, err := Param[int64](pContext, "dstId")
	if err != nil {
		redirect(dstId, err)
		return
	}

	// src id
	srcId, err := Param[int64](pContext, "srcId")
	if err != nil {
		redirect(dstId, err)
		return
	}

	// dst
	if dstId != 0 {
		f, _, err := DbQry[typ.File](pContext, "SELECT f.id, f.pid, f.`name`, f.`type`, f.`size`, f.`add_time`, f.`upd_time` FROM `file` f WHERE f.`del` = 0 AND f.`id` = ?", dstId)
		if err != nil || typ.FileTypeD != typ.FileTypeOf(f.Type) {
			redirect(dstId, err)
			return
		}
	}

	// update
	_, err = DbUpd(pContext, "UPDATE `file` SET `pid` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `pid` <> ?",
		dstId,
		time.Now().Unix(),
		srcId,
		dstId)

	redirect(dstId, err)
	return
}

// FileDel 删除文件
func FileDel(pContext *gin.Context) {
	redirect := func(id int64, msg any) {
		Redirect(pContext, fmt.Sprintf("/?id=%d", id), nil, msg)
	}

	// id
	id, err := Param[int64](pContext, "id")
	if err != nil {
		redirect(0, err)
		return
	}

	// pid
	pid, _, err := DbQry[int64](pContext, "SELECT f.pid FROM `file` f WHERE f.del = 0 AND f.id = ?", id)
	if err != nil {
		redirect(pid, err)
		return
	}

	// update
	_, err = DbDel(pContext, "UPDATE `file` SET `del` = 1, `upd_time` = ? WHERE `id` = ?", time.Now().Unix(), id)

	redirect(pid, err)
	return
}

func FileView(pContext *gin.Context) {
	// id
	id, err := Param[int64](pContext, "id")
	if err != nil {
		log.Println(err)
		return
	}

	// f
	f, count, err := DbQry[typ.File](pContext, "SELECT f.`id`, f.`pid`, f.`name`, f.`type`, f.`size`, f.`add_time`, f.`upd_time` FROM `file` f WHERE f.`del` = 0 AND f.`id` = ?", id)
	if err != nil || count == 0 {
		log.Println(err)
		return
	}

	// 排除目录
	if typ.FileTypeD == typ.FileTypeOf(f.Type) {
		return
	}

	// path
	fPath, err := FilePath(pContext, f)
	if err != nil {
		log.Println(err)
		return
	}

	// read
	buf, err := os.ReadFile(fPath)
	if err != nil {
		log.Println(err)
		return
	}

	pContext.Writer.Write(buf)
	return
}

// FileViewPage 查看文件页面
func FileViewPage(pContext *gin.Context) {
	// id
	id, err := Param[int64](pContext, "id")
	if err != nil {
		FileDefaultViewPage(pContext, typ.File{}, err)
		return
	}

	// query
	f, count, err := DbQry[typ.File](pContext, "SELECT f.`id`, f.`pid`, f.`name`, f.`type`, f.`size`, f.`add_time`, f.`upd_time` FROM `file` f WHERE f.`del` = 0 AND f.`id` = ?", id)
	if err != nil || count == 0 {
		FileDefaultViewPage(pContext, f, err)
		return
	}

	// type
	switch typ.FileTypeOf(f.Type) {
	// markdown
	case typ.FileTypeMd:
		FileMdViewPage(pContext, f)

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
	Html(pContext, "file/default/view.html", gin.H{"f": f}, err)
}

// FileMdViewPage 查看md文件
// https://github.com/russross/blackfriday
// https://pkg.go.dev/github.com/russross/blackfriday/v2
func FileMdViewPage(pContext *gin.Context, f typ.File) {
	html := func(html string, msg any) {
		Html(pContext, "file/md/view.html", gin.H{"f": f, "html": html}, msg)
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

// FilePdfViewPage 查看pdf文件
func FilePdfViewPage(pContext *gin.Context, f typ.File) {
	Html(pContext, "file/pdf/view.html", gin.H{"f": f}, nil)
}

// FileEditPage 文件修改页
func FileEditPage(pContext *gin.Context) {
	// id
	id, err := Param[int64](pContext, "id")
	if err != nil {
		FileDefaultEditPage(pContext, typ.File{}, err)
		return
	}

	// query
	f, count, err := DbQry[typ.File](pContext, "SELECT f.id, f.pid, f.`name`, f.`type`, f.`size`, f.add_time, f.upd_time FROM `file` f WHERE f.id = ?", id)
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
	Html(pContext, "file/default/edit.html", gin.H{"f": f}, err)
}

// FileMdEditPage md文件修改页
func FileMdEditPage(pContext *gin.Context, f typ.File) {
	html := func(content string, msg any) {
		Html(pContext, "file/md/edit.html", gin.H{"f": f, "content": content}, msg)
	}

	// read
	buf, err := FileRead(pContext, f)
	content := ""
	if err == nil {
		content = string(buf)
	}

	html(content, err)
}

// FileUpdContent 修改文件内容
func FileUpdContent(pContext *gin.Context) {
	json := func(err error) {
		if err != nil {
			pContext.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}

		pContext.JSON(http.StatusOK, nil)
	}

	// id
	id, err := PostForm[int64](pContext, "id")
	if err != nil {
		json(err)
		return
	}
	//log.Println("id", id)

	// f
	f, count, err := DbQry[typ.File](pContext, "SELECT f.id, f.pid, f.`name`, f.`type`, f.`size`, f.add_time, f.upd_time FROM `file` f WHERE f.del = 0 AND f.id = ?", id)
	if count == 0 || typ.FileTypeOf(f.Type) != typ.FileTypeMd {
		json(nil)
		return
	}

	// content
	content, err := PostForm[string](pContext, "content")
	if err != nil {
		json(err)
		return
	}
	//log.Println("content", content)

	// os file
	fPath, err := FilePath(pContext, f)
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
	_, err = DbUpd(pContext, "UPDATE `file` SET `size` = ?, `upd_time` = ? WHERE id = ?", size, time.Now().Unix(), id)
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
	fPath, err := FilePath(pContext, f)
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

// FilePath 获取文件物理路径
func FilePath(pContext *gin.Context, f typ.File) (string, error) {
	// dir
	dataDir := DataDir(pContext)
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
