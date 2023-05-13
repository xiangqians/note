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
	//log.Println("id", id)

	// f
	f, count, err := DbQry(context, typ.Note{Abs: typ.Abs{Id: id}, Pid: -1})
	if count == 0 || typ.ExtNameOf(f.Type) != typ.FtMd {
		json(nil)
		return
	}

	// content
	content, err := api_common_context.PostForm[string](context, "content")
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
	_, err = db.Upd(context, "UPDATE `note` SET `size` = ?, `upd_time` = ? WHERE id = ?", size, time.NowUnix(), id)
	if err != nil {
		json(err)
		return
	}

	json(nil)
	return
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

	// name
	note.Name = strings.TrimSpace(note.Name)
	err = validate.FileName(note.Name)
	if err != nil {
		redirect(pid, err)
		return
	}

	// update
	_, err = db.Upd(context, "UPDATE `note` SET `name` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `name` <> ?", note.Name, time.NowUnix(), note.Id, note.Name)

	redirect(pid, err)
	return
}
