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

func Get(request *http.Request, writer http.ResponseWriter, session *session.Session, table string) {
	vars := mux.Vars(request)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return
	}

	db := db.Get()
	result, err := db.Get(fmt.Sprintf("SELECT `id`, `name`, `type`, `size`, `del`, `add_time`, `upd_time` FROM `%s` WHERE `del` = 0 AND `id` = ? LIMIT 1", table), id)
	if err != nil {
		return
	}

	var name string
	var filetype string
	switch table {
	case TableImage:
		var image model.Image
		err = result.Scan(&image)
		id = image.Id
		name = image.Name
		filetype = image.Type
	}

	if err != nil || id == 0 {
		return
	}

	file, err := os.Open(util_os.Path(dataDir, table, fmt.Sprintf("%d", id)))
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

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.%s", url.QueryEscape(name), filetype))
	writer.Header().Set("Content-Type", contentType)
	writer.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))
	//writer.Header().Set("Content-Transfer-Encoding", "binary")

	file.Seek(0, 0)
	_, err = io.Copy(writer, file)

	return
}
