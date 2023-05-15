// note upload
// @author xiangqian
// @date 21:01 2023/04/09
package note

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	api_common_context "note/src/api/common/context"
	"note/src/api/common/db"
	"note/src/typ"
	"note/src/util/os"
	"note/src/util/str"
	"note/src/util/time"
	"note/src/util/validate"
	"strings"
)

// ReUpload 重新上传文件
func ReUpload(context *gin.Context) {
	// 重传文件，需要id
	id, err := api_common_context.PostForm[int64](context, "id")

	// redirect
	redirect := func(err any) {
		resp := typ.Resp[any]{Msg: str.ConvTypeToStr(err)}
		api_common_context.Redirect(context, fmt.Sprintf("/note/%d/view", id), resp)
	}

	// err ?
	if err != nil {
		redirect(err)
		return
	}

	// file header
	fh, err := context.FormFile("file")
	if err != nil || fh == nil {
		redirect(err)
		return
	}

	// name
	name, err := validateName(fh)
	if err != nil {
		redirect(err)
		return
	}

	// validate type
	ft, err := validateType(fh)
	if err != nil {
		redirect(err)
		return
	}
	// type
	_type := string(ft)

	// size
	size := fh.Size

	// 校验 id 是否存在
	note, count, err := DbQry(context, id, 0, 0)
	if err != nil || count == 0 || typ.ExtNameOf(note.Type) == typ.FtD { // 重传上传文件不能是目录类型
		redirect(err)
		return
	}

	// 笔记历史记录
	histNotes, err := DeserializeHist(note.Hist)
	if err != nil {
		redirect(err)
		return
	}
	if histNotes == nil {
		histNotes = make([]typ.Note, 0, 1)
	}

	// 将原笔记添加到历史记录
	histNote := typ.Note{
		Abs: typ.Abs{
			Id:      note.Id,
			AddTime: note.AddTime,
			UpdTime: note.UpdTime,
		},
		Pid:  note.Pid,
		Name: note.Name,
		Type: note.Type,
		Size: note.Size,
	}
	histNotes = append(histNotes, histNote)
	Sort(&histNotes)

	// 备份最近一条历史记录
	// src
	var srcPath string
	srcPath, err = Path(context, histNote)
	if err != nil {
		redirect(err)
		return
	}
	// dst
	var dstPath string
	dstPath, err = HistPath(context, histNote)
	if err != nil {
		redirect(err)
		return
	}
	// copy
	_, err = os.CopyFile(dstPath, srcPath)
	if err != nil {
		redirect(err)
		return
	}

	// 笔记历史记录至多保存15张，超过15张则删除最早地历史笔记
	max := 15
	l := len(histNotes)
	if l > max {
		for i := max; i < l; i++ {
			DelHistNote(context, histNotes[i])
		}
		histNotes = histNotes[:max]
	}

	// hist size
	var histSize int64 = 0
	for _, histNote := range histNotes {
		histSize += histNote.Size
	}

	// serialize
	hist, err := SerializeHist(histNotes)
	if err != nil {
		redirect(err)
		return
	}

	// new note
	newNote := typ.Note{
		Abs: typ.Abs{
			Id:      id,
			UpdTime: time.NowUnix(),
		},
		Name:     name,
		Type:     _type,
		Size:     size,
		Hist:     hist,
		HistSize: histSize,
	}

	// update
	_, err = db.Upd(context, "UPDATE `note` SET `name` = ?, `type` = ?, `size` = ?, `hist` = ?, `hist_size` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ?",
		newNote.Name, newNote.Type, newNote.Size, newNote.Hist, newNote.HistSize, newNote.UpdTime, newNote.Id)
	if err != nil {
		redirect(err)
		return
	}

	// 清空文件
	path, err := ClearNote(context, newNote)
	if err != nil {
		redirect(err)
		return
	}

	// 保存文件
	err = context.SaveUploadedFile(fh, path)

	// 保存文件成功时，判断如果重传不是同一个文件类型，则删除之前文件
	if note.Type != newNote.Type {
		_, err = DelNote(context, note)
	}

	// redirect
	redirect(err)
}

// Upload 上传文件
func Upload(context *gin.Context) {
	// 上传文件，需要pid
	pid, err := api_common_context.PostForm[int64](context, "pid")

	// redirect func
	redirect := func(err any) {
		resp := typ.Resp[any]{
			Msg: str.ConvTypeToStr(err),
		}
		api_common_context.Redirect(context, fmt.Sprintf("/note/list?pid=%d", pid), resp)
	}

	// err ?
	if err != nil || pid < 0 {
		redirect(err)
		return
	}

	// file header
	fh, err := context.FormFile("file")
	if err != nil || fh == nil {
		redirect(err)
		return
	}

	// name
	name, err := validateName(fh)
	if err != nil {
		redirect(err)
		return
	}

	// validate type
	ft, err := validateType(fh)
	if err != nil {
		redirect(err)
		return
	}
	// type
	_type := string(ft)

	// size
	size := fh.Size

	// 校验 pid 是否存在
	if pid != 0 {
		note, count, err := DbQry(context, pid, 0, 0)
		if err != nil || count == 0 || typ.ExtNameOf(note.Type) != typ.FtD { // 父节点必须是目录
			redirect(err)
			return
		}
	}

	// 查询是否有永久删除的笔记记录id，以复用
	id, count, err := DbQryPermlyDelId(context)
	// 新id
	if err != nil || count == 0 {
		id, err = db.Add(context, "INSERT INTO `note` (`pid`, `name`, `type`, `size`, `add_time`) VALUES (?, ?, ?, ?, ?)", pid, name, _type, size, time.NowUnix())
	} else
	// 复用id
	{
		_, err = db.Upd(context, "UPDATE `note` SET `pid` = ?, `name` = ?, `type` = ?, `size` = ?, `hist` = '', `hist_size` = 0, `del` = 0, `add_time` = ?, `upd_time` = 0 WHERE `id` = ?",
			pid, name, _type, size, time.NowUnix(), id)
	}
	if err != nil {
		redirect(err)
		return
	}

	// 清空文件
	path, err := ClearNote(context, typ.Note{Abs: typ.Abs{Id: id}, Type: _type})
	if err != nil {
		redirect(err)
		return
	}

	// 保存文件
	err = context.SaveUploadedFile(fh, path)

	// redirect
	redirect(err)
}

// validateType 校验上传数据类型
// 校验文件类型，只支持上传 html/pdf/zip
func validateType(fh *multipart.FileHeader) (ft typ.Ft, err error) {
	contentType := fh.Header.Get("Content-Type")
	ft = typ.ContentTypeOf(contentType)
	if ft == typ.FtUnk || !(ft == typ.FtHtml || ft == typ.FtPdf || ft == typ.FtZip) {
		err = errors.New(fmt.Sprintf("%s, %s", i18n.MustGetMessage("i18n.fileTypeUnsupported"), contentType))
	}
	return
}

// validateType 校验上传文件名称
func validateName(fh *multipart.FileHeader) (name string, err error) {
	// name
	name = strings.TrimSpace(fh.Filename)
	// validate name
	err = validate.FileName(name)
	return
}

// DbQryPermlyDelId 查询永久删除的笔记记录id，以复用
func DbQryPermlyDelId(context *gin.Context) (int64, int64, error) {
	id, count, err := db.Qry[int64](context, "SELECT `id` FROM `note` WHERE `del` = 2 LIMIT 1")
	return id, count, err
}
