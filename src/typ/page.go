// page type
// @author xiangqian
// @date 10:52 2023/02/04
package typ

// Page 分页数据
type Page[T any] struct {
	Current int64   `json:"current"` // 当前页
	Size    uint8   `json:"size"`    // 页数量
	Total   int64   `json:"total"`   // 总数
	Pages   int64   `json:"pages"`   // 总页数
	Indexes []int64 `json:"indexes"` // 页数索引
	Data    []T     `json:"data"`    // 数据
}
