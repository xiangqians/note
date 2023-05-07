// note md view
// @author xiangqian
// @date 15:44 2023/04/30
package note

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday/v2"
	"note/src/api/common"
	"note/src/typ"
	util_str "note/src/util/str"
)

// ViewMd 查看md文件
// https://github.com/russross/blackfriday
// https://pkg.go.dev/github.com/russross/blackfriday/v2
func ViewMd(context *gin.Context, note typ.Note) {
	// read
	buf, err := Read(context, note)
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
		Msg:  util_str.ConvTypeToStr(err),
		Data: note,
	}

	// html
	common.HtmlOk(context, "note/md/view.html", resp)
}
