// note type
// @author xiangqian
// @date 11:52 2023/05/07
package typ

import "encoding/gob"

// Note 笔记
type Note struct {
	Abs
	Pid      int64  `json:"pid" form:"pid" binding:"gte=0"`                   // 父id
	Name     string `json:"name" form:"name" binding:"required,min=1,max=60"` // 文件名称
	Type     string `json:"type" form:"type"`                                 // 文件类型
	Size     int64  `json:"size" form:"size"`                                 // 文件大小，单位：byte
	Hist     string `json:"hist" form:"hist"`                                 // history（历史记录）
	HistSize int64  `json:"histSize" form:"histSize"`                         // history（历史记录）文件大小，单位：byte
	QryPath  int8   `json:"qryPath"`                                          // 查询路径，0-不查询，1-查询，2-查询并包含自身的
	Path     string `json:"path"`                                             // 笔记路径
	PathLink string `json:"pathLink"`                                         // 笔记路径链接
	Url      string `json:"url"`                                              // 笔记url
	Content  string `json:"content"`                                          // 笔记内容
	Hists    []Note `json:"hists"`                                            // 图片历史记录
	HistIdx  int8   `json:"histIdx"`                                          // Hists Index
	Sub      int8   `json:"sub" form:"sub"`                                   // 是否包含所有子集，0-否，1-是
	Deleted  int8   `json:"deleted" form:"deleted"`                           // 是否包含已删除文件，0-否，1-是
	Children []Note `json:"children"`                                         // 子集
}

// 注册模型
func init() {
	gob.Register(Note{})
}
