// io
// @author xiangqian
// @date 22:31 2022/12/20
package util

import (
	"bufio"
	"fmt"
	"io"
)

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
