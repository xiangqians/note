// @author xiangqian
// @date 10:24 2023/10/28
package common

import "note/src/util/filetype"

// 图片文件类型
var imageFileTypes = [...]string{
	filetype.Ico,
	filetype.Gif,
	filetype.Jpg,
	filetype.Jpeg,
	filetype.Png,
	filetype.Webp,
}

// 是否是图片文件类型
func isImage(fileType string) bool {
	for _, imageFileType := range imageFileTypes {
		if imageFileType == fileType {
			return true
		}
	}
	return false
}
