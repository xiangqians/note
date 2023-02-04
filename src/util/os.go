// os
// @author xiangqian
// @date 11:08 2023/02/04
package util

import (
	"errors"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os"
	"os/exec"
	"runtime"
)

var FileSeparator string

// OperatingSystem 操作系统标识
type OperatingSystem int8

const (
	OSWindows OperatingSystem = iota // Windows
	OSLinux                          // Linux
	OSUnknown                        // Unknown
)

func init() {
	switch OS() {
	case OSWindows:
		FileSeparator = "\\"

	case OSLinux:
		FileSeparator = "/"

	default:
		FileSeparator = "/"
	}
}

func OS() OperatingSystem {
	if runtime.GOOS == "windows" {
		return OSWindows
	}

	if runtime.GOOS == "linux" {
		return OSLinux
	}

	return OSUnknown
}

func Cd(path string) (string, error) {
	switch OS() {
	case OSWindows:
		return fmt.Sprintf("cd /d %s", path), nil

	case OSLinux:
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
	return exec.Command("cmd", "/C", cmd)
}

// Command 执行命令行
func Command(cmd string) (*exec.Cmd, error) {
	switch OS() {
	case OSWindows:
		return CommandWindows(cmd), nil

	case OSLinux:
		return CommandLinux(cmd), nil

	default:
		return nil, errors.New(fmt.Sprintf("The current system is not supported, %v", runtime.GOOS))
	}
}

func CopyDir(srcDir, dstDir string) *exec.Cmd {
	switch OS() {
	case OSWindows:
		return CommandWindows(fmt.Sprintf("xcopy %s %s /s /e /h /i /y", srcDir, dstDir))

	case OSLinux:
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
	case OSWindows:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(buf)
		return string(decodeBytes)

	default:
		return string(buf)
	}
}
