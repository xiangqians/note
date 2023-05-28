// stat type
// @author xiangqian
// @date 11:54 2023/05/07
package typ

import "encoding/gob"

// Stat 统计note/lib
type Stat struct {
	Type     string `json:"type"`     // 文件类型
	Num      int64  `json:"num"`      // 文件数量
	Size     int64  `json:"size"`     // 文件大小
	HistSize int64  `json:"histSize"` // 文件历史大小
}

// 注册模型
func init() {
	gob.Register(Stat{})
}
