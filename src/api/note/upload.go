// note upload
// @author xiangqian
// @date 21:01 2023/04/09
package note

import (
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"net/http"
	"note/src/api/common"
	typ_api "note/src/typ/api"
	typ_ft "note/src/typ/ft"
	typ_resp "note/src/typ/resp"
	util_json "note/src/util/json"
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

	// 获取 id / pid
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
		redirect(id, pid, fmt.Sprintf("%s: %s", i18n.MustGetMessage("i18n.fileTypeUnsupported"), contentType))
		return
	}

	// size
	fs := fh.Size

	// 原笔记信息
	var oldNote typ_api.Note

	// 操作数据库
	switch method {
	// 新增笔记
	case http.MethodPost:
		// 校验 pid 是否存在
		if pid != 0 {
			var note typ_api.Note
			var count int64
			note, count, err = DbQry(context, pid, 0)
			if err != nil || count == 0 || typ_ft.ExtNameOf(note.Type) != typ_ft.FtD { // 父节点必须是目录
				redirect(id, pid, err)
				return
			}
		}

		// add
		id, err = common.DbAdd(context, "INSERT INTO `note` (`pid`, `name`, `type`, `size`, `add_time`) VALUES (?, ?, ?, ?, ?)",
			pid, fn, ft, fs, util_time.NowUnix())

	// 修改笔记
	case http.MethodPut:
		// 校验 id 是否存在
		var count int64
		oldNote, count, err = DbQry(context, id, 0)
		if err != nil || count == 0 || typ_ft.ExtNameOf(oldNote.Type) == typ_ft.FtD { // 上传文件不能是目录类型
			redirect(id, pid, err)
			return
		}

		// 笔记历史记录
		hist := oldNote.Hist
		histSize := oldNote.HistSize
		histNotes := make([]typ_api.Note, 0, 1) // len 0, cap ?
		if hist != "" {
			err = util_json.Deserialize(hist, &histNotes)
			if err != nil {
				redirect(id, pid, err)
				return
			}
		}

		// 将原笔记添加到历史记录
		histNote := typ_api.Note{
			Abs: typ_api.Abs{
				Id:      oldNote.Id,
				AddTime: oldNote.AddTime,
				UpdTime: oldNote.UpdTime,
			},
			Pid:  oldNote.Pid,
			Name: oldNote.Name,
			Type: oldNote.Type,
			Size: oldNote.Size,
		}
		histNotes = append(histNotes, histNote)
		// src
		var srcPath string
		srcPath, err = Path(context, histNote)
		if err != nil {
			redirect(id, pid, err)
			return
		}
		// dst
		var dstPath string
		dstPath, err = HistPath(context, histNote)
		if err != nil {
			redirect(id, pid, err)
			return
		}
		// copy
		err = util_os.CopyFile(srcPath, dstPath)
		if err != nil {
			redirect(id, pid, err)
			return
		}

		// 笔记历史记录至多保存15条，超过15条则删除最早地历史笔记
		max := 15
		if len(histNotes) > max {
			l := len(histNotes) - max
			for i := 0; i < l; i++ {
				path, err := HistPath(context, histNotes[i])
				if err == nil {
					util_os.DelFile(path)
				}
			}
			histNotes = histNotes[l:]
		}

		// hist size
		for _, histNote := range histNotes {
			histSize += histNote.Size
		}

		// serialize
		hist, err = util_json.Serialize(histNotes)
		if err != nil {
			redirect(id, pid, err)
			return
		}

		// update
		_, err = common.DbUpd(context, "UPDATE `note` SET `name` = ?, `type` = ?, `size` = ?, `hist` = ?, `hist_size` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ?",
			fn, ft, fs, hist, histSize, util_time.NowUnix(), id)
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
