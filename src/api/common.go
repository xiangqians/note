// Route
// @author xiangqian
// @date 21:47 2022/12/23
package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"note/src/db"
	"note/src/session"
	"note/src/typ"
	util_os "note/src/util/os"
)

// Db 获取数据库操作实例
func Db(ctx *gin.Context) (*gorm.DB, error) {
	dataDir := typ.GetArg().DataDir
	if ctx == nil {
		return db.Db(util_os.Path(dataDir, "database.db"))
	}

	user, err := session.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	return db.Db(util_os.Path(dataDir, fmt.Sprintf("%d", user.Id), "database.db"))
}
