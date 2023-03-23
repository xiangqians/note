// api type
// @author xiangqian
// @date 14:07 2023/02/04
package typ

import (
	"encoding/gob"
	"strings"
)

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

	Path     string // 目录路径
	PathLink string // 目录路径链接
}

// Img 图片
type Img struct {
	Abs
	Name string `form:"name" binding:"required,min=1,max=60"` // 图片名称
	Type string `form:"type"`                                 // 图片类型
	Size int64  `form:"size"`                                 // 图片大小，单位：byte

	Url string // 图片url
}

// FileType 文件类型
type FileType string

const (
	FileTypeD    FileType = "d"    // 目录
	FileTypeMd            = "md"   // md文件
	FileTypeHtml          = "html" // html文件
	FileTypePdf           = "pdf"  // pdf文件
	FileTypeZip           = "zip"  // zip文件
	FileTypeIco           = "ico"  // ico文件
	FileTypeGif           = "gif"  // gif文件
	FileTypeJpg           = "jpg"  // jpg文件
	FileTypeJpeg          = "jpeg" // jpeg文件
	FileTypePng           = "png"  // png文件
	FileTypeWebp          = "webp" //webp文件
	FileTypeUnk           = "unk"  // unknown
)

var fileTypes = [...]FileType{
	FileTypeD,
	FileTypeMd, FileTypeHtml, FileTypePdf, FileTypeZip,
	FileTypeIco, FileTypeGif, FileTypeJpg, FileTypeJpeg, FileTypePng, FileTypeWebp,
}

func FileTypeOf(value string) FileType {
	for _, fileType := range fileTypes {
		if strings.ToLower(string(fileType)) == strings.ToLower(value) {
			return fileType
		}
	}

	return FileTypeUnk
}

var imgFileTypes = [...]FileType{FileTypeIco, FileTypeGif, FileTypeJpg, FileTypeJpeg, FileTypePng, FileTypeWebp}

func FileTypeImgOf(value string) FileType {
	for _, imgFileType := range imgFileTypes {
		if strings.ToLower(string(imgFileType)) == strings.ToLower(value) {
			return imgFileType
		}
	}

	return FileTypeUnk
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
	gob.Register(File{})
	gob.Register(Img{})
}
