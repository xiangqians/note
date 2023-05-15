// get
// @author xiangqian
// @date 22:56 2023/05/15
package note

import (
	"github.com/gin-gonic/gin"
	"log"
	api_common_context "note/src/api/common/context"
	"note/src/typ"
)

// GetHist 获取历史笔记
func GetHist(context *gin.Context) {
	// id
	id, err := api_common_context.Param[int64](context, "id")
	if err != nil {
		log.Println(err)
		return
	}

	// idx
	idx, err := api_common_context.Param[int](context, "idx")
	if err != nil || idx < 0 {
		log.Println(err)
		return
	}

	// note
	note, count, err := DbQry(context, id, 0, 0)
	if err != nil || count == 0 || typ.FtD == typ.ExtNameOf(note.Type) { // 排除目录
		log.Println(err)
		return
	}

	// hists
	histNotes, err := DeserializeHist(note.Hist)
	if err != nil || histNotes == nil {
		log.Println("hist is empty")
		return
	}

	// 校验idx是否合法
	if idx >= len(histNotes) {
		log.Println(err)
		return
	}

	// hist note
	histNote := histNotes[idx]

	// path
	path, err := HistPath(context, histNote)
	if err != nil {
		log.Println(err)
		return
	}

	// write
	context.File(path)
}

func Get(context *gin.Context) {
	// id
	id, err := api_common_context.Param[int64](context, "id")
	if err != nil {
		log.Println(err)
		return
	}

	// note
	note, count, err := DbQry(context, id, 0, 0)
	if err != nil || count == 0 || typ.FtD == typ.ExtNameOf(note.Type) { // 排除目录
		log.Println(err)
		return
	}

	// path
	path, err := Path(context, note)
	if err != nil {
		log.Println(err)
		return
	}

	/**
	// read all
	buf, err := os.ReadFile(fPath)
	if err != nil {
		log.Println(err)
		return
	}
	writer := context.Writer
	writer.Write(buf)
	writer.Flush()
	*/

	/**
	// open
	pFile, err := os.Open(fPath)
	if err != nil {
		log.Println(err)
		return
	}

	// write
	err = util.IOCopy(pFile, context.Writer, 0)
	if err != nil {
		log.Println(err)
		return
	}
	*/

	// write
	context.File(path)
}
