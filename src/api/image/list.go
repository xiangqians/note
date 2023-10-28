// 图片列表
// @author xiangqian
// @date 20:29 2023/04/27
package image

import (
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	"note/src/model"
)

func List(ctx *gin.Context) {
	common.List[model.Image](ctx)
	return
}
