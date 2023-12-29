// @author xiangqian
// @date 16:22 2023/12/28
package common

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"note/src/db"
	"note/src/model"
	"note/src/session"
	util_filetype "note/src/util/filetype"
	util_os "note/src/util/os"
	util_string "note/src/util/string"
	"note/src/util/time"
	"os"
	"strconv"
	"strings"
)

func ReUpload(request *http.Request, writer http.ResponseWriter, session *session.Session, table string) (string, model.Response) {
	// id
	vars := mux.Vars(request)
	id, err := strconv.ParseInt(vars["id"], 10, 64)

	// 重定向函数
	redirect := func(err any) (string, model.Response) {
		return fmt.Sprintf("redirect:/%s/%d/view", table, id), model.Response{Msg: util_string.String(err)}
	}

	if err != nil || id <= 0 {
		return redirect(err)
	}

	// note.pid
	pid, _ := strconv.ParseInt(request.PostFormValue("pid"), 10, 64)
	if pid < 0 {
		pid = 0
	}

	// 读取上传文件
	file, fileHeader, err := request.FormFile("file")
	if err != nil {
		return redirect(err)
	}
	defer file.Close()

	// 读取文件字节内容
	bytes, err := io.ReadAll(file)
	if err != nil {
		return redirect(err)
	}

	// 文件名
	name := strings.TrimSpace(fileHeader.Filename)

	// 获取文件类型
	filetype := util_filetype.GetType(name, bytes)
	err = validateFiletype(session, table, filetype)
	if err != nil {
		return redirect(err)
	}

	// 去除文件后缀名
	suffix := "." + filetype
	if strings.HasSuffix(name, suffix) {
		name = name[:len(name)-len(suffix)]
	}

	// 文件大小，单位：字节
	size := fileHeader.Size

	var result *db.Result
	db := db.Get()

	// 开启事务
	err = db.Begin()
	if err != nil {
		return redirect(err)
	}

	// 校验id
	result, err = db.Get(fmt.Sprintf("SELECT `id` FROM `%s` WHERE `del` = 0 AND `id` = ? LIMIT 1", table), id)
	if err != nil {
		db.Rollback()
		return redirect(err)
	}

	id = 0
	err = result.Scan(&id)
	if err != nil || id == 0 {
		db.Rollback()
		return redirect(err)
	}

	_, err = db.Upd(fmt.Sprintf("UPDATE `%s` SET `name` = ?, `type` = ?, `size` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ?", table),
		name,
		filetype,
		size,
		time.NowUnix(),
		id)
	if err != nil {
		db.Rollback()
		return redirect(err)
	}

	// 数据目录
	dataDir := util_os.Path(DataDir, table)
	fileInfo, err := os.Stat(dataDir)
	// 数据目录不存在或者不是文件目录，则创建数据目录
	if (err != nil && !os.IsExist(err)) || !fileInfo.IsDir() {
		err = os.MkdirAll(dataDir, os.ModePerm)
	}
	if err != nil {
		db.Rollback()
		return redirect(err)
	}

	// 保存文件
	newFile, err := os.OpenFile(util_os.Path(dataDir, fmt.Sprintf("%d", id)),
		os.O_CREATE| // 创建文件，如果文件不存在的话
			os.O_WRONLY| // 只写
			os.O_TRUNC, // 清空文件，如果文件存在的话
		0666)
	if err != nil {
		db.Rollback()
		return redirect(err)
	}
	defer newFile.Close()
	_, err = newFile.Write(bytes)
	if err != nil {
		db.Rollback()
		return redirect(err)
	}

	// 提交事务
	db.Commit()

	return redirect(err)
}
