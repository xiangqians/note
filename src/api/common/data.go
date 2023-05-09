// data
// @author xiangqian
// @date 19:59 2023/03/22
package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/api/common/session"
	util_os "note/src/util/os"
)

func DataDir(context *gin.Context) string {
	if context == nil {
		return AppArg.DataDir
	}

	user, _ := session.GetSessionUser(context)
	return fmt.Sprintf("%s%s%d", AppArg.DataDir, util_os.FileSeparator(), user.Id)
}
