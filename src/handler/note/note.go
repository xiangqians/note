// @author xiangqian
// @date 21:34 2023/12/10
package note

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"note/src/db"
	"note/src/handler/common"
	"note/src/model"
	"note/src/session"
	util_filetype "note/src/util/filetype"
	util_os "note/src/util/os"
	util_string "note/src/util/string"
	util_time "note/src/util/time"
	util_validate "note/src/util/validate"
	"os"
	"strconv"
	"strings"
	"time"
)

// 内存值
var v int64 = 1

// 期望值
var a int64

// 新值
const b int64 = 0

func init() {
	// 启动一个 goroutine 来执行定时任务
	go func() {
		for {
			// 设定检测时间间隔
			time.Sleep(2 * time.Second)
			if v != b {
				// 设置期望值
				a = v

				startTime := time.Now()
				log.Println("开始执行定时任务【计算目录大小】...")
				err := calculateFolderSizeTask()
				endTime := time.Now()
				if err != nil {
					log.Println("执行定时任务【计算目录大小】发生错误", err)
				} else {
					duration := endTime.Sub(startTime)
					log.Println("执行定时任务【计算目录大小】完成，耗时：", duration)
				}

				// 如果内存值等于期望值，那么就将内存值更新为新值
				if v == a {
					v = b
				}
			}
		}
	}()
}

// TriggerCalculateFolderSizeTask 触发计算目录大小任务
func TriggerCalculateFolderSizeTask() {
	log.Println("==============触发计算目录大小任务================")
	v++
}

// 计算目录大小任务
func calculateFolderSizeTask() error {
	db := db.Get()
	result, err := db.Get("SELECT `id`, `del` FROM `note` WHERE `del` IN (0, 1) AND `type` = 'folder'")
	if err != nil {
		return err
	}

	var notes []model.Note
	err = result.Scan(&notes)
	if err != nil {
		return err
	}

	sql := "WITH RECURSIVE `tmp`(`id`, `pid`, `type`, `size`, `del`) AS (" +
		" SELECT t.`id`, t.`pid`, t.`type`, t.`size`, t.`del`" +
		" FROM `note` t" +
		" WHERE t.`pid` = %d" + // 起点条件
		" UNION ALL" +
		" SELECT t.`id`, t.`pid`, t.`type`, t.`size`, t.`del`" +
		" FROM `note` t" +
		" INNER JOIN `tmp` ON t.pid = `tmp`.id)" + // 关联递归查询结果
		" SELECT IFNULL(SUM(t.`size`), 0) FROM tmp t WHERE t.`del` IN (%s) AND t.`type` != 'folder'"

	for _, note := range notes {
		id := note.Id
		del := note.Del
		var dels = ""
		if del == 0 {
			dels = "0"
		} else {
			dels = "0, 1"
		}
		result, err = db.Get(fmt.Sprintf(sql, id, dels))
		if err != nil {
			return err
		}

		var size int64
		err = result.Scan(&size)
		if err != nil {
			return err
		}

		_, err = db.Upd("UPDATE `note` SET `size` = ? WHERE `del` = ? AND `id` = ?", size, del, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func AddFolder(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	return add(request, writer, session, util_filetype.Folder)
}

func AddMdFile(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	return add(request, writer, session, util_filetype.Md)
}

func Paste(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	var pid int64 = 0
	redirect := func(err any) (string, model.Response) {
		return fmt.Sprintf("redirect:/%s/%d/list", common.TableNote, pid), model.Response{Msg: util_string.String(err)}
	}

	// 文件id
	fromId, err := strconv.ParseInt(request.PostFormValue("fromId"), 10, 64)
	if err != nil || fromId <= 0 {
		return redirect(err)
	}

	// 目标文件夹id
	toId, err := strconv.ParseInt(request.PostFormValue("toId"), 10, 64)
	if err != nil || toId < 0 {
		return redirect(err)
	}
	pid = toId

	var result *db.Result
	db := db.Get()

	// 校验目标文件夹是否存在
	if toId > 0 {
		result, err = db.Get("SELECT `id`, `type` FROM `note` WHERE `del` = 0 AND `id` = ? LIMIT 1", toId)
		if err != nil {
			return redirect(err)
		}
		var note model.Note
		err = result.Scan(&note)
		if err != nil || note.Id == 0 || note.Type != util_filetype.Folder {
			return redirect(err)
		}
	}

	// 校验文件是否存在
	result, err = db.Get("SELECT `id` FROM `note` WHERE `del` = 0 AND `id` = ? LIMIT 1", fromId)
	if err != nil {
		return redirect(err)
	}
	fromId = 0
	err = result.Scan(&fromId)
	if err != nil || fromId == 0 {
		return redirect(err)
	}

	_, err = db.Upd("UPDATE `note` SET `pid` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ?", toId, util_time.NowUnix(), fromId)
	return redirect(err)
}

func Upd(request *http.Request, writer http.ResponseWriter, session *session.Session) (templateName string, response model.Response) {
	write := func(err any) {
		//writer.Header().Set("Content-Type", "application/json")
		writer.Write([]byte(util_string.String(err)))
	}

	vars := mux.Vars(request)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil || id <= 0 {
		write(err)
		return
	}

	db := db.Get()
	result, err := db.Get("SELECT `id`, `name`, `type`, `size`, `del`, `add_time`, `upd_time` FROM `note` WHERE `del` = 0 AND `id` = ? LIMIT 1", id)
	if err != nil {
		write(err)
		return
	}

	var note model.Note
	err = result.Scan(&note)
	if err != nil || note.Id == 0 || note.Type != util_filetype.Md {
		write(err)
		return
	}

	// 数据目录
	dataDir := util_os.Path(common.DataDir, common.TableNote)
	fileInfo, err := os.Stat(dataDir)
	// 数据目录不存在或者不是文件目录，则创建数据目录
	if (err != nil && !os.IsExist(err)) || !fileInfo.IsDir() {
		err = os.MkdirAll(dataDir, os.ModePerm)
	}
	if err != nil {
		write(err)
		return
	}

	// 打开文件
	file, err := os.OpenFile(util_os.Path(dataDir, fmt.Sprintf("%d", id)),
		os.O_CREATE| // 创建文件，如果文件不存在的话
			os.O_WRONLY| // 只写
			os.O_TRUNC, // 清空文件，如果文件存在的话
		0666)
	if err != nil {
		write(err)
		return
	}
	defer file.Close()

	content := request.PostFormValue("content")
	bytes := []byte(content)

	// 写入到文件
	_, err = file.Write(bytes)
	if err != nil {
		write(err)
		return
	}

	_, err = db.Upd("UPDATE `note` SET `size` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ?",
		len(bytes),
		util_time.NowUnix(),
		id)
	write(err)
	return
}

func add(request *http.Request, writer http.ResponseWriter, session *session.Session, Type string) (string, model.Response) {
	redirect := func(pid int64, err any) (string, model.Response) {
		return fmt.Sprintf("redirect:/note/%d/list", pid), model.Response{Msg: util_string.String(err)}
	}

	// pid
	pid, err := strconv.ParseInt(strings.TrimSpace(request.PostFormValue("pid")), 10, 64)
	if err != nil || pid < 0 {
		return redirect(pid, err)
	}

	// 名称
	name := strings.TrimSpace(request.PostFormValue("name"))
	if name == "" {
		return redirect(pid, nil)
	}
	err = util_validate.FileName(name, session.GetLanguage())
	if err != nil {
		return redirect(pid, err)
	}

	db := db.Get()

	// 校验父id是否存在
	if pid > 0 {
		result, err := db.Get("SELECT `id`, `type` FROM `note` WHERE `del` = 0 AND `id` = ? LIMIT 1", pid)
		if err != nil {
			return redirect(pid, err)
		}

		var note model.Note
		err = result.Scan(&note)
		if err != nil || note.Id == 0 || note.Type != util_filetype.Folder {
			return redirect(pid, err)
		}
	}

	// 获取永久删除id，以复用
	result, err := db.Get("SELECT `id` FROM `note` WHERE `del` = 2 LIMIT 1")
	if err != nil {
		return redirect(pid, err)
	}
	var id int64
	err = result.Scan(&id)
	if err != nil {
		return redirect(pid, err)
	}

	// 新id
	if id == 0 {
		_, id, err = db.Add("INSERT INTO `note` (`pid`, `name`, `type`, `add_time`) VALUES (?, ?, ?, ?)", pid, name, Type, util_time.NowUnix())
	} else
	// 复用id
	{
		_, err = db.Upd("UPDATE `note` SET `pid` = ?, `name` = ?, `type` = ?, `size` = 0, `del` = 0, `add_time` = ?, `upd_time` = 0 WHERE `id` = ?",
			pid,
			name,
			Type,
			util_time.NowUnix(),
			id)
	}

	return redirect(pid, err)
}
