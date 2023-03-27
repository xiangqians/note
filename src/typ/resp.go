// response
// @author xiangqian
// @date 20:24 2023/03/22
package typ

import "encoding/gob"

// Resp 响应数据
type Resp[T any] struct {
	Msg  string // 消息（没有消息就是最好的消息）
	Data T      // 数据
}

// 注册模型
func init() {
	gob.Register(Resp[any]{})
}
