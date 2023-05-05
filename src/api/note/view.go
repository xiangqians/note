// @author xiangqian
// @date 13:44 2023/04/08
package note

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"note/src/api/common"
	typ_api "note/src/typ/api"
	typ_ft "note/src/typ/ft"
	typ_resp "note/src/typ/resp"
	util_str "note/src/util/str"
	util_time "note/src/util/time"
)

// View 查看文件页面
func View(context *gin.Context) {
	// id
	id, err := common.Param[int64](context, "id")
	if err != nil {
		ViewUnsupported(context, typ_api.Note{}, err)
		return
	}

	// query
	note, count, err := DbQryNew(context, id, 1, typ_api.NotDeleted)
	if err != nil || count == 0 {
		ViewUnsupported(context, note, err)
		return
	}

	// url
	note.Url = fmt.Sprintf("/note/%d?t=%d", id, util_time.NowUnix())

	// 笔记历史记录
	note.Hists, err = DeserializeHist(note.Hist)
	if err != nil {
		ViewUnsupported(context, note, err)
		return
	}

	// type
	switch typ_ft.ExtNameOf(note.Type) {
	// markdown
	case typ_ft.FtMd:
		ViewMd(context, note)

	// html
	case typ_ft.FtHtml:
		ViewHtml(context, note)

	// pdf
	case typ_ft.FtPdf:
		ViewPdf(context, note)

	// default
	default:
		ViewUnsupported(context, note, err)
	}
}

// ViewUnsupported 不支持查看
func ViewUnsupported(context *gin.Context, note typ_api.Note, err any) {
	resp := typ_resp.Resp[typ_api.Note]{
		Msg:  util_str.ConvTypeToStr(err),
		Data: note,
	}
	common.HtmlOk(context, "note/unsupported/view.html", resp)
}

func Get(context *gin.Context) {
	// id
	id, err := common.Param[int64](context, "id")
	if err != nil {
		log.Println(err)
		return
	}

	// note
	note, count, err := DbQry(context, typ_api.Note{Abs: typ_api.Abs{Id: id}, Pid: -1})
	if err != nil || count == 0 {
		log.Println(err)
		return
	}

	// 排除目录
	if typ_ft.FtD == typ_ft.ExtNameOf(note.Type) {
		return
	}

	// path
	path, err := Path(context, note)
	if err != nil {
		log.Println(err)
		return
	}

	/**
	// read all
	buf, err := os.ReadFile(fPath)
	if err != nil {
		log.Println(err)
		return
	}
	writer := context.Writer
	writer.Write(buf)
	writer.Flush()
	*/

	/**
	// open
	pFile, err := os.Open(fPath)
	if err != nil {
		log.Println(err)
		return
	}

	// write
	err = util.IOCopy(pFile, context.Writer, 0)
	if err != nil {
		log.Println(err)
		return
	}
	*/

	context.File(path)
}
