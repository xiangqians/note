// note upload
// @author xiangqian
// @date 21:01 2023/04/09
package note

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"net/http"
	"note/src/api/common"
	typ_api "note/src/typ/api"
	typ_ft "note/src/typ/ft"
	typ_resp "note/src/typ/resp"
	util_os "note/src/util/os"
	util_str "note/src/util/str"
	util_time "note/src/util/time"
	"os"
)

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
