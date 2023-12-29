// @author xiangqian
// @date 20:00 2023/12/04
package filetype

import (
	"net/http"
	"strings"
)

// 文件类型
const (
	Folder = "folder" // 文件夹
	Md     = "md"     // md文件
	Doc    = "doc"    // doc文件
	Docx   = "docx"   // docx文件
	Pdf    = "pdf"    // pdf文件
	Zip    = "zip"    // zip文件
	TarGz  = "tar.gz" // tar.gz文件
	Ico    = "ico"    // ico文件
	Gif    = "gif"    // gif文件
	Jpg    = "jpg"    // jpg文件
	Jpeg   = "jpeg"   // jpeg文件
	Png    = "png"    // png文件
	Webp   = "webp"   // webp文件
)

func GetType(name string, data []byte) string {
	// Go标准库提供一个基于 mimesniff 算法的 http.DetectContentType 函数，只需要读取文件的前512个字节就能够判定文件类型。
	// 请注意，这种方法并不是绝对准确的，因为文件头部信息可能会被修改或伪造。
	contentType := http.DetectContentType(data)
	switch contentType {
	case "text/plain; charset=utf-8":
		if name != "" {
			if strings.HasSuffix(name, ".doc") {
				return Doc
			} else if strings.HasSuffix(name, ".docx") {
				return Docx
			}
		}
	case "application/pdf":
		return Pdf
	case "application/zip", "application/x-gzip":
		if name != "" && strings.HasSuffix(name, ".tar.gz") {
			return TarGz
		}
		return Zip
	case "image/x-icon":
		return Ico
	case "image/gif":
		return Gif
	case "image/jpg":
		return Jpg
	case "image/jpeg":
		return Jpeg
	case "image/png":
		return Png
	case "image/webp":
		return Webp
	}
	return contentType
}
