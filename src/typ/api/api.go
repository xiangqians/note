// api type
// @author xiangqian
// @date 14:07 2023/02/04
package api

import (
	"encoding/gob"
)

// Abs 抽象实体定义
type Abs struct {
	Id      int64  `form:"id" binding:"gte=0"`    // 主键id
	Rem     string `form:"rem" binding:"max=200"` // 备注
	Del     byte   `form:"del"`                   // 删除标识，0-正常，1-删除
	AddTime int64  `form:"addTime"`               // 创建时间（时间戳，s）
	UpdTime int64  `form:"updTime"`               // 修改时间（时间戳，s）
}

// User 用户
type User struct {
	Abs
	Name     string `form:"name" binding:"required,excludes= ,min=1,max=60"`               // 用户名
	Nickname string `form:"nickname"binding:"max=60"`                                      // 昵称
	Passwd   string `form:"passwd" binding:"required,excludes= ,max=100"`                  // 密码
	RePasswd string `form:"rePasswd" binding:"required,excludes= ,max=100,eqfield=Passwd"` // retype Passwd
}

// Note 笔记
type Note struct {
	Abs
	Pid      int64  `form:"pid" binding:"gte=0"`                  // 父id
	Name     string `form:"name" binding:"required,min=1,max=60"` // 文件名称
	Type     string `form:"type"`                                 // 文件类型
	Size     int64  `form:"size"`                                 // 文件大小，单位：byte
	Hist     string `form:"hist"`                                 // history（历史记录）
	HistSize int64  `form:"histSize"`                             // history（历史记录）文件大小，单位：byte
	Path     string // 目录路径
	Url      string // 笔记url
}

// Img 图片
type Img struct {
	Abs
	Name     string `form:"name" binding:"required,min=1,max=60"` // 图片名称
	Type     string `form:"type"`                                 // 图片类型
	Size     int64  `form:"size"`                                 // 图片大小，单位：byte
	Hist     string `form:"hist"`                                 // history（历史记录）
	HistSize int64  `form:"histSize"`                             // history（历史记录）文件大小，单位：byte
	Url      string // 图片url
	Hists    []Img  // 图片历史记录
}

// Stat 统计（file和img）
type Stat struct {
	Type string // 文件类型
	Num  int64  // 文件数量
	Size int64  // 文件大小
}

// 注册模型
func init() {
	gob.Register(User{})
	gob.Register(Note{})
	gob.Register(Img{})
	gob.Register(Stat{})
}
