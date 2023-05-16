// @author xiangqian
// @date 13:44 2023/04/08
package note

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	api_common_context "note/src/api/common/context"
	"note/src/typ"
	"note/src/util/str"
	"note/src/util/time"
)

// HistView 查看历史页面
func HistView(context *gin.Context) {
	view(context, true)
}

// View 查看页面
func View(context *gin.Context) {
	view(context, false)
}

// view 查看页面
// hist: 是否是历史记录
func view(context *gin.Context, hist bool) {
	// id
	id, err := api_common_context.Param[int64](context, "id")
	if err != nil {
		common.DataNotExist(context, err)
		return
	}

	// query
	note, count, err := DbQry(context, id, 1, 0)
	if err != nil || count == 0 {
		common.DataNotExist(context, err)
		return
	}

	// url
	note.Url = fmt.Sprintf("/note/%d?t=%d", id, time.NowUnix())

	// 笔记历史记录
	note.Hists, err = DeserializeHist(note.Hist)
	note.HistIdx = -1
	if err != nil {
		common.DataNotExist(context, err)
		return
	}

	// 如果查询的是历史记录
	if hist {
		histNotes := note.Hists
		if histNotes == nil {
			common.DataNotExist(context, err)
			return
		}

		// 校验idx是否合法
		l := len(histNotes)
		var idx int
		idx, err = api_common_context.Param[int](context, "idx")
		if err != nil || idx < 0 || idx >= l {
			common.DataNotExist(context, err)
			return
		}

		// hist note
		histNote := histNotes[idx]
		histNote.Path = note.Path
		histNote.PathLink = note.PathLink
		histNote.Url = fmt.Sprintf("/note/%d/hist/%d?t=%d", id, idx, time.NowUnix())
		histNote.Hists = histNotes
		histNote.HistIdx = int8(idx)
		note = histNote
	}

	// type
	switch typ.ExtNameOf(note.Type) {
	// markdown
	case typ.FtMd:
		ViewMd(context, note, hist)

	// html
	case typ.FtHtml:
		ViewHtml(context, note, hist)

	// pdf
	case typ.FtPdf:
		ViewPdf(context, note)

	// unsupported，不支持查看
	default:
		resp := typ.Resp[typ.Note]{
			Msg:  str.ConvTypeToStr(err),
			Data: note,
		}
		api_common_context.HtmlOk(context, "note/unsupported/view.html", resp)
	}
}
