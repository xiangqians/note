// img view
// @author xiangqian
// @date 21:25 2023/04/10
package img

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"note/src/api/common"
	typ_api "note/src/typ/api"
	typ_resp "note/src/typ/resp"
	util_str "note/src/util/str"
	util_time "note/src/util/time"
)

// HistView 查看图片历史页面
func HistView(context *gin.Context) {
	html := func(img typ_api.Img, err any) {
		resp := typ_resp.Resp[typ_api.Img]{
			Msg:  util_str.ConvTypeToStr(err),
			Data: img,
		}
		common.HtmlOk(context, "img/view.html", resp)
	}

	// id
	id, err := common.Param[int64](context, "id")
	if err != nil || id <= 0 {
		html(typ_api.Img{}, err)
		return
	}

	// idx
	idx, err := common.Param[int](context, "idx")
	if err != nil || idx < 0 {
		idx = 0
	}

	// img
	img, count, err := DbQry(context, id, 0)
	if err != nil || count == 0 {
		html(img, err)
		return
	}

	// hists
	hists, err := DeserializeHist(img.Hist)
	if err != nil || hists == nil {
		html(img, err)
		return
	}

	// 校验idx是否合法
	l := len(hists)
	if idx >= l {
		idx = l - 1
	}

	// hist img
	histImg := hists[idx]
	histImg.Url = fmt.Sprintf("/img/%d/hist/%d?t=%d", id, idx, util_time.NowUnix())
	histImg.Hists = hists
	histImg.HistIdx = int8(idx)

	// html
	html(histImg, err)
}

// View 查看图片页面
func View(context *gin.Context) {
	html := func(img typ_api.Img, err any) {
		resp := typ_resp.Resp[typ_api.Img]{
			Msg:  util_str.ConvTypeToStr(err),
			Data: img,
		}
		common.HtmlOk(context, "img/view.html", resp)
	}

	// id
	id, err := common.Param[int64](context, "id")
	if err != nil {
		html(typ_api.Img{}, err)
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
	img.Url = fmt.Sprintf("/img/%d?t=%d", id, util_time.NowUnix())

	// hists
	img.Hists, err = DeserializeHist(img.Hist)

	// html
	html(img, err)
}

// GetHist 获取历史图片
func GetHist(context *gin.Context) {
	// id
	id, err := common.Param[int64](context, "id")
	if err != nil {
		log.Println(err)
		return
	}

	// idx
	idx, err := common.Param[int](context, "idx")
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
	hists, err := DeserializeHist(img.Hist)
	if err != nil || hists == nil {
		log.Println("hist is empty")
		return
	}

	// 校验idx是否合法
	if idx >= len(hists) {
		log.Println(err)
		return
	}

	// hist img
	histImg := hists[idx]

	// path
	path, err := HistPath(context, histImg)
	if err != nil {
		log.Println(err)
		return
	}

	// write
	common.Write(context, path)
}

// Get 获取图片
func Get(context *gin.Context) {
	// id
	id, err := common.Param[int64](context, "id")
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
	common.Write(context, path)
}
