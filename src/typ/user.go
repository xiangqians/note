// user type
// @author xiangqian
// @date 11:51 2023/05/07
package typ

import "encoding/gob"

// User 用户
type User struct {
	Abs
	Name     string `json:"name" form:"name" binding:"required,excludes= ,min=1,max=60"`                  // 用户名
	Nickname string `json:"nickname" form:"nickname"binding:"max=60"`                                     // 昵称
	Passwd   string `json:"passwd" form:"passwd" binding:"required,excludes= ,max=60"`                    // 密码
	RePasswd string `json:"rePasswd" form:"rePasswd" binding:"required,excludes= ,max=60,eqfield=Passwd"` // retype Passwd
}

// 注册模型
func init() {
	gob.Register(User{})
}
