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

// Abs 抽象类型定义
type Abs struct {
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
	Abs
}

// Audio 音频
type Audio struct {
	Abs
}

// Video 视频
type Video struct {
	Abs
}

// Note 笔记
type Note struct {
	Abs
	Pid       int64    `json:"pid"`       // 父id
	PidsStr   string   `json:"pidsStr"`   // 父节点id集字符串
	PnamesStr string   `json:"pnamesStr"` // 父节点名称集字符串
	Pids      []string `json:"pids"`      // 父节点id集
	Pnames    []string `json:"pnames"`    // 父节点名称集
	Content   string   `json:"content"`   // 文件内容
}

// PNote 笔记父节点
type PNote struct {
	Id       int64    `json:"id"`       // 父节点id
	IdsStr   string   `json:"idsStr"`   // 父节点id集字符串
	NamesStr string   `json:"namesStr"` // 父节点名称集字符串
	Ids      []string `json:"ids"`      // 父节点id集
	Names    []string `json:"names"`    // 父节点名称集
	C        bool     `json:"c"`        // contain & child，是否包含子目录
}

// Stats 统计
type Stats struct {
	Type  string `json:"type"`  // 文件类型
	Count int64  `json:"count"` // 文件数量
	Size  int64  `json:"size"`  // 文件大小
}

// 注册模型
func init() {
	gob.Register(System{})
	gob.Register(Image{})
	gob.Register(Audio{})
	gob.Register(Video{})
	gob.Register(PNote{})
	gob.Register(Note{})
	gob.Register(Stats{})
}
