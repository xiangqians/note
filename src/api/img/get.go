// get
// @author xiangqian
// @date 14:41 2023/05/20
package img

import (
	"github.com/gin-gonic/gin"
	"log"
	api_common_context "note/src/api/common/context"
)

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
	context.File(path)
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
	context.File(path)
}
