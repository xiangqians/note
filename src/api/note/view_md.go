// note md view
// @author xiangqian
// @date 15:44 2023/04/30
package note

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday/v2"
	api_common_context "note/src/api/common/context"
	"note/src/typ"
	"note/src/util/str"
)

// ViewMd 查看md文件
// https://github.com/russross/blackfriday
// https://pkg.go.dev/github.com/russross/blackfriday/v2
func ViewMd(context *gin.Context, note typ.Note, hist bool) {
	// read
	var buf []byte
	var err error
	if hist {
		buf, err = ReadHist(context, note)
	} else {
		buf, err = Read(context, note)
	}
	if err == nil && len(buf) > 0 {
		note.Content = string(buf)
	}
	if err == nil && len(buf) > 0 {
		//output := blackfriday.Run(input)
		//output := blackfriday.Run(input, blackfriday.WithNoExtensions())
		//output := blackfriday.Run(input, blackfriday.WithExtensions(blackfriday.CommonExtensions))

		// https://github.com/russross/blackfriday/issues/394
		buf = bytes.Replace(buf, []byte("\r"), nil, -1)
		//output := blackfriday.Run(input, blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.HardLineBreak))
		buf = blackfriday.Run(buf, blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.HardLineBreak|blackfriday.AutoHeadingIDs|blackfriday.Autolink))

		// 安全过滤
		//buf = bluemonday.UGCPolicy().SanitizeBytes(buf)

		note.Content = string(buf)
	}

	// resp
	resp := typ.Resp[typ.Note]{
		Msg:  str.ConvTypeToStr(err),
		Data: note,
	}

	// html
	api_common_context.HtmlOk(context, "note/md/view.html", resp)
}
