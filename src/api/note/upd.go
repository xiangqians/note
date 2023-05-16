// note update
// @author xiangqian
// @date 22:58 2023/04/25
package note

import (
	"bufio"
	"github.com/gin-gonic/gin"
	api_common_context "note/src/api/common/context"
	"note/src/api/common/db"
	"note/src/typ"
	util_os "note/src/util/os"
	"note/src/util/str"
	"note/src/util/time"
	"note/src/util/validate"
	"os"
	"strings"
)

// UpdContent 修改文件内容
func UpdContent(context *gin.Context) {
	json := func(err error) {
		if err != nil {
			api_common_context.JsonBadRequest(context, typ.Resp[any]{Msg: str.ConvTypeToStr(err)})
			return
		}

		api_common_context.JsonOk(context, typ.Resp[any]{})
	}

	// id
	id, err := api_common_context.PostForm[int64](context, "id")
	if err != nil {
		json(err)
		return
	}

	// qry
	note, count, err := DbQry(context, id, 0, 0)
	if err != nil || count == 0 || typ.ExtNameOf(note.Type) != typ.FtMd {
		json(err)
		return
	}

	// content
	content, err := api_common_context.PostForm[string](context, "content")
	if err != nil {
		json(err)
		return
	}

	// 笔记历史记录
	histNotes, err := DeserializeHist(note.Hist)
	if err != nil {
		json(err)
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
		json(err)
		return
	}
	// dst
	var dstPath string
	dstPath, err = HistPath(context, histNote)
	if err != nil {
		json(err)
		return
	}
	// copy
	_, err = util_os.CopyFile(dstPath, srcPath)
	if err != nil {
		json(err)
		return
	}

	// 笔记历史记录至多保存15条，超过15条则删除最早地历史笔记
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
		json(err)
		return
	}

	// path
	path, err := Path(context, note)
	if err != nil {
		json(err)
		return
	}

	// file
	file, err := os.OpenFile(path,
		os.O_WRONLY|os.O_TRUNC, // 只写（O_WRONLY） & 清空文件（O_TRUNC）
		0666)
	if err != nil {
		json(err)
		return
	}
	defer file.Close()

	// write
	writer := bufio.NewWriter(file)
	writer.WriteString(content)
	writer.Flush()

	// file info
	fileInfo, err := file.Stat()
	if err != nil {
		json(err)
		return
	}

	// size
	size := fileInfo.Size()

	// update
	_, err = db.Upd(context, "UPDATE `note` SET `size` = ?, `hist` = ?, `hist_size` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ?", size, hist, histSize, time.NowUnix(), id)
	if err != nil {
		json(err)
		return
	}

	// json
	json(nil)
}

// UpdName 文件重命名
func UpdName(context *gin.Context) {
	redirect := func(pid int64, err any) {
		RedirectToList(context, pid, err)
	}

	// note
	note := typ.Note{}
	err := api_common_context.ShouldBind(context, &note)
	pid := note.Pid
	if err != nil {
		redirect(pid, err)
		return
	}

	// id
	id := note.Id

	// name
	name := strings.TrimSpace(note.Name)
	// validate name
	err = validate.FileName(name)
	if err != nil {
		redirect(pid, err)
		return
	}

	// update
	_, err = db.Upd(context, "UPDATE `note` SET `name` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `name` <> ?", name, time.NowUnix(), id, name)

	redirect(pid, err)
}
