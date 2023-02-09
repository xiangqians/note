// api type
// @author xiangqian
// @date 14:07 2023/02/04
package typ

import "encoding/gob"

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

// User 用户
type User struct {
	Abs
	Name     string `form:"name" binding:"required,excludes= ,min=1,max=60"`               // 用户名
	Nickname string `form:"nickname"binding:"max=60"`                                      // 昵称
	Passwd   string `form:"passwd" binding:"required,excludes= ,max=100"`                  // 密码
	RePasswd string `form:"rePasswd" binding:"required,excludes= ,max=100,eqfield=Passwd"` // retype Passwd
}

// File 文件
type File struct {
	Abs
	Pid  int64  `form:"pid" binding:"gte=0"`                  // 父id
	Name string `form:"name" binding:"required,min=1,max=60"` // 文件名称
	Type string `form:"type"`                                 // 文件类型
	Size int64  `form:"size"`                                 // 文件大小，单位：byte

	Path string // 目录路径
}

// FileType 文件类型
type FileType string

const (
	FileTypeD   FileType = "d"   // 目录
	FileTypeMd           = "md"  // md文件
	FileTypePdf          = "pdf" // pdf文件
	FileTypeZip          = "zip" // zip文件
	FileTypeUnk          = "unk" // unknown
)

var fileTypes = [...]FileType{FileTypeD, FileTypeMd, FileTypePdf, FileTypeZip}

func FileTypeOf(value string) FileType {
	for _, fileType := range fileTypes {
		if string(fileType) == value {
			return fileType
		}
	}

	return FileTypeUnk
}

// 注册模型
func init() {
	gob.Register(User{})
	gob.Register(File{})
}
