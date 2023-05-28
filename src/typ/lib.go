// lib type
// @author xiangqian
// @date 11:53 2023/05/07
package typ

import "encoding/gob"

// Lib 库（library）
type Lib struct {
	Abs
	Name     string `json:"name" form:"name" binding:"required,min=1,max=60"` // 名称
	Type     string `json:"type" form:"type"`                                 // 类型
	Size     int64  `json:"size" form:"size"`                                 // 大小，单位：byte
	Hist     string `json:"hist" form:"hist"`                                 // history（历史记录）
	HistSize int64  `json:"histSize" form:"histSize"`                         // history（历史记录）文件大小，单位：byte
	Url      string `json:"url"`                                              // url
	Hists    []Lib  `json:"hists"`                                            // 历史记录
	HistIdx  int8   `json:"histIdx"`                                          // Hists Index
}

// 注册模型
func init() {
	gob.Register(Lib{})
}
