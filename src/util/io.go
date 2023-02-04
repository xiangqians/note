// io
// @author xiangqian
// @date 22:31 2022/12/20
package util

import (
	"bufio"
	"io"
)

// Copy
// bufsize: 缓存大小，byte
func Copy(dst io.Writer, src io.Reader, bufsize int) error {
	var pReader *bufio.Reader
	var pWriter *bufio.Writer

	if bufsize <= 0 {
		bufsize = 1024 * 4 // bufio.defaultBufSize
	}

	// 块缓存大小
	buf := make([]byte, bufsize)

	// <= 4KB
	if bufsize <= 1024*4 {
		pReader = bufio.NewReader(src)
		pWriter = bufio.NewWriter(dst)
	} else
	// > 4KB
	{
		pReader = bufio.NewReaderSize(src, bufsize)
		pWriter = bufio.NewWriterSize(dst, bufsize)
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
