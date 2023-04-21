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
	typ_os "note/src/typ/os"
	"os"
	"os/exec"
	"runtime"
)

var fileSeparator string
var osTyp typ_os.OS

func init() {
	// init os type
	switch runtime.GOOS {
	// windows
	case "windows":
		osTyp = typ_os.Windows

	// linux
	case "linux":
		fallthrough // 执行穿透
	case "android":
		osTyp = typ_os.Linux

	// unknown
	default:
		osTyp = typ_os.Unk
	}

	// init file separator
	switch osTyp {
	case typ_os.Windows:
		fileSeparator = "\\"

	case typ_os.Linux:
		fileSeparator = "/"

	default:
		fileSeparator = "/"
	}
}

// FileSeparator 文件分隔符
func FileSeparator() string {
	return fileSeparator
}

// OS 获取操作系统标识
func OS() typ_os.OS {
	return osTyp
}

// NotSupportedError 不支持当前系统错误
func NotSupportedError() error {
	return errors.New(fmt.Sprintf("The current system is not supported, %v", runtime.GOOS))
}

// Cmd 执行命令
func Cmd(cmd string) (*exec.Cmd, error) {
	switch OS() {
	case typ_os.Windows:
		return exec.Command("cmd", "/C", cmd), nil

	case typ_os.Linux:
		return exec.Command("bash", "-c", cmd), nil

	default:
		return nil, NotSupportedError()
	}
}

// Cd 执行cd命令
func Cd(path string) (*exec.Cmd, error) {
	switch OS() {
	case typ_os.Windows:
		return Cmd(fmt.Sprintf("cd /d %s", path))

	case typ_os.Linux:
		return Cmd(fmt.Sprintf("cd %s", path))

	default:
		return nil, NotSupportedError()
	}
}

// IsExist 判断所给路径（文件/文件夹）是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// DelFile 删除文件
func DelFile(path string) error {
	return os.Remove(path)
}

// DelDir 删除文件夹
func DelDir(path string) error {
	return os.RemoveAll(path)
}

// MkDir make directories, 创建目录
func MkDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// CopyDir 拷贝目录
func CopyDir(srcDir, dstDir string) (*exec.Cmd, error) {
	switch OS() {
	case typ_os.Windows:
		return Cmd(fmt.Sprintf("xcopy %s %s /s /e /h /i /y", srcDir, dstDir))

	case typ_os.Linux:
		return Cmd(fmt.Sprintf("cp -r %s %s", srcDir, dstDir))

	default:
		return nil, NotSupportedError()
	}
}

func DecodeBuf(buf []byte) string {
	if buf == nil || len(buf) == 0 {
		return ""
	}

	switch OS() {
	// 解决windows乱码问题
	// GB18030编码
	case typ_os.Windows:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(buf)
		return string(decodeBytes)

	default:
		return string(buf)
	}
}

// CopyFile 拷贝文件
func CopyFile(srcPath, dstPath string) error {
	// src
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	// dst
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// copy
	return CopyIo(src, dst, 0)
}

// CopyIo 流拷贝
// src: io.Reader
// dst: io.Writer
// bufSize: 缓存大小，byte
func CopyIo(src io.Reader, dst io.Writer, bufSize int) error {
	var reader *bufio.Reader
	var writer *bufio.Writer

	if bufSize <= 0 {
		bufSize = 1024 * 4 // bufio.defaultBufSize
	}

	// 块缓存大小
	buf := make([]byte, bufSize)

	// <= 4KB
	if bufSize <= 1024*4 {
		reader = bufio.NewReader(src)
		writer = bufio.NewWriter(dst)
	} else
	// > 4KB
	{
		reader = bufio.NewReaderSize(src, bufSize)
		writer = bufio.NewWriterSize(dst, bufSize)
	}

	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			if n > 0 {
				writer.Write(buf[:n])
				writer.Flush()
			}
			break
		}

		if err != nil {
			return err
		}

		writer.Write(buf[:n])
		writer.Flush()
	}

	return nil
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
	return fmt.Sprintf("%v B", size)
}
