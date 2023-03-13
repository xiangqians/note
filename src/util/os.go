// os
// @author xiangqian
// @date 11:08 2023/02/04
package util

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"note/src/typ"
	"os"
	"os/exec"
	"regexp"
	"runtime"
)

// FileSeparator 文件分隔符
var FileSeparator string

func init() {
	switch OS() {
	case typ.WindowsOS:
		FileSeparator = "\\"

	case typ.LinuxOS:
		FileSeparator = "/"

	default:
		FileSeparator = "/"
	}
}

// OS 获取操作系统标识
func OS() typ.OS {
	os := runtime.GOOS
	switch os {
	// windows
	case "windows":
		return typ.WindowsOS

	// linux
	case "linux":
		fallthrough // 执行穿透
	case "android":
		return typ.LinuxOS

	// unknown
	default:
		return typ.UnknownOS
	}
}

// CdCmd cd命令
func CdCmd(path string) (string, error) {
	switch OS() {
	case typ.WindowsOS:
		return fmt.Sprintf("cd /d %s", path), nil

	case typ.LinuxOS:
		return fmt.Sprintf("cd %s", path), nil

	default:
		return "", errors.New(fmt.Sprintf("The current system is not supported, %v", runtime.GOOS))
	}
}

// IsExistOfPath 判断所给路径（文件/文件夹）是否存在
func IsExistOfPath(path string) bool {
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

// Mkdir 创建目录
func Mkdir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func CommandWindows(cmd string) *exec.Cmd {
	return exec.Command("cmd", "/C", cmd)
}

func CommandLinux(cmd string) *exec.Cmd {
	return exec.Command("bash", "-c", cmd)
}

// Command 执行命令行
func Command(cmd string) (*exec.Cmd, error) {
	switch OS() {
	case typ.WindowsOS:
		return CommandWindows(cmd), nil

	case typ.LinuxOS:
		return CommandLinux(cmd), nil

	default:
		return nil, errors.New(fmt.Sprintf("The current system is not supported, %v", runtime.GOOS))
	}
}

func CopyDir(srcDir, dstDir string) *exec.Cmd {
	switch OS() {
	case typ.WindowsOS:
		return CommandWindows(fmt.Sprintf("xcopy %s %s /s /e /h /i /y", srcDir, dstDir))

	case typ.LinuxOS:
		return CommandLinux(fmt.Sprintf("cp -r %s %s", srcDir, dstDir))

	default:
		panic(fmt.Sprintf("The current system is not supported, %v", runtime.GOOS))
	}
}

func DecodeBuf(buf []byte) string {
	if buf == nil || len(buf) == 0 {
		return ""
	}

	switch OS() {
	// 解决windows乱码问题
	// GB18030编码
	case typ.WindowsOS:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(buf)
		return string(decodeBytes)

	default:
		return string(buf)
	}
}

// IOCopy 流拷贝
// src: io.Reader
// dst: io.Writer
// bufSize: 缓存大小，byte
func IOCopy(src io.Reader, dst io.Writer, bufSize int) error {
	var pReader *bufio.Reader
	var pWriter *bufio.Writer

	if bufSize <= 0 {
		bufSize = 1024 * 4 // bufio.defaultBufSize
	}

	// 块缓存大小
	buf := make([]byte, bufSize)

	// <= 4KB
	if bufSize <= 1024*4 {
		pReader = bufio.NewReader(src)
		pWriter = bufio.NewWriter(dst)
	} else
	// > 4KB
	{
		pReader = bufio.NewReaderSize(src, bufSize)
		pWriter = bufio.NewWriterSize(dst, bufSize)
	}

	for {
		n, err := pReader.Read(buf)
		if err == io.EOF {
			if n > 0 {
				pWriter.Write(buf[:n])
				pWriter.Flush()
			}
			break
		}

		if err != nil {
			return err
		}

		pWriter.Write(buf[:n])
		pWriter.Flush()
	}

	return nil
}

// HumanizFileSize 人性化文件大小
// size: 文件大小，单位：byte
func HumanizFileSize(size int64) string {

	// 1B  = 8b
	// 1KB = 1024B
	// 1MB = 1024KB
	// 1GB = 1024MB
	// 1TB = 1024GB

	if size <= 0 {
		return "0 B"
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

// VerifyFileName 校验文件名
func VerifyFileName(dirName string) error {
	// 文件名不能包含字符：
	// \ / : * ? " < > |

	// ^[^\\/:*?"<>|]*$
	matched, err := regexp.MatchString("^[^\\\\/:*?\"<>|]*$", dirName)
	if err != nil {
		return err
	}

	if !matched {
		return errors.New("文件名不能包含字符：\\ / : * ? \" < > |")
	}

	return nil
}
