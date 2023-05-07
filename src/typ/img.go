// img type
// @author xiangqian
// @date 11:53 2023/05/07
package typ

import "encoding/gob"

// Img 图片
type Img struct {
	Abs
	Name     string `json:"name" form:"name" binding:"required,min=1,max=60"` // 图片名称
	Type     string `json:"type" form:"type"`                                 // 图片类型
	Size     int64  `json:"size" form:"size"`                                 // 图片大小，单位：byte
	Hist     string `json:"hist" form:"hist"`                                 // history（历史记录）
	HistSize int64  `json:"histSize" form:"histSize"`                         // history（历史记录）文件大小，单位：byte
	Url      string `json:"url"`                                              // 图片url
	Hists    []Img  `json:"hists"`                                            // 图片历史记录
	HistIdx  int8   `json:"histIdx"`                                          // Hists Index
}

// 注册模型
func init() {
	gob.Register(Img{})
}
