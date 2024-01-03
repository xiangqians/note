// @author xiangqian
// @date 13:41 2023/11/19
package index

import (
	"net/http"
	"note/src/db"
	"note/src/model"
	"note/src/session"
	util_string "note/src/util/string"
)

func Index(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	var statsArr []model.Stats
	html := func(err any) (string, model.Response) {
		return "index", model.Response{
			Msg:  util_string.String(err),
			Data: statsArr,
		}
	}

	db := db.Get()
	result, err := db.Get("SELECT 'note' AS 'type', COUNT(`id`) AS 'count', SUM(`size`) AS 'size' FROM `note` WHERE `del` = 0 AND `type` != 'folder'" +
		" UNION ALL SELECT 'image' AS 'type', COUNT(`id`) AS 'count', SUM(`size`) AS 'size' FROM `image` WHERE `del` = 0" +
		" UNION ALL SELECT 'audio' AS 'type', COUNT(`id`) AS 'count', SUM(`size`) AS 'size' FROM `audio` WHERE `del` = 0" +
		" UNION ALL SELECT 'video' AS 'type', COUNT(`id`) AS 'count', SUM(`size`) AS 'size' FROM `video` WHERE `del` = 0")
	if err != nil {
		return html(err)
	}

	err = result.Scan(&statsArr)
	return html(err)
}

/*
SELECT `type`,
COUNT(`id`) AS 'count',
SUM(`size`) AS 'size'
FROM `note`
WHERE `del` = 0
GROUP BY `type`
*/
