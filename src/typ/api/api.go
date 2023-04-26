// api type
// @author xiangqian
// @date 14:07 2023/02/04
package api

import (
	"encoding/gob"
)

// Abs 抽象实体定义
type Abs struct {
	Id      int64  `json:"id" form:"id" binding:"gte=0"`     // 主键id
	Rem     string `json:"rem" form:"rem" binding:"max=250"` // 备注
	Del     byte   `json:"del" form:"del"`                   // 删除标识，0-正常，1-删除，2-永久删除
	AddTime int64  `json:"addTime" form:"addTime"`           // 创建时间（时间戳，s）
	UpdTime int64  `json:"updTime" form:"updTime"`           // 修改时间（时间戳，s）
}

// User 用户
type User struct {
	Abs
	Name     string `json:"name" form:"name" binding:"required,excludes= ,min=1,max=60"`                  // 用户名
	Nickname string `json:"nickname" form:"nickname"binding:"max=60"`                                     // 昵称
	Passwd   string `json:"passwd" form:"passwd" binding:"required,excludes= ,max=60"`                    // 密码
	RePasswd string `json:"rePasswd" form:"rePasswd" binding:"required,excludes= ,max=60,eqfield=Passwd"` // retype Passwd
}

// Note 笔记
type Note struct {
	Abs
	Pid      int64  `json:"pid" form:"pid" binding:"gte=0"`                   // 父id
	Name     string `json:"name" form:"name" binding:"required,min=1,max=60"` // 文件名称
	Type     string `json:"type" form:"type"`                                 // 文件类型
	Size     int64  `json:"size" form:"size"`                                 // 文件大小，单位：byte
	Hist     string `json:"hist" form:"hist"`                                 // history（历史记录）
	HistSize int64  `json:"histSize" form:"histSize"`                         // history（历史记录）文件大小，单位：byte
	QryPath  int8   `json:"qryPath"`                                          // 查询路径，0-不查询，1-查询，2-查询并包含自身的
	Path     string `json:"path"`                                             // 笔记路径
	PathLink string `json:"pathLink"`                                         // 笔记路径链接
	Url      string `json:"url"`                                              // 笔记url
	Hists    []Note `json:"hists"`                                            // 图片历史记录
	Sub      int8   `json:"sub" form:"sub"`                                   // 是否包含所有子集，0-否，1-是
	Deleted  int8   `json:"deleted" form:"deleted"`                           // 是否包含已删除文件，0-否，1-是
	Children []Note `json:"children"`                                         // 子集
}

// Img 图片
type Img struct {
	Abs
	Name     string `json:"name" form:"name" binding:"required,min=1,max=60"` // 图片名称
	Type     string `json:"type" form:"type"`                                 // 图片类型
	Size     int64  `json:"size" form:"size"`                                 // 图片大小，单位：byte
	Hist     string `json:"hist" form:"hist"`                                 // history（历史记录）
	HistSize int64  `json:"histSize" form:"histSize"`                         // history（历史记录）文件大小，单位：byte
	Url      string `json:"url"`                                              // 图片url
	Hists    []Img  `json:"hists"`                                            // 图片历史记录
}

// Stat 统计（file和img）
type Stat struct {
	Type     string `json:"type"`     // 文件类型
	Num      int64  `json:"num"`      // 文件数量
	Size     int64  `json:"size"`     // 文件大小
	HistSize int64  `json:"histSize"` // 文件历史大小
}

// 注册模型
func init() {
	gob.Register(User{})
	gob.Register(Note{})
	gob.Register(Img{})
	gob.Register(Stat{})
}
