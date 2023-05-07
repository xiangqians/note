// index
// @author xiangqian
// @date 17:21 2023/02/04
package index

import (
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	"note/src/typ"
	util_str "note/src/util/str"
)

// Index index页面
func Index(context *gin.Context) {
	html := func(noteStats, imgStats []typ.Stat, err error) {
		resp := typ.Resp[map[string][]typ.Stat]{
			Msg: util_str.ConvTypeToStr(err),
			Data: map[string][]typ.Stat{
				"noteStats": noteStats,
				"imgStats":  imgStats,
			},
		}
		common.HtmlOk(context, "index.html", resp)
	}

	// stat
	noteStats := []typ.Stat{}
	imgStats := []typ.Stat{}

	// note
	stats, count, err := common.DbQry[[]typ.Stat](context, "SELECT `type`, COUNT(`id`) AS 'num', SUM(`size`) AS 'size', SUM(`hist_size`) AS 'hist_size' FROM `note` WHERE `del` = 0 GROUP BY `type` ORDER BY COUNT(`id`) DESC")
	if err != nil {
		html(noteStats, imgStats, err)
		return
	}
	if count > 0 {
		noteStats = stats
	}

	// img
	stats, count, err = common.DbQry[[]typ.Stat](context, "SELECT `type`, COUNT(`id`) AS 'num', SUM(`size`) AS 'size', SUM(`hist_size`) AS 'hist_size' FROM `img` WHERE `del` = 0 GROUP BY `type` ORDER BY COUNT(`id`) DESC")
	if count > 0 {
		imgStats = stats
	}

	// html
	html(noteStats, imgStats, err)
}
