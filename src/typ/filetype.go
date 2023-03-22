// file type
// @author xiangqian
// @date 21:14 2023/03/22
package typ

import "strings"

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
