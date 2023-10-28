// @author xiangqian
// @date 23:23 2023/10/22
package image

import (
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	"note/src/model"
)

// Del 图片删除
func Del(ctx *gin.Context) {
	common.Del[model.Image](ctx)
}
