// index
// @author xiangqian
// @date 17:21 2023/02/04
package index

import (
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	typ_api "note/src/typ/api"
	typ_resp "note/src/typ/resp"
	util_str "note/src/util/str"
)

// Index index页面
func Index(context *gin.Context) {
	html := func(noteStats, imgStats []typ_api.Stat, err error) {
		resp := typ_resp.Resp[map[string][]typ_api.Stat]{
			Msg: util_str.ConvTypeToStr(err),
			Data: map[string][]typ_api.Stat{
				"noteStats": noteStats,
				"imgStats":  imgStats,
			},
		}
		common.HtmlOk(context, "index.html", resp)
	}

	// stat
	noteStats := []typ_api.Stat{}
	imgStats := []typ_api.Stat{}

	// note
	stats, count, err := common.DbQry[[]typ_api.Stat](context, "SELECT `type`, COUNT(`id`) AS 'num', SUM(`size`) AS 'size', SUM(`hist_size`) AS 'hist_size' FROM `note` WHERE `del` = 0 GROUP BY `type` ORDER BY COUNT(`id`) DESC")
	if err != nil {
		html(noteStats, imgStats, err)
		return
	}
	if count > 0 {
		noteStats = stats
	}

	// img
	stats, count, err = common.DbQry[[]typ_api.Stat](context, "SELECT `type`, COUNT(`id`) AS 'num', SUM(`size`) AS 'size', SUM(`hist_size`) AS 'hist_size' FROM `img` WHERE `del` = 0 GROUP BY `type` ORDER BY COUNT(`id`) DESC")
	if count > 0 {
		imgStats = stats
	}

	// html
	html(noteStats, imgStats, err)
}
