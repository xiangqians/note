// data
// @author xiangqian
// @date 19:59 2023/03/22
package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/api/common/session"
	"note/src/typ"
	"note/src/util/os"
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
