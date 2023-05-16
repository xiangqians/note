// data
// @author xiangqian
// @date 19:59 2023/03/22
package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	api_common_context "note/src/api/common/context"
	"note/src/api/common/session"
	"note/src/typ"
	"note/src/util/os"
	"note/src/util/str"
)

var AppArg typ.AppArg

// DataDir 根据 context 获取数据目录
func DataDir(context *gin.Context) string {
	if context == nil {
		return AppArg.DataDir
	}

	user, _ := session.GetUser(context)
	return DataDirOnUserId(user.Id)
}

// DataDirOnUserId 根据用户id获取数据目录
func DataDirOnUserId(userId int64) string {
	return fmt.Sprintf("%s%s%d", AppArg.DataDir, os.FileSeparator(), userId)
}

// PageReq 获取分页查询参数
func PageReq(context *gin.Context) (current int64, size uint8) {
	//  current
	current, _ = api_common_context.Param[int64](context, "current")
	if current <= 0 {
		current = 1
	}

	// size
	size, _ = api_common_context.Param[uint8](context, "size")
	if size <= 0 {
		size = 10
	}

	return
}

// DataNotExist 数据不存在
func DataNotExist(context *gin.Context, err error) {
	api_common_context.HtmlOk(context, "dataNotExist.html", typ.Resp[any]{
		Msg: str.ConvTypeToStr(err),
	})
}
