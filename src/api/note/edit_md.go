// edit md
// @author xiangqian
// @date 15:38 2023/05/13
package note

import (
	"bufio"
	"github.com/gin-gonic/gin"
	api_common_context "note/src/api/common/context"
	"note/src/api/common/db"
	"note/src/typ"
	"note/src/util/str"
	"note/src/util/time"
	"os"
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

	// note
	note, count, err := DbQry(context, id, 0, 0)
	if err != nil || count == 0 || typ.ExtNameOf(note.Type) != typ.FtMd {
		json(nil)
		return
	}

	// content
	content, err := api_common_context.PostForm[string](context, "content")
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
	_, err = db.Upd(context, "UPDATE `note` SET `size` = ?, `upd_time` = ? WHERE id = ?", size, time.NowUnix(), id)
	if err != nil {
		json(err)
		return
	}

	// json
	json(nil)
}

// EditMd md文件修改页
func EditMd(context *gin.Context, note typ.Note) {
	html := func(content string, err any) {
		resp := typ.Resp[map[string]any]{
			Msg: str.ConvTypeToStr(err),
			Data: map[string]any{
				"note":    note,
				"content": content,
			},
		}

		api_common_context.HtmlOk(context, "note/md/edit.html", resp)
	}

	// read
	buf, err := Read(context, note)
	content := ""
	if err == nil {
		content = string(buf)
	}

	html(content, err)
}
