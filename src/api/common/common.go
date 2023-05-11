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

func DataDir(context *gin.Context) string {
	if context == nil {
		return AppArg.DataDir
	}

	user, _ := session.GetUser(context)
	return fmt.Sprintf("%s%s%d", AppArg.DataDir, os.FileSeparator(), user.Id)
}
