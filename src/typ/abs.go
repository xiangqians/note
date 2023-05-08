// abs type
// @author xiangqian
// @date 11:51 2023/05/07
package typ

// Abs 抽象实体定义
type Abs struct {
	Id      int64  `json:"id" form:"id" binding:"gte=0"`     // 主键id
	Rem     string `json:"rem" form:"rem" binding:"max=250"` // 备注
	Del     byte   `json:"del" form:"del"`                   // 删除标识，0-正常，1-删除，2-永久删除
	AddTime int64  `json:"addTime" form:"addTime"`           // 创建时间（时间戳，s）
	UpdTime int64  `json:"updTime" form:"updTime"`           // 修改时间（时间戳，s）
}