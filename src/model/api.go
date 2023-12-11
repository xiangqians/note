// @author xiangqian
// @date 11:53 2023/05/07
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

// 抽象类型定义
type abs struct {
	Id      int64  `json:"id"`      // 主键id
	Name    string `json:"name"`    // 名称
	Type    string `json:"type"`    // 类型
	Size    int64  `json:"size"`    // 大小，单位：byte
	Del     byte   `json:"del"`     // 删除标识，0-正常，1-删除，2-永久删除
	AddTime int64  `json:"addTime"` // 创建时间（时间戳，单位s）
	UpdTime int64  `json:"updTime"` // 修改时间（时间戳，单位s）
}

// Image 图片
type Image struct {
	abs
}

// Audio 音频
type Audio struct {
	abs
}

// Video 视频
type Video struct {
	abs
}

// Note 笔记
type Note struct {
	abs
	Pid int64 `json:"pid"` // 父id
}

// PNote 笔记父节点
type PNote struct {
	Id       int64    `json:"id"`       // 父节点id
	IdsStr   string   `json:"idsStr"`   // 父节点id集
	NamesStr string   `json:"namesStr"` // 父节点名称集
	Ids      []string `json:"ids"`      // 父节点id集
	Names    []string `json:"names"`    // 父节点名称集
}

// 注册模型
func init() {
	gob.Register(System{})
	gob.Register(Image{})
	gob.Register(Audio{})
	gob.Register(Video{})
	gob.Register(Note{})
}
