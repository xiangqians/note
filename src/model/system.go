// @author xiangqian
// @date 20:35 2023/06/10
package model

import "encoding/gob"

// System 系统信息表
type System struct {
	Passwd            string // 密码
	Try               byte   // 尝试输入次数
	LastSignInIp      string // 上一次登录IP
	LastSignInTime    int64  // 上一次登录时间（时间戳，单位s）
	CurrentSignInIp   string // 当前次登录IP
	CurrentSignInTime int64  // 当前登录时间（时间戳，单位s）
	UpdTime           int64  // 修改时间（时间戳，单位s）
}

// 注册模型
func init() {
	gob.Register(System{})
}
