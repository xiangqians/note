// api type
// @author xiangqian
// @date 14:07 2023/02/04
package typ

const (
	LocaleZh = "zh"
	LocaleEn = "en"
)

// Abs 抽象实体定义
type Abs struct {
	Id      int64  `form:"id" binding:"gte=0"`    // 主键id
	Rem     string `form:"rem" binding:"max=200"` // 备注
	Del     byte   `form:"del"`                   // 删除标识，0-正常，1-删除
	AddTime int64  `form:"addTime"`               // 创建时间（时间戳，s）
	UpdTime int64  `form:"updTime"`               // 修改时间（时间戳，s）
}
