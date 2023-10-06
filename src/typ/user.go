// user
// @author xiangqian
// @date 20:35 2023/06/10
package typ

import "encoding/gob"

// User 用户
type User struct {
	Abs
	Name       string `json:"name" form:"name" binding:"required,excludes= ,min=1,max=32"`                   // 用户名
	Nickname   string `json:"nickname" form:"nickname" binding:"max=60"`                                     // 昵称
	OrigPasswd string `json:"origPasswd" form:"origPasswd"`                                                  // Original Password，原密码
	Passwd     string `json:"passwd" form:"passwd" binding:"required,excludes= ,max=120"`                    // 密码
	RePasswd   string `json:"rePasswd" form:"rePasswd" binding:"required,excludes= ,max=120,eqfield=Passwd"` // retype Passwd，再次输入密码
	Try        byte   `json:"try" form:"try"`                                                                // 尝试输入次数
}

// 注册模型
func init() {
	gob.Register(User{})
}
