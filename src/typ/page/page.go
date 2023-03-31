// page type
// @author xiangqian
// @date 10:52 2023/02/04
package page

import "encoding/gob"

// PageReq 分页请求
type PageReq struct {
	Current int64 `json:"current" form:"current"  binding:"gt=0"` // 当前页
	Size    uint8 `json:"size" form:"size" binding:"gt=0"`        // 页数量
}

// Page 分页数据
type Page[T any] struct {
	Current int64 `json:"current"` // 当前页
	Size    uint8 `json:"size"`    // 页数量
	Pages   int64 `json:"pages"`   // 总页数
	Total   int64 `json:"total"`   // 总数
	Data    []T   `json:"data"`    // 数据
}

// 注册模型
func init() {
	gob.Register(PageReq{})
}
