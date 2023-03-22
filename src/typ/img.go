// img
// @author xiangqian
// @date 21:12 2023/03/22
package typ

import "encoding/gob"

// Img 图片
type Img struct {
	Abs
	Name string `form:"name" binding:"required,min=1,max=60"` // 图片名称
	Type string `form:"type"`                                 // 图片类型
	Size int64  `form:"size"`                                 // 图片大小，单位：byte

	Url string // 图片url
}

// 注册模型
func init() {
	gob.Register(Img{})
}
