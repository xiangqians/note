// @author xiangqian
// @date 21:55 2023/12/04
package common

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/url"
	"note/src/db"
	"note/src/model"
	"note/src/session"
	util_os "note/src/util/os"
	"os"
	"strconv"
)

func Get(request *http.Request, writer http.ResponseWriter, session *session.Session, table string) (templateName string, response model.Response) {
	vars := mux.Vars(request)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil || id <= 0 {
		return
	}

	db := db.Get()
	result, err := db.Get(fmt.Sprintf("SELECT `id`, `name`, `type`, `size`, `del`, `add_time`, `upd_time` FROM `%s` WHERE `del` = 0 AND `id` = ? LIMIT 1", table), id)
	if err != nil {
		return
	}

	var abs model.Abs
	err = result.Scan(&abs)
	if err != nil || abs.Id == 0 {
		return
	}

	name := abs.Name
	filetype := abs.Type

	file, err := os.Open(util_os.Path(DataDir, table, fmt.Sprintf("%d", id)))
	if err != nil {
		return
	}
	defer file.Close()

	// Go标准库提供一个基于 mimesniff 算法的 http.DetectContentType 函数，只需要读取文件的前512个字节就能够判定文件类型
	fileHeader := make([]byte, 512)
	_, err = file.Read(fileHeader)
	if err != nil {
		return
	}
	contentType := http.DetectContentType(fileHeader)

	fileStat, err := file.Stat()
	if err != nil {
		return
	}

	fileSize := fileStat.Size()

	// 分段获取数据
	if table == TableAudio || table == TableVideo {
		// 获取请求范围
		start, end, partial := getRange(request, fileSize)

		// 设置响应头
		writer.Header().Set("Accept-Ranges", "bytes")
		writer.Header().Set("Content-Type", contentType)
		writer.Header().Set("Content-Length", strconv.FormatInt(end-start+1, 10))

		// 如果是分段请求，则设置状态码为 206 Partial Content
		if partial {
			writer.WriteHeader(http.StatusPartialContent)
		}

		// 分段传输文件内容
		http.ServeContent(writer, request, "", fileStat.ModTime(), file)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.%s", url.QueryEscape(name), filetype))
	writer.Header().Set("Content-Type", contentType)
	writer.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10))
	//writer.Header().Set("Content-Transfer-Encoding", "binary")

	file.Seek(0, 0)
	_, err = io.Copy(writer, file)

	return
}

func getRange(request *http.Request, fileSize int64) (start, end int64, partial bool) {
	rangeHeader := request.Header.Get("Range")
	if rangeHeader == "" {
		return 0, fileSize - 1, false
	}

	_, err := fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)
	if err != nil {
		return 0, 0, false
	}

	if start >= fileSize || end >= fileSize {
		return 0, 0, false
	}

	return start, end, true
}
