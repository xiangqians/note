// index
// @author xiangqian
// @date 17:21 2023/02/04
package index

import (
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	"note/src/typ"
	"note/src/util"
)

// Page index页面
func Page(context *gin.Context) {
	html := func(fileStats, imgStats []typ.Stat, err error) {
		resp := typ.Resp[map[string][]typ.Stat]{
			Msg: util.TypeAsStr(err),
			Data: map[string][]typ.Stat{
				"FileStats": fileStats,
				"ImgStats":  imgStats,
			},
		}
		common.HtmlOkNew(context, "index.html", resp)
	}

	// stat
	fileStats := []typ.Stat{}
	imgStats := []typ.Stat{}

	// file
	stats, count, err := common.DbQry[[]typ.Stat](context, "SELECT `type`, COUNT(id) AS 'num', SUM(`size`) AS 'size' FROM file WHERE del = 0 GROUP BY `type` ORDER BY COUNT(id) DESC")
	if err != nil {
		html(fileStats, imgStats, err)
		return
	}
	if count > 0 {
		fileStats = stats
	}

	// img
	stats, count, err = common.DbQry[[]typ.Stat](context, "SELECT `type`, COUNT(id) AS 'num', SUM(`size`) AS 'size' FROM img WHERE del = 0 GROUP BY `type` ORDER BY COUNT(id) DESC")
	if count > 0 {
		imgStats = stats
	}

	// html
	html(fileStats, imgStats, err)
}
