// @author xiangqian
// @date 19:40 2023/12/04
package common

import (
	"fmt"
	"io"
	"net/http"
	"note/src/db"
	"note/src/model"
	"note/src/session"
	util_filetype "note/src/util/filetype"
	util_i18n "note/src/util/i18n"
	util_os "note/src/util/os"
	"note/src/util/time"
	"os"
	"strings"
)

func Upload(request *http.Request, writer http.ResponseWriter, session *session.Session, table string) (string, model.Response) {
	// 读取上传文件
	file, fileHeader, err := request.FormFile("file")
	if err != nil {
		return redirect(table, err)
	}
	defer file.Close()

	// 读取文件字节内容
	bytes, err := io.ReadAll(file)
	if err != nil {
		return redirect(table, err)
	}

	// 获取文件类型
	filetype := util_filetype.GetType(bytes)
	if table == "image" {
		if filetype != util_filetype.Ico &&
			filetype != util_filetype.Gif &&
			filetype != util_filetype.Jpg &&
			filetype != util_filetype.Jpeg &&
			filetype != util_filetype.Png &&
			filetype != util_filetype.Webp {
			return redirect(table, fmt.Sprintf(util_i18n.GetMessage("i18n.fileTypeUnsupportedUpload", session.GetLanguage()), filetype))
		}
	}

	// 文件名
	name := fileHeader.Filename
	// 去除文件后缀名
	suffix := "." + filetype
	if strings.HasSuffix(name, suffix) {
		name = name[:len(name)-len(suffix)]
	}

	// 文件大小，单位：字节
	size := fileHeader.Size

	// 获取永久删除id，以复用
	id, err := getPermlyDelId(table)
	if err != nil {
		return redirect(table, err)
	}

	db := db.Get()

	// 开启事务
	err = db.Begin()
	if err != nil {
		return redirect(table, err)
	}

	// 入库
	// 新id
	if id == 0 {
		_, id, err = db.Add(fmt.Sprintf("INSERT INTO `%s` (`name`, `type`, `size`, `add_time`) VALUES (?, ?, ?, ?)", table), name, filetype, size, time.NowUnix())
	} else
	// 复用id
	{
		_, err = db.Upd(fmt.Sprintf("UPDATE `%s` SET `name` = ?, `type` = ?, `size` = ?, `hist` = '', `hist_size` = 0, `del` = 0, `add_time` = ?, `upd_time` = 0 WHERE `id` = ?", table), name, filetype, size, time.NowUnix(), id)
	}
	if err != nil {
		db.Rollback()
		return redirect(table, err)
	}

	// 数据目录
	dataDir := util_os.Path(dataDir, table)
	fileInfo, err := os.Stat(dataDir)
	// 数据目录不存在或者不是文件目录，则创建数据目录
	if (err != nil && !os.IsExist(err)) || !fileInfo.IsDir() {
		err = os.MkdirAll(dataDir, os.ModePerm)
	}
	if err != nil {
		db.Rollback()
		return redirect(table, err)
	}

	// 保存文件
	newFile, err := os.OpenFile(util_os.Path(dataDir, fmt.Sprintf("%d", id)),
		os.O_CREATE| // 创建文件，如果文件不存在的话
			os.O_WRONLY| // 只写
			os.O_TRUNC, // 清空文件，如果文件存在的话
		0666)
	if err != nil {
		db.Rollback()
		return redirect(table, err)
	}
	defer newFile.Close()
	_, err = newFile.Write(bytes)
	if err != nil {
		db.Rollback()
		return redirect(table, err)
	}

	// 提交事务
	db.Commit()

	return redirect(table, err)
}
