// img view
// @author xiangqian
// @date 21:25 2023/04/10
package img

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
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

// view 查看图片页面
// hist: 是否是历史记录
func view(context *gin.Context, hist bool) {
	// html
	html := func(img typ.Img, err any) {
		resp := typ.Resp[typ.Img]{
			Msg:  str.ConvTypeToStr(err),
			Data: img,
		}
		api_common_context.HtmlOk(context, "img/view.html", resp)
	}

	// id
	id, err := api_common_context.Param[int64](context, "id")
	if err != nil {
		html(typ.Img{}, err)
		return
	}

	// img
	img, count, err := DbQry(context, id, 0)
	// current img
	img.HistIdx = -1
	// err ? / count == 0 ?
	if err != nil || count == 0 {
		html(img, err)
		return
	}

	// url
	img.Url = fmt.Sprintf("/img/%d?t=%d", id, time.NowUnix())

	// hists
	img.Hists, err = DeserializeHist(img.Hist)

	// 如果查询的是历史记录
	if hist {
		hists := img.Hists
		if err != nil || hists == nil {
			html(img, err)
			return
		}

		// 校验idx是否合法
		l := len(hists)
		var idx int
		idx, err = api_common_context.Param[int](context, "idx")
		if err != nil || idx < 0 {
			idx = 0
		}
		if idx >= l {
			idx = l - 1
		}

		// hist img
		histImg := hists[idx]
		histImg.Url = fmt.Sprintf("/img/%d/hist/%d?t=%d", id, idx, time.NowUnix())
		histImg.Hists = hists
		histImg.HistIdx = int8(idx)
		img = histImg
	}

	// html
	html(img, err)
}

// GetHist 获取历史图片
func GetHist(context *gin.Context) {
	// id
	id, err := api_common_context.Param[int64](context, "id")
	if err != nil {
		log.Println(err)
		return
	}

	// idx
	idx, err := api_common_context.Param[int](context, "idx")
	if err != nil || idx < 0 {
		log.Println(err)
		return
	}

	// img
	img, count, err := DbQry(context, id, 0)
	if err != nil || count == 0 {
		log.Println(err)
		return
	}

	// hists
	histImgs, err := DeserializeHist(img.Hist)
	if err != nil || histImgs == nil {
		log.Println("hist is empty")
		return
	}

	// 校验idx是否合法
	if idx >= len(histImgs) {
		log.Println(err)
		return
	}

	// hist img
	histImg := histImgs[idx]

	// path
	path, err := HistPath(context, histImg)
	if err != nil {
		log.Println(err)
		return
	}

	// write
	api_common_context.WriteFile(context, path)
}

// Get 获取图片
func Get(context *gin.Context) {
	// id
	id, err := api_common_context.Param[int64](context, "id")
	if err != nil {
		log.Println(err)
		return
	}

	// img
	img, count, err := DbQry(context, id, 0)
	if err != nil || count == 0 {
		log.Println(err)
		return
	}

	// path
	path, err := Path(context, img)
	if err != nil {
		log.Println(err)
		return
	}

	// write
	api_common_context.WriteFile(context, path)
}
