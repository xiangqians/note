// @author xiangqian
// @date 11:08 2023/02/04
package os

import (
	"errors"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// IsWindows 是否是windows系统
var IsWindows = false

// IsLinux 是否是linux系统
var IsLinux = false

func init() {
	switch runtime.GOOS {
	// windows
	case "windows":
		IsWindows = true

	// linux
	case "linux":
		fallthrough // 执行穿透
	case "android":
		IsLinux = true

	// 不支持当前操作系统
	default:
		panic(fmt.Sprintf("不支持当前操作系统：%s", runtime.GOOS))
	}
}

// Path 拼接路径
func Path(more ...string) string {
	if IsWindows {
		return strings.Join(more, "\\")
	}

	if IsLinux {
		return strings.Join(more, "/")
	}

	panic(errors.New(fmt.Sprintf("Not Implemented, %s", runtime.GOOS)))
}

// Cmd 执行命令
func Cmd(cmd string) (*exec.Cmd, error) {
	if IsWindows {
		return exec.Command("cmd", "/C", cmd), nil
	}

	if IsLinux {
		return exec.Command("bash", "-c", cmd), nil
	}

	panic(errors.New(fmt.Sprintf("不支持当前操作系统：%s", runtime.GOOS)))
}

// Cd 执行cd命令
func Cd(path string) (*exec.Cmd, error) {
	if IsWindows {
		return Cmd(fmt.Sprintf("cd /d %s", path))
	}

	if IsLinux {
		return Cmd(fmt.Sprintf("cd %s", path))
	}

	panic(errors.New(fmt.Sprintf("不支持当前操作系统：%s", runtime.GOOS)))
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
	if IsWindows {
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(buf)
		return string(decodeBytes)
	}

	return string(buf)
}

type Byte float64

const (
	B  Byte = 1         // B（Byte），字节，1B = 8b（bit）
	KB      = 1024 * B  // KB(Kilobyte)，千字节
	MB      = 1024 * KB // MB（Megabyte），兆字节
	GB      = 1024 * MB // GB（Gigabyte），千兆字节
	TB      = 1024 * GB // TB（Terabyte）
	PB      = 1024 * TB // PB（Petabyte）
)

func ParseByte(s string) (Byte, error) {
	s = strings.ToUpper(s)
	if strings.HasSuffix(s, "PB") {
		f, err := strconv.ParseFloat(s[:len(s)-2], 64)
		return Byte(f) * PB, err
	}

	if strings.HasSuffix(s, "TB") {
		f, err := strconv.ParseFloat(s[:len(s)-2], 64)
		return Byte(f) * TB, err
	}

	if strings.HasSuffix(s, "GB") {
		f, err := strconv.ParseFloat(s[:len(s)-2], 64)
		return Byte(f) * GB, err
	}

	if strings.HasSuffix(s, "MB") {
		f, err := strconv.ParseFloat(s[:len(s)-2], 64)
		return Byte(f) * MB, err
	}

	if strings.HasSuffix(s, "KB") {
		f, err := strconv.ParseFloat(s[:len(s)-2], 64)
		return Byte(f) * KB, err
	}

	if strings.HasSuffix(s, "B") {
		f, err := strconv.ParseFloat(s[:len(s)-1], 64)
		return Byte(f) * B, err
	}

	f, err := strconv.ParseFloat(s[:len(s)-1], 64)
	return Byte(f), err
}

func (byte Byte) String() string {
	if byte <= 0 {
		return "0 B"
	}

	// GB
	gb := byte / GB
	if gb >= 1 {
		return fmt.Sprintf("%.2f GB", gb)
	}

	// MB
	mb := byte / MB
	if mb >= 1 {
		return fmt.Sprintf("%.2f MB", mb)
	}

	// KB
	kb := byte / KB
	if kb >= 1 {
		return fmt.Sprintf("%.2f KB", kb)
	}

	// B
	return fmt.Sprintf("%f B", byte)
}
