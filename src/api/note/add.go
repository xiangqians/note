// note add
// @author xiangqian
// @date 20:58 2023/04/09
package note

import (
	"fmt"
	"github.com/gin-contrib/i18n"
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

// Add 新增文件
func Add(context *gin.Context) {
	// note
	note := typ.Note{}
	err := api_common_context.ShouldBind(context, &note)
	pid := note.Pid

	// redirect
	redirect := func(err any) {
		resp := typ.Resp[any]{Msg: str.ConvTypeToStr(err)}
		api_common_context.Redirect(context, fmt.Sprintf("/note/list?pid=%d", pid), resp)
	}

	// note err ?
	if err != nil {
		redirect(err)
		return
	}

	// name
	note.Name = strings.TrimSpace(note.Name)
	err = validate.FileName(note.Name)
	if err != nil {
		redirect(err)
		return
	}

	// 校验文件类型
	// 只支持添加 目录 和 md文件
	ft := typ.ExtNameOf(strings.TrimSpace(note.Type))
	if !(ft == typ.FtD || ft == typ.FtMd) {
		redirect(fmt.Sprintf("%s: %s", i18n.MustGetMessage("i18n.fileTypeUnsupported"), note.Type))
		return
	}
	note.Type = string(ft)

	// add
	id, err := db.Add(context, "INSERT INTO `note` (`pid`, `name`, `type`, `add_time`) VALUES (?, ?, ?, ?)", note.Pid, note.Name, note.Type, time.NowUnix())
	if err != nil {
		redirect(err)
		return
	}
	note.Id = id

	// 如果不是目录，则创建物理文件
	if ft != typ.FtD {
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
