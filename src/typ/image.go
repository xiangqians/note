// lib type
// @author xiangqian
// @date 11:53 2023/05/07
package typ

import "encoding/gob"

// Image 图片
type Image struct {
	Abs
	Name        string `json:"name" form:"name" binding:"required,min=1,max=60"` // 名称
	Type        string `json:"type" form:"type"`                                 // 类型
	Size        int64  `json:"size" form:"size"`                                 // 大小，单位：byte
	History     string `json:"history" form:"history"`                           // 历史
	HistorySize int64  `json:"historySize" form:"historySize"`                   // 历史大小，单位：byte
	//histories   []Image `json:"histories"`                                        // 历史记录
}

// 注册模型
func init() {
	gob.Register(Image{})
}
