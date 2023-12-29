// @author xiangqian
// @date 22:39 2023/12/04
package common

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"note/src/db"
	"note/src/model"
	"note/src/session"
	"note/src/util/filetype"
	util_os "note/src/util/os"
	util_string "note/src/util/string"
	"os"
	"strconv"
	"strings"
)

func View(request *http.Request, writer http.ResponseWriter, session *session.Session, table string) (string, model.Response) {
	vars := mux.Vars(request)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil || id <= 0 {
		return NotFound(request, writer, session, err)
	}

	switch table {
	case TableImage, TableAudio, TableVideo:
		return absView(request, writer, session, table, id)

	case TableNote:
		return noteView(request, writer, session, table, id)
	}

	return NotFound(request, writer, session, err)
}

func noteView(request *http.Request, writer http.ResponseWriter, session *session.Session, table string, id int64) (string, model.Response) {
	db := db.Get()
	result, err := db.Get(fmt.Sprintf("SELECT `id`, `pid`, `name`, `type`, `size`, `del`, `add_time`, `upd_time` FROM `%s` WHERE `del` = 0 AND `id` = ? LIMIT 1", table), id)
	if err != nil {
		return NotFound(request, writer, session, err)
	}

	var note model.Note
	err = result.Scan(&note)
	if err != nil || note.Id == 0 || note.Type == filetype.Folder {
		return NotFound(request, writer, session, err)
	}

	pid := note.Pid
	var pNote model.PNote
	if pid > 0 {
		result, err = db.Get(getPNoteSql(), pid)
		if err != nil {
			return NotFound(request, writer, session, err)
		}

		err = result.Scan(&pNote)
		if err != nil {
			return NotFound(request, writer, session, err)
		}

		if pNote.IdsStr != "" {
			pNote.Ids = strings.Split(pNote.IdsStr, "/")[1:]
			pNote.Names = strings.Split(pNote.NamesStr, "/")[1:]
		}
	}
	pNote.Id = pid

	var templateName string
	switch note.Type {
	case filetype.Md:
		templateName = fmt.Sprintf("%s/md/view", table)
		file, err := os.Open(util_os.Path(DataDir, table, fmt.Sprintf("%d", id)))
		if err == nil {
			defer file.Close()
			bytes, err := io.ReadAll(file)
			if err == nil {
				note.Content = string(bytes)
			}
		}

	case filetype.Pdf:
		version := strings.TrimSpace(request.URL.Query().Get("version"))
		if version != "v1" && version != "v2" {
			version = "v1"
		}
		templateName = fmt.Sprintf("%s/pdf/view/%s", table, version)

	default:
		templateName = fmt.Sprintf("%s/default/view", table)
	}

	return templateName, model.Response{Msg: session.GetMsg() + util_string.String(err), Data: map[string]any{
		"table": table,
		"data":  note,
		"pNote": pNote,
	}}
}

func absView(request *http.Request, writer http.ResponseWriter, session *session.Session, table string, id int64) (string, model.Response) {
	db := db.Get()
	result, err := db.Get(fmt.Sprintf("SELECT `id`, `name`, `type`, `size`, `del`, `add_time`, `upd_time` FROM `%s` WHERE `del` = 0 AND `id` = ? LIMIT 1", table), id)
	if err != nil {
		return NotFound(request, writer, session, err)
	}

	var abs model.Abs
	err = result.Scan(&abs)
	if err != nil || abs.Id == 0 {
		return NotFound(request, writer, session, err)
	}

	return fmt.Sprintf("%s/view", table),
		model.Response{Msg: session.GetMsg() + util_string.String(err), Data: map[string]any{
			"table": table,
			"data":  abs,
		}}
}
