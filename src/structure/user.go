// @author xiangqian
// @date 20:35 2023/06/10
package structure

import "encoding/gob"

// User 用户
type User struct {
	Abs
	Name              string `json:"name" form:"name" binding:"required,excludes= ,min=1,max=60"` // 用户名
	Nickname          string `json:"nickname" form:"nickname" binding:"max=60"`                   // 昵称
	Passwd            string `json:"passwd" form:"passwd" binding:"required,excludes= ,max=120"`  // 密码
	Try               byte   `json:"try" form:"try"`                                              // 尝试输入次数
	LastSignInIp      string `json:"last_sign_in_ip" form:"last_sign_in_ip"`                      // 上一次登录IP
	LastSignInTime    int64  `json:"last_sign_in_time" form:"last_sign_in_time"`                  // 上一次登录时间（时间戳，单位s）
	CurrentSignInIp   string `json:"current_sign_in_ip" form:"current_sign_in_ip"`                // 当前次登录IP
	CurrentSignInTime int64  `json:"current_sign_in_time" form:"current_sign_in_time"`            // 当前登录时间（时间戳，单位s）
}

// AddUser 新增用户
type AddUser struct {
	User
	RePasswd string `json:"rePasswd" form:"rePasswd" binding:"required,excludes= ,max=120,eqfield=Passwd"` // retype Passwd，再次输入密码
}

// UpdUser 修改用户
type UpdUser struct {
	User
	OrigPasswd  string `json:"origPasswd" form:"origPasswd" binding:"required,excludes= ,max=120"`                     // Original Password，原密码
	NewPasswd   string `json:"newPasswd" form:"newPasswd" binding:"required,excludes= ,max=120"`                       // 新密码
	ReNewPasswd string `json:"reNewPasswd" form:"reNewPasswd" binding:"required,excludes= ,max=120,eqfield=NewPasswd"` // retype new Passwd，再次输入新密码
}

// 注册模型
func init() {
	gob.Register(User{})
	gob.Register(AddUser{})
	gob.Register(UpdUser{})
}
