// api type
// @author xiangqian
// @date 14:07 2023/02/04
package typ

// Ft file type，文件类型
type Ft string

const (
	FtUnk  Ft = ""     // unknown
	FtD       = "d"    // 目录
	FtMd      = "md"   // md文件
	FtHtml    = "html" // html文件
	FtPdf     = "pdf"  // pdf文件
	FtZip     = "zip"  // zip文件
	FtIco     = "ico"  // ico文件
	FtGif     = "gif"  // gif文件
	FtJpg     = "jpg"  // jpg文件
	FtJpeg    = "jpeg" // jpeg文件
	FtPng     = "png"  // png文件
	FtWebp    = "webp" // webp文件
)

var fts = [...]Ft{
	FtD,
	FtMd, FtHtml, FtPdf, FtZip,
	FtIco, FtGif, FtJpg, FtJpeg, FtPng, FtWebp,
}

// ExtNameOf 根据文件扩展名（Filename Extension）获取文件类型
func ExtNameOf(extName string) Ft {
	for _, ft := range fts {
		if string(ft) == extName {
			return ft
		}
	}

	return FtUnk
}

// ContentTypeOf 根据contentType获取文件类型
func ContentTypeOf(contentType string) Ft {
	switch contentType {
	// img
	case "image/x-icon":
		return FtIco
	case "image/gif":
		return FtGif
	case "image/jpg":
		return FtJpg
	case "image/jpeg":
		return FtJpeg
	case "image/png":
		return FtPng
	case "image/webp":
		return FtWebp
	// note
	case "text/html":
		return FtHtml
	case "application/pdf":
		return FtPdf
	case "application/x-zip-compressed":
		return FtZip
	// unknown
	default:
		return FtUnk
	}
}

var imgFts = [...]Ft{FtIco, FtGif, FtJpg, FtJpeg, FtPng, FtWebp}

// IsImg 是否是图片文件类型
func IsImg(ft Ft) bool {
	for _, imgFt := range imgFts {
		if imgFt == ft {
			return true
		}
	}
	return false
}
