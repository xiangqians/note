// note add
// @author xiangqian
// @date 20:58 2023/04/09
package note

import (
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	typ_api "note/src/typ/api"
	typ_ft "note/src/typ/ft"
	typ_resp "note/src/typ/resp"
	util_str "note/src/util/str"
	util_time "note/src/util/time"
	"os"
	"strings"
)

// Add 新增文件
func Add(context *gin.Context) {
	// note
	note := typ_api.Note{}
	err := common.ShouldBind(context, &note)
	pid := note.Pid

	// redirect
	redirect := func(err any) {
		resp := typ_resp.Resp[any]{Msg: util_str.ConvTypeToStr(err)}
		common.Redirect(context, fmt.Sprintf("/note/list?pid=%d", pid), resp)
	}

	// note err ?
	if err != nil {
		redirect(err)
		return
	}

	// name
	note.Name = strings.TrimSpace(note.Name)
	err = common.VerifyName(note.Name)
	if err != nil {
		redirect(err)
		return
	}

	// 校验文件类型
	// 只支持添加 目录 和 md文件
	ft := typ_ft.ExtNameOf(strings.TrimSpace(note.Type))
	if !(ft == typ_ft.FtD || ft == typ_ft.FtMd) {
		redirect(fmt.Sprintf("%s: %s", i18n.MustGetMessage("i18n.fileTypeUnsupported"), note.Type))
		return
	}
	note.Type = string(ft)

	// add
	id, err := common.DbAdd(context, "INSERT INTO `note` (`pid`, `name`, `type`, `add_time`) VALUES (?, ?, ?, ?)", note.Pid, note.Name, note.Type, util_time.NowUnix())
	if err != nil {
		redirect(err)
		return
	}
	note.Id = id

	// 如果不是目录，则创建物理文件
	if ft != typ_ft.FtD {
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
