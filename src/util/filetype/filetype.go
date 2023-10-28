// 文件类型
// @author xiangqian
// @date 14:07 2023/02/04
package filetype

// 文件类型
const (
	Folder         = "folder"  // 文件夹
	Md             = "md"      // md文件
	Html           = "html"    // html文件
	Pdf            = "pdf"     // pdf文件
	Doc            = "doc"     // doc文件
	Zip            = "zip"     // zip文件
	Ico            = "ico"     // ico文件
	Gif            = "gif"     // gif文件
	Jpg            = "jpg"     // jpg文件
	Jpeg           = "jpeg"    // jpeg文件
	Png            = "png"     // png文件
	Webp           = "webp"    // webp文件
	Unknown string = "unknown" // 未知
)

// ContentTypeOf 根据文件内容类型获取文件类型
func ContentTypeOf(contentType string) string {
	switch contentType {
	case "text/html":
		return Html
	case "application/pdf":
		return Pdf
	case "application/x-zip-compressed":
		return Zip
	// image
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

	// audio

	//video

	// 未知
	default:
		return Unknown
	}
}
