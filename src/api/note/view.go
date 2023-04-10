// @author xiangqian
// @date 13:44 2023/04/08
package note

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday/v2"
	"note/src/api/common"
	typ_api "note/src/typ/api"
	typ_ft "note/src/typ/ft"
	typ_resp "note/src/typ/resp"
	util_str "note/src/util/str"
	"strings"
)

// View 查看文件页面
func View(context *gin.Context) {
	// id
	id, err := common.Param[int64](context, "id")
	if err != nil {
		DefaultView(context, typ_api.Note{}, err)
		return
	}

	// query
	note, count, err := DbQry(context, id, true)
	if err != nil || count == 0 {
		DefaultView(context, note, err)
		return
	}

	// type
	switch typ_ft.ExtNameOf(note.Type) {
	// markdown
	case typ_ft.FtMd:
		MdView(context, note)

	// html
	case typ_ft.FtHtml:
		HtmlView(context, note)

	// pdf
	case typ_ft.FtPdf:
		PdfView(context, note)

	// default
	default:
		DefaultView(context, note, err)
	}
}

// PdfView 查看pdf文件
func PdfView(context *gin.Context, note typ_api.Note) {
	v, _ := common.Query[string](context, "v")
	v = strings.TrimSpace(v)
	switch v {
	case "1.0":
		// v1.0
	case "2.0":
		// v2.0
	default:
		v = "2.0"
	}

	note.Url = fmt.Sprintf("/note/%v", note.Id)

	resp := typ_resp.Resp[typ_api.Note]{
		Data: note,
	}
	common.HtmlOk(context, fmt.Sprintf("note/pdf/view_v%s.html", v), resp)
}

// HtmlView 查看html文件
func HtmlView(context *gin.Context, note typ_api.Note) {
	html := func(html string, err any) {
		resp := typ_resp.Resp[map[string]any]{
			Msg: util_str.TypeToStr(err),
			Data: map[string]any{
				"note": note,
				"html": html,
			},
		}
		common.HtmlOk(context, "note/html/view.html", resp)
	}

	// read
	buf, err := Read(context, note)
	if err != nil {
		html("", err)
		return
	}

	html(string(buf), nil)
}

// MdView 查看md文件
// https://github.com/russross/blackfriday
// https://pkg.go.dev/github.com/russross/blackfriday/v2
func MdView(context *gin.Context, note typ_api.Note) {
	html := func(html string, err any) {
		resp := typ_resp.Resp[map[string]any]{
			Msg: util_str.TypeToStr(err),
			Data: map[string]any{
				"note": note,
				"html": html,
			},
		}
		common.HtmlOk(context, "note/md/view.html", resp)
	}

	// read
	buf, err := Read(context, note)
	if err != nil {
		html("", err)
		return
	}

	//output := blackfriday.Run(input)
	//output := blackfriday.Run(input, blackfriday.WithNoExtensions())
	//output := blackfriday.Run(input, blackfriday.WithExtensions(blackfriday.CommonExtensions))

	// https://github.com/russross/blackfriday/issues/394
	buf = bytes.Replace(buf, []byte("\r"), nil, -1)
	//output := blackfriday.Run(input, blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.HardLineBreak))
	buf = blackfriday.Run(buf, blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.HardLineBreak|blackfriday.AutoHeadingIDs|blackfriday.Autolink))

	// 安全过滤
	//buf = bluemonday.UGCPolicy().SanitizeBytes(buf)

	html(string(buf), nil)
}

// DefaultView 默认查看文件
func DefaultView(context *gin.Context, note typ_api.Note, err error) {
	resp := typ_resp.Resp[typ_api.Note]{
		Msg:  util_str.TypeToStr(err),
		Data: note,
	}
	common.HtmlOk(context, "note/default/view.html", resp)
}
