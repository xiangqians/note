// page type
// @author xiangqian
// @date 10:52 2023/02/04
package model

// Page 分页数据
type Page[T any] struct {
	Current     int64   `json:"current"`     // 当前页
	Size        uint8   `json:"size"`        // 页数量
	Total       int64   `json:"total"`       // 总数
	PageCount   int64   `json:"pageCount"`   // 总页数
	PageIndexes []int64 `json:"pageIndexes"` // 页数索引集
	Data        []T     `json:"data"`        // 数据
	Search      string  `json:"search"`      // 检索条件
}
