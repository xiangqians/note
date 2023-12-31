// @author xiangqian
// @date 19:40 2023/12/04
package common

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"note/src/db"
	"note/src/model"
	"note/src/session"
	util_filetype "note/src/util/filetype"
	util_i18n "note/src/util/i18n"
	util_json "note/src/util/json"
	util_os "note/src/util/os"
	util_string "note/src/util/string"
	"note/src/util/time"
	"os"
	"strconv"
	"strings"
)

func Upload(request *http.Request, writer http.ResponseWriter, session *session.Session, table string) (string, model.Response) {
	var id int64
	var name string

	// note.pid
	pid, _ := strconv.ParseInt(request.PostFormValue("pid"), 10, 64)
	if pid < 0 {
		pid = 0
	}

	// 重定向函数
	redirect := func(err any) (string, model.Response) {
		if table == TableNote {
			return fmt.Sprintf("redirect:/%s/%d/list", table, pid), model.Response{Msg: util_string.String(err)}
		}

		if strings.TrimSpace(request.PostFormValue("dataType")) == "json" {
			response := model.Response{Msg: util_string.String(err)}
			response.Data = map[string]any{
				"id":   id,
				"name": name,
			}
			json, _ := util_json.Serialize(response, false)
			writer.Write([]byte(json))
			return "", response
		}

		return fmt.Sprintf("redirect:/%s", table), model.Response{Msg: util_string.String(err)}
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
	name = strings.TrimSpace(fileHeader.Filename)

	// 获取文件类型
	filetype := util_filetype.GetType(name)
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

	// 获取永久删除id，以复用
	result, err = db.Get(fmt.Sprintf("SELECT `id` FROM `%s` WHERE `del` = 2 LIMIT 1", table))
	if err != nil {
		db.Rollback()
		return redirect(err)
	}
	//var id int64
	err = result.Scan(&id)
	if err != nil {
		db.Rollback()
		return redirect(err)
	}

	// 校验note.pid
	if table == TableNote && pid > 0 {
		result, err = db.Get("SELECT `id`, `type` FROM `note` WHERE `del` = 0 AND `id` = ? LIMIT 1", pid)
		if err != nil {
			db.Rollback()
			return redirect(err)
		}
		var note model.Note
		err = result.Scan(&note)
		if err != nil || note.Id == 0 || note.Type != util_filetype.Folder {
			db.Rollback()
			return redirect(err)
		}
	}

	// 复用id
	if id > 0 {
		// len 0, cap ?
		capacity := 6
		columns := make([]string, 0, capacity)
		values := make([]any, 0, capacity)

		if table == TableNote {
			columns = append(columns, "`pid` = ?")
		}
		columns = append(columns, "`name` = ?")
		columns = append(columns, "`type` = ?")
		columns = append(columns, "`size` = ?")
		columns = append(columns, "`del` = 0")
		columns = append(columns, "`add_time` = ?")
		columns = append(columns, "`upd_time` = 0")

		if table == TableNote {
			values = append(values, pid)
		}
		values = append(values, name)
		values = append(values, filetype)
		values = append(values, size)
		values = append(values, time.NowUnix())
		values = append(values, id)

		_, err = db.Upd(fmt.Sprintf("UPDATE `%s` SET %s  WHERE `id` = ?", table, strings.Join(columns, ", ")), values...)
	} else
	// 新id
	{
		// len 0, cap ?
		capacity := 5
		columns := make([]string, 0, capacity)
		values := make([]any, 0, capacity)

		if table == TableNote {
			columns = append(columns, "`pid`")
		}
		columns = append(columns, "`name`")
		columns = append(columns, "`type`")
		columns = append(columns, "`size`")
		columns = append(columns, "`add_time`")

		if table == TableNote {
			values = append(values, pid)
		}
		values = append(values, name)
		values = append(values, filetype)
		values = append(values, size)
		values = append(values, time.NowUnix())

		placeholders := make([]string, 0, capacity)
		for i, length := 0, len(columns); i < length; i++ {
			placeholders = append(placeholders, "?")
		}

		_, id, err = db.Add(fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", table, strings.Join(columns, ", "), strings.Join(placeholders, ", ")), values...)
	}

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
	err = db.Commit()

	return redirect(err)
}

func validateFiletype(session *session.Session, table, filetype string) error {
	switch table {
	case TableImage:
		if filetype != util_filetype.Ico &&
			filetype != util_filetype.Gif &&
			filetype != util_filetype.Jpg &&
			filetype != util_filetype.Jpeg &&
			filetype != util_filetype.Png &&
			filetype != util_filetype.Webp {
			return errors.New(fmt.Sprintf(util_i18n.GetMessage("i18n.fileTypeUnsupportedUpload", session.GetLanguage()), filetype))
		}

	case TableAudio:
		if filetype != util_filetype.Mp3 &&
			filetype != util_filetype.Wav &&
			filetype != util_filetype.Flac &&
			filetype != util_filetype.Aac &&
			filetype != util_filetype.Ogg {
			return errors.New(fmt.Sprintf(util_i18n.GetMessage("i18n.fileTypeUnsupportedUpload", session.GetLanguage()), filetype))
		}

	case TableVideo:
		if filetype != util_filetype.Mp4 &&
			filetype != util_filetype.Avi &&
			filetype != util_filetype.Mov &&
			filetype != util_filetype.Wmv &&
			filetype != util_filetype.Mkv &&
			filetype != util_filetype.Flv {
			return errors.New(fmt.Sprintf(util_i18n.GetMessage("i18n.fileTypeUnsupportedUpload", session.GetLanguage()), filetype))
		}

	case TableNote:
		if filetype != util_filetype.Doc &&
			filetype != util_filetype.Docx &&
			filetype != util_filetype.Pdf &&
			filetype != util_filetype.Zip &&
			filetype != util_filetype.TarGz {
			return errors.New(fmt.Sprintf(util_i18n.GetMessage("i18n.fileTypeUnsupportedUpload", session.GetLanguage()), filetype))
		}
	}

	return nil
}
