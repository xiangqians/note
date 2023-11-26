// @author xiangqian
// @date 10:52 2023/02/04
package model

import (
	"sort"
)

// Page 分页数据
type Page struct {
	Current int64   `json:"current"` // 当前页
	Size    uint8   `json:"size"`    // 页数量
	Search  string  `json:"search"`  // 检索条件
	Total   int64   `json:"total"`   // 总数
	Indexes []int64 `json:"indexes"` // 页数索引集
	Data    any     `json:"data"`    // 数据
}

// InitIndexes 初始化页数索引集
func (page *Page) InitIndexes() {
	// 总页数
	pageCount := page.Total / int64(page.Size)
	if page.Total%int64(page.Size) != 0 {
		pageCount += 1
	}

	if pageCount == 0 {
		return
	}

	if page.Current == 1 || page.Current > pageCount {
		indexes := make([]int64, 0, 8)
		var index int64 = 1
		count := cap(indexes)
		for {
			count--
			if count < 0 || index > pageCount {
				break
			}
			indexes = append(indexes, index)
			index++
		}

		length := len(indexes)
		if indexes[length-1] != pageCount {
			indexes[length-2] = 0
			indexes[length-1] = pageCount
		}
		page.Indexes = indexes

	} else if page.Current == pageCount {
		indexes := make([]int64, 0, 8)
		var index int64 = pageCount
		count := cap(indexes)
		for {
			count--
			if count < 0 || index <= 0 {
				break
			}
			indexes = append(indexes, index)
			index--
		}

		// 排序：升序
		sort.Slice(indexes, func(i, j int) bool {
			return i > j
		})

		if indexes[0] != 1 {
			indexes[0] = 1
			indexes[1] = 0
		}
		page.Indexes = indexes

	} else {
		indexes := make([]int64, 0, 6+1+6)
		var index int64 = page.Current - 6
		if index <= 0 {
			index = 1
		}
		i := 0 // 当前页索引在数组中位置
		count := cap(indexes)
		for {
			count--
			if count < 0 || index > pageCount {
				break
			}
			indexes = append(indexes, index)
			if page.Current == index {
				i = len(indexes) - 1
			}
			index++
		}

		length := len(indexes)
		// ... 在右侧
		if indexes[0] == 1 && i < 4 {
			if length >= 8 && indexes[8-1] != pageCount {
				indexes[8-2] = 0
				indexes[8-1] = pageCount
				indexes = indexes[0:8]
			}

		} else
		// ... 在左侧
		if length >= 8 && indexes[length-1] == pageCount && i >= length-4 {
			if indexes[length-8] != 1 {
				indexes[length-8] = 1
				indexes[length-8+1] = 0
				indexes = indexes[length-8:]
			}
		} else
		// ... 在左右两侧
		if length > 8 {
			indexes = indexes[i-4 : i+4+1]
			length = len(indexes)
			indexes[0] = 1
			indexes[1] = 0
			indexes[length-2] = 0
			indexes[length-1] = pageCount
		}
		page.Indexes = indexes
	}
}
