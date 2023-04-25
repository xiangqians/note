// @author xiangqian
// @date 13:44 2023/04/08
package note

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday/v2"
	"log"
	"note/src/api/common"
	typ_api "note/src/typ/api"
	typ_ft "note/src/typ/ft"
	typ_resp "note/src/typ/resp"
	util_json "note/src/util/json"
	util_str "note/src/util/str"
	util_time "note/src/util/time"
	"os"
	"sort"
	"strings"
)

// View 查看文件页面
func View(context *gin.Context) {
	// id
	id, err := common.Param[int64](context, "id")
	View0(context, id, err)
}

func View0(context *gin.Context, id int64, err any) {
	if err != nil {
		DefaultView(context, typ_api.Note{}, err)
		return
	}

	// query
	note, count, err := DbQry(context, id, 1)
	if err != nil || count == 0 {
		DefaultView(context, note, err)
		return
	}

	// 笔记历史记录
	hist := note.Hist
	if hist != "" {
		// hists
		hists := make([]typ_api.Note, 0, 1) // len 0, cap ?
		err = util_json.Deserialize(hist, &hists)
		if err != nil {
			DefaultView(context, note, err)
			return
		}

		// sort
		sort.Slice(hists, func(i, j int) bool {
			return hists[i].UpdTime > hists[j].UpdTime
		})

		note.Hists = hists
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
func DefaultView(context *gin.Context, note typ_api.Note, err any) {
	resp := typ_resp.Resp[typ_api.Note]{
		Msg:  util_str.TypeToStr(err),
		Data: note,
	}
	common.HtmlOk(context, "note/default/view.html", resp)
}

// Edit 文件修改页
func Edit(context *gin.Context) {
	// id
	id, err := common.Param[int64](context, "id")
	if err != nil {
		FileDefaultEditPage(context, typ_api.Note{}, err)
		return
	}

	// query
	f, count, err := DbQry(context, id, 0)
	if err != nil || count == 0 {
		FileDefaultEditPage(context, f, err)
		return
	}

	// type
	switch typ_ft.ExtNameOf(f.Type) {
	// markdown
	case typ_ft.FtMd:
		FileMdEditPage(context, f)

	// default
	default:
		FileDefaultEditPage(context, f, err)
	}
}

// FileDefaultEditPage 默认文件修改页
func FileDefaultEditPage(context *gin.Context, note typ_api.Note, err error) {
	resp := typ_resp.Resp[typ_api.Note]{
		Msg:  util_str.TypeToStr(err),
		Data: note,
	}
	common.HtmlOk(context, "note/default/edit.html", resp)
}

// FileMdEditPage md文件修改页
func FileMdEditPage(context *gin.Context, note typ_api.Note) {
	html := func(content string, err any) {
		resp := typ_resp.Resp[map[string]any]{
			Msg: util_str.TypeToStr(err),
			Data: map[string]any{
				"note":    note,
				"content": content,
			},
		}

		common.HtmlOk(context, "note/md/edit.html", resp)
	}

	// read
	buf, err := Read(context, note)
	content := ""
	if err == nil {
		content = string(buf)
	}

	html(content, err)
}

// UpdContent 修改文件内容
func UpdContent(context *gin.Context) {
	json := func(err error) {
		if err != nil {
			common.JsonBadRequest(context, typ_resp.Resp[any]{Msg: util_str.TypeToStr(err)})
			return
		}

		common.JsonOk(context, typ_resp.Resp[any]{})
	}

	// id
	id, err := common.PostForm[int64](context, "id")
	if err != nil {
		json(err)
		return
	}
	//log.Println("id", id)

	// f
	f, count, err := DbQry(context, id, 0)
	if count == 0 || typ_ft.ExtNameOf(f.Type) != typ_ft.FtMd {
		json(nil)
		return
	}

	// content
	content, err := common.PostForm[string](context, "content")
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
	_, err = common.DbUpd(context, "UPDATE `note` SET `size` = ?, `upd_time` = ? WHERE id = ?", size, util_time.NowUnix(), id)
	if err != nil {
		json(err)
		return
	}

	json(nil)
	return
}

func Get(context *gin.Context) {
	// id
	id, err := common.Param[int64](context, "id")
	if err != nil {
		log.Println(err)
		return
	}

	// note
	note, count, err := DbQry(context, id, 0)
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

	return
}
