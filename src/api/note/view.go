// @author xiangqian
// @date 13:44 2023/04/08
package note

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	typ2 "note/app/typ"
)

// View 查看文件页面
func View(context *gin.Context) {
	// id
	id, err := context.Param[int64](context, "id")
	if err != nil {
		ViewUnsupported(context, typ2.Note{}, err)
		return
	}

	// query
	note, count, err := DbQryNew(context, id, 1, 0)
	if err != nil || count == 0 {
		ViewUnsupported(context, note, err)
		return
	}

	// url
	note.Url = fmt.Sprintf("/note/%d?t=%d", id, time.NowUnix())

	// 笔记历史记录
	note.Hists, err = DeserializeHist(note.Hist)
	if err != nil {
		ViewUnsupported(context, note, err)
		return
	}

	// type
	switch typ2.ExtNameOf(note.Type) {
	// markdown
	case typ2.FtMd:
		ViewMd(context, note)

	// html
	case typ2.FtHtml:
		ViewHtml(context, note)

	// pdf
	case typ2.FtPdf:
		ViewPdf(context, note)

	// default
	default:
		ViewUnsupported(context, note, err)
	}
}

// ViewUnsupported 不支持查看
func ViewUnsupported(context *gin.Context, note typ2.Note, err any) {
	resp := typ2.Resp[typ2.Note]{
		Msg:  str.ConvTypeToStr(err),
		Data: note,
	}
	context.HtmlOk(context, "note/unsupported/view.html", resp)
}

func Get(context *gin.Context) {
	// id
	id, err := context.Param[int64](context, "id")
	if err != nil {
		log.Println(err)
		return
	}

	// note
	note, count, err := DbQry(context, typ2.Note{Abs: typ2.Abs{Id: id}, Pid: -1})
	if err != nil || count == 0 {
		log.Println(err)
		return
	}

	// 排除目录
	if typ2.FtD == typ2.ExtNameOf(note.Type) {
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
