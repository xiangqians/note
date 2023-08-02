// os
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
	"time"
)

type File interface {
	Err() error         // 文件错误
	IsExist() bool      // 判断文件（普通文件、目录文件）是否存在
	Name() string       // base name of the file
	Size() int64        // length in bytes for regular files; system-dependent for others
	Mode() os.FileMode  // file mode bits
	ModTime() time.Time // modification time
	IsDir() bool        // abbreviation for Mode().IsDir()
	Sys() any           // underlying data source (can return nil)
}

type fileImpl struct {
	fileInfo os.FileInfo
	err      error
}

func (file *fileImpl) Err() error {
	return file.err
}

func (file *fileImpl) IsExist() bool {
	if file.err != nil {
		return os.IsExist(file.err)
	}
	return true
}

func (file *fileImpl) Name() string {
	return file.fileInfo.Name()
}

func (file *fileImpl) Size() int64 {
	return file.fileInfo.Size()
}

func (file *fileImpl) Mode() os.FileMode {
	return file.fileInfo.Mode()
}

func (file *fileImpl) ModTime() time.Time {
	return file.fileInfo.ModTime()
}

func (file *fileImpl) IsDir() bool {
	return file.fileInfo.IsDir()
}

func (file *fileImpl) Sys() any {
	return file.fileInfo.Sys()
}

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

// CurrentOsIsNotSupportedError 不支持当前操作系统错误
func CurrentOsIsNotSupportedError() error {
	return errors.New(fmt.Sprintf("The current os is not supported, %s", runtime.GOOS))
}

func Path(more ...string) string {
	if isWindows {
		return strings.Join(more, "\\")
	}

	if isLinux {
		return strings.Join(more, "/")
	}

	panic(CurrentOsIsNotSupportedError())
}

// Stat 文件信息
func Stat(path string) File {
	fileInfo, err := os.Stat(path)
	return &fileImpl{fileInfo, err}
}

// MkDir (make directories) 创建目录文件
func MkDir(path string) error {
	file := Stat(path)
	if file.IsExist() {
		return errors.New(fmt.Sprintf("File exists, %s", path))
	}

	return os.MkdirAll(path, os.ModePerm) // 0777
}

// Rm 删除文件（普通文件或者目录文件）
func Rm(path string) error {
	file := Stat(path)
	if !file.IsExist() {
		return errors.New(fmt.Sprintf("File not found, %s", path))
	}

	// 删除目录文件
	if file.IsDir() {
		return os.RemoveAll(path)
	}

	// 删除普通文件
	return os.Remove(path)
}

// CopyFile 拷贝文件
func CopyFile(srcPath, dstPath string) (written int64, err error) {
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
	written, err = io.Copy(dst, src)
	return
}

// Deprecated: use io.Copy
// CopyIo 拷贝数据流
// src: io.Reader
// dst: io.Writer
// bufSize: 缓存大小，byte。默认 bufio.defaultBufSize = 4KB
func CopyIo(src io.Reader, dst io.Writer, bufSize int) (int, error) {
	// buf size
	if bufSize <= 0 {
		bufSize = 1024 * 4 // 4KB
	}

	// r & w
	reader := bufio.NewReaderSize(src, bufSize)
	writer := bufio.NewWriterSize(dst, bufSize)

	// 已写入的字节数
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

// Cmd 执行命令
func Cmd(cmd string) (*exec.Cmd, error) {
	if isWindows {
		return exec.Command("cmd", "/C", cmd), nil
	}

	if isLinux {
		return exec.Command("bash", "-c", cmd), nil
	}

	return nil, CurrentOsIsNotSupportedError()
}

// Cd 执行cd命令
func Cd(path string) (*exec.Cmd, error) {
	if isWindows {
		return Cmd(fmt.Sprintf("cd /d %s", path))
	}

	if isLinux {
		return Cmd(fmt.Sprintf("cd %s", path))
	}

	return nil, CurrentOsIsNotSupportedError()
}

// CopyDir 拷贝目录
func CopyDir(srcPath, dstPath string) (*exec.Cmd, error) {
	if isWindows {
		return Cmd(fmt.Sprintf("xcopy %s %s /s /e /h /i /y", srcPath, dstPath))
	}

	if isLinux {
		return Cmd(fmt.Sprintf("cp -r %s %s", srcPath, dstPath))
	}

	return nil, CurrentOsIsNotSupportedError()
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

	// 解决windows乱码问题
	// GB18030编码
	if isWindows {
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(buf)
		return string(decodeBytes)
	}

	return string(buf)
}
