// io
// @author xiangqian
// @date 22:31 2022/12/20
package util

import (
	"bufio"
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
