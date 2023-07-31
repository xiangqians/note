// os util
// @author xiangqian
// @date 11:08 2023/02/04
package os

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var isWindows = false
var isLinux = false

func init() {
	switch runtime.GOOS {
	// windows
	case "windows":
		isWindows = true

	// linux
	case "linux":
		fallthrough // 执行穿透
	case "android":
		isLinux = true
	}
}

func IsWindows() bool {
	return isWindows
}

func IsLinux() bool {
	return isLinux
}

func Path(more ...string) string {
	//if isWindows {
	//	return strings.Join(more, "\\")
	//}

	if isLinux {
		return strings.Join(more, "/")
	}

	panic(CurrentOSNotSupportedError())
}

// CurrentOSNotSupportedError 当前操作系统不支持错误
func CurrentOSNotSupportedError() error {
	return errors.New(fmt.Sprintf("The current os is not supported, %v", runtime.GOOS))
}

// Cmd 执行命令
func Cmd(cmd string) (*exec.Cmd, error) {
	if isWindows {
		return exec.Command("cmd", "/C", cmd), nil
	}

	if isLinux {
		return exec.Command("bash", "-c", cmd), nil
	}

	return nil, CurrentOSNotSupportedError()
}

// Cd 执行cd命令
func Cd(path string) (*exec.Cmd, error) {
	if isWindows {
		return Cmd(fmt.Sprintf("cd /d %s", path))
	}

	if isLinux {
		return Cmd(fmt.Sprintf("cd %s", path))
	}

	return nil, CurrentOSNotSupportedError()
}

// IsExist 判断所给路径（文件/文件夹）是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// Rm 删除普通文件或者目录文件
func Rm(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	// 删除目录文件
	if fileInfo.IsDir() {
		return os.RemoveAll(path)
	}

	// 删除普通文件
	return os.Remove(path)
}

// MkDir (make directories) 创建目录文件
func MkDir(path string) error {
	return os.MkdirAll(path, os.ModePerm) // 0777
}

// Cp 拷贝
func Cp(srcPath, dstPath string) (*exec.Cmd, error) {

}

// CopyDir 拷贝目录
func CopyDir(srcPath, dstPath string) (*exec.Cmd, error) {
	if isWindows {
		return Cmd(fmt.Sprintf("xcopy %s %s /s /e /h /i /y", srcPath, dstPath))
	}

	if isLinux {
		return Cmd(fmt.Sprintf("cp -r %s %s", srcPath, dstPath))
	}

	return nil, CurrentOSNotSupportedError()
}

// CopyFile 拷贝文件
func CopyFile(srcPath, dstPath string) (int64, error) {
	// src
	src, err := os.Open(srcPath)
	if err != nil {
		return 0, err
	}
	defer src.Close()

	// dst
	dst, err := os.Create(dstPath)
	if err != nil {
		return 0, err
	}
	defer dst.Close()

	// copy
	return io.Copy(dst, src)
}

// CopyIo 流拷贝
// src: io.Reader
// dst: io.Writer
// bufSize: 缓存大小，byte。默认 bufio.defaultBufSize = 4KB
func CopyIo(dst io.Writer, src io.Reader, bufSize int) (int, error) {
	// buf size
	if bufSize <= 0 {
		bufSize = 1024 * 4 // 4KB
	}

	// w & r
	writer := bufio.NewWriterSize(dst, bufSize)
	reader := bufio.NewReaderSize(src, bufSize)

	// write func
	var written int

	// 块缓存大小
	buf := make([]byte, bufSize)

	// write
	for {
		// Read reads data into buf.
		// It returns the number of bytes read into buf.
		// The bytes are taken from at most one Read on the underlying Reader, hence n may be less than len(buf).
		rn, rerr := reader.Read(buf)

		// write
		if rn > 0 {
			wn, werr := writer.Write(buf[:rn])
			if werr == nil && (wn < 0 || rn < wn) {
				werr = errors.New("invalid write result")
			}

			if werr == nil && rn != wn {
				werr = errors.New("short write")
			}

			if werr != nil {
				return written, werr
			}

			writer.Flush()
			written += wn
		}

		// If the underlying Reader can return a non-zero count with io.EOF,
		// then this Read method can do so as well; see the [io.Reader] docs.
		if rerr == io.EOF {
			break
		}

		// err ?
		if rerr != nil {
			return written, rerr
		}
	}
	return written, nil
}

// HumanizFileSize 人性化文件大小
// size: 文件大小，单位：byte
func HumanizFileSize(size int64) string {

	// B, Byte
	// 1B  = 8b
	// 1KB = 1024B
	// 1MB = 1024KB
	// 1GB = 1024MB
	// 1TB = 1024GB

	if size <= 0 {
		return "0 B"
	}

	// GB
	gb := float64(size) / (1024 * 1024 * 1024)
	if gb > 1 {
		return fmt.Sprintf("%.2f GB", gb)
	}

	// MB
	mb := float64(size) / (1024 * 1024)
	if mb > 1 {
		return fmt.Sprintf("%.2f MB", mb)
	}

	// KB
	kb := float64(size) / 1024
	if kb > 1 {
		return fmt.Sprintf("%.2f KB", kb)
	}

	// B
	return fmt.Sprintf("%d B", size)
}

// DecodeBuf 解码buffer
func DecodeBuf(buf []byte) string {
	if buf == nil || len(buf) == 0 {
		return ""
	}

	switch GetOS() {
	// 解决windows乱码问题
	// GB18030编码
	case Windows:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(buf)
		return string(decodeBytes)

	default:
		return string(buf)
	}
}
