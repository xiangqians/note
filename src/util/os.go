// os
// @author xiangqian
// @date 11:08 2023/02/04
package util

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// OperatingSystem 操作系统标识
type OperatingSystem int8

const (
	OSWindows OperatingSystem = iota // Windows
	OSLinux                          // Linux
	OSUnknown                        // Unknown
)

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

// Command 执行命令行
func Command(cmd string) (*exec.Cmd, error) {
	switch OS() {
	case OSWindows:
		return exec.Command("cmd", "/C", cmd), nil

	case OSLinux:
		return exec.Command("bash", "-c", cmd), nil

	default:
		return nil, errors.New(fmt.Sprintf("The current system is not supported, %v", runtime.GOOS))
	}
}
