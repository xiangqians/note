// file
// @author xiangqian
// @date 17:50 2023/02/04
package api

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
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

	// 校验文件类型
	fType := typ.FileTypeOf(strings.TrimSpace(f.Type))
	if fType == typ.FileTypeUnk {
		redirect(pid, errors.New(i18n.MustGetMessage("i18n.fileTypeUnsupported")))
	}

	// add
	id, err := DbAdd(pContext, "INSERT INTO `file` (`pid`, `name`, `type`, `add_time`) VALUES (?, ?, ?, ?)", f.Pid, f.Name, fType, time.Now().Unix())
	if err != nil {
		redirect(pid, err)
		return
	}

	f.Id = id

	// 如果不是目录，则创建物理文件
	if fType != typ.FileTypeD {
		fPath := FilePath(pContext, f)
		pFile, fErr := os.Create(fPath)
		if fErr != nil {
			log.Println(fErr)
		}
		defer pFile.Close()
	}

	redirect(pid, nil)
	return
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

	// update
	_, err = DbUpd(pContext, "UPDATE `file` SET `name` = ?, `upd_time` = ? WHERE id = ?", f.Name, time.Now().Unix(), f.Id)
	if err != nil {
		redirect(pid, err)
		return
	}

	redirect(pid, nil)
	return
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
	fPath := FilePath(pContext, f)
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

// FileDel 删除文件
func FileDel(pContext *gin.Context) {
	redirect := func(id int64, msg any) {
		Redirect(pContext, fmt.Sprintf("/?id=%d", id), nil, msg)
	}

	// id
	id, _ := Param[int64](pContext, "id")

	// pid
	pid, _, _ := DbQry[int64](pContext, "SELECT f.pid FROM `file` f WHERE f.del = 0 AND f.id = ?", id)

	// update
	_, err := DbDel(pContext, "UPDATE `file` SET del = 1, `upd_time` = ? WHERE id = ?", time.Now().Unix(), id)
	if err != nil {
		redirect(pid, err)
		return
	}

	redirect(pid, nil)
	return
}

// FileRead 读取文件
func FileRead(pContext *gin.Context, f typ.File) ([]byte, error) {
	// open file
	fPath := FilePath(pContext, f)
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
func FilePath(pContext *gin.Context, f typ.File) string {
	// fType
	fType := f.Type

	// fDir
	dataDir := DataDir(pContext)
	fDir := fmt.Sprintf("%s%s%s", dataDir, util.FileSeparator, fType)
	if !util.IsExistOfPath(fDir) {
		util.Mkdir(fDir)
	}

	// fName -> f.id

	// fPath
	fPath := fmt.Sprintf("%s%s%d", fDir, util.FileSeparator, f.Id)
	return fPath
}
