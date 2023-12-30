// @author xiangqian
// @date 20:00 2023/12/04
package filetype

import (
	"net/http"
	"strings"
)

// 文件类型
const (
	Folder = "folder" // folder
	Md     = "md"     // md
	Doc    = "doc"    // doc
	Docx   = "docx"   // docx
	Pdf    = "pdf"    // pdf
	Zip    = "zip"    // zip
	TarGz  = "tar.gz" // tar.gz
	Ico    = "ico"    // ico
	Gif    = "gif"    // gif
	Jpg    = "jpg"    // jpg
	Jpeg   = "jpeg"   // jpeg
	Png    = "png"    // png
	Webp   = "webp"   // webp
	Mp3    = "mp3"    // mp3
	Wav    = "wav"    // wav
	Flac   = "flac"   // flac
	Aac    = "aac"    // aac
	Ogg    = "ogg"    // Ogg Vorbis
	Mp4    = "mp4"    // mp4
	Avi    = "avi"    // Audio Video Interleaved
	Mov    = "mov"    // QuickTime
	Wmv    = "wmv"    // Windows Media Video
	Mkv    = "mkv"    // Matroska Video
	Flv    = "flv"    // Flash Video
)

func GetType(name string) string {
	if name == "" {
		return ""
	}

	index := strings.Index(name, ".")
	var suffix string
	if index >= 0 {
		suffix = name[index:]
	} else {
		suffix = name
	}
	suffix = strings.ToLower(suffix)

	if strings.HasSuffix(suffix, Doc) {
		return Doc

	} else if strings.HasSuffix(suffix, Docx) {
		return Docx

	} else if strings.HasSuffix(suffix, Pdf) {
		return Pdf

	} else if strings.HasSuffix(suffix, Zip) {
		return Zip

	} else if strings.HasSuffix(suffix, TarGz) {
		return TarGz

	} else
	// image
	if strings.HasSuffix(suffix, Ico) {
		return Ico

	} else if strings.HasSuffix(suffix, Gif) {
		return Gif

	} else if strings.HasSuffix(suffix, Jpg) {
		return Jpg

	} else if strings.HasSuffix(suffix, Jpeg) {
		return Jpeg

	} else if strings.HasSuffix(suffix, Png) {
		return Png

	} else if strings.HasSuffix(suffix, Webp) {
		return Webp

	} else
	// audio
	if strings.HasSuffix(suffix, Mp3) {
		return Mp3

	} else if strings.HasSuffix(suffix, Wav) {
		return Wav

	} else if strings.HasSuffix(suffix, Flac) {
		return Flac

	} else if strings.HasSuffix(suffix, Aac) {
		return Aac

	} else if strings.HasSuffix(suffix, Ogg) {
		return Ogg
	} else
	// video
	if strings.HasSuffix(suffix, Mp4) {
		return Mp4

	} else if strings.HasSuffix(suffix, Avi) {
		return Avi

	} else if strings.HasSuffix(suffix, Mov) {
		return Mov

	} else if strings.HasSuffix(suffix, Wmv) {
		return Wmv

	} else if strings.HasSuffix(suffix, Mkv) {
		return Mkv

	} else if strings.HasSuffix(suffix, Flv) {
		return Flv
	}

	return suffix
}

func GetTypeBak(data []byte) string {
	// Go标准库提供一个基于 mimesniff 算法的 http.DetectContentType 函数，只需要读取文件的前512个字节就能够判定文件类型。
	// 请注意，这种方法并不是绝对准确的，因为文件头部信息可能会被修改或伪造。
	contentType := http.DetectContentType(data)
	switch contentType {
	case "application/pdf":
		return Pdf
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
		return contentType
	}
}
