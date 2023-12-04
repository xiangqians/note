// @author xiangqian
// @date 20:00 2023/12/04
package filetype

import "net/http"

// 文件类型
const (
	Folder  = "folder"  // 文件夹
	Md      = "md"      // md文件
	Html    = "html"    // html文件
	Pdf     = "pdf"     // pdf文件
	Doc     = "doc"     // doc文件
	Zip     = "zip"     // zip文件
	Ico     = "ico"     // ico文件
	Gif     = "gif"     // gif文件
	Jpg     = "jpg"     // jpg文件
	Jpeg    = "jpeg"    // jpeg文件
	Png     = "png"     // png文件
	Webp    = "webp"    // webp文件
	Unknown = "unknown" // 未知
)

func GetType(data []byte) string {
	// Go标准库提供一个基于 mimesniff 算法的 http.DetectContentType 函数，只需要读取文件的前512个字节就能够判定文件类型
	contentType := http.DetectContentType(data)
	switch contentType {
	case "text/html":
		return Html
	case "application/pdf":
		return Pdf
	case "application/?":
		return Doc
	case "application/zip", "application/x-gzip":
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
	default:
		return Unknown
	}
}
