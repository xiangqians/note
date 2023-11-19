// @author xiangqian
// @date 11:53 2023/05/07
package model

import "encoding/gob"

// Image 图片
type Image struct {
	Abs1
	//histories   []Image `json:"histories"` // 历史记录
}

// 注册模型
func init() {
	gob.Register(Image{})
}
