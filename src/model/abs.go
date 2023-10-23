// abs
// @author xiangqian
// @date 20:34 2023/06/10
package model

// Abs 抽象类型定义
type Abs struct {
	Id      int64  `json:"id" form:"id" binding:"gte=0"`     // 主键id
	Rem     string `json:"rem" form:"rem" binding:"max=250"` // 备注
	Del     byte   `json:"del" form:"del"`                   // 删除标识。0-正常，1-删除，2-永久删除
	AddTime int64  `json:"addTime" form:"addTime"`           // 创建时间（时间戳，s）
	UpdTime int64  `json:"updTime" form:"updTime"`           // 修改时间（时间戳，s）
}

type Abs1 struct {
	Abs
	Name        string `json:"name" form:"name" binding:"required,min=1,max=60"` // 名称
	Type        string `json:"type" form:"type"`                                 // 类型
	Size        int64  `json:"size" form:"size"`                                 // 大小，单位：byte
	History     string `json:"history"`                                          // 历史
	HistorySize int64  `json:"historySize"`                                      // 历史大小，单位：byte
}
