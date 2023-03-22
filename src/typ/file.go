// file
// @author xiangqian
// @date 21:12 2023/03/22
package typ

import "encoding/gob"

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

// 注册模型
func init() {
	gob.Register(File{})
}
