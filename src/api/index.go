// index
// @author xiangqian
// @date 17:21 2023/02/04
package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"note/src/typ"
)

func IndexPage(pContext *gin.Context) {
	//pid, err := Param[int64](pContext, "pid")
	pid := 0
	log.Printf("pid = %d\n", pid)

	// 查询目录
	dfs, count, err := DbQry[[]typ.DF](pContext, "SELECT d.* FROM `dir` d WHERE d.`del` = 0 AND d.`pid` = ?", pid)
	if count == 0 {
		dfs = nil
	}

	tmpDfs, count, err := DbQry[[]typ.DF](pContext, "SELECT f.* FROM file f JOIN dir_file df ON df.file_id = f.id WHERE f.del = 0 AND df.dir_id = ?", pid)
	if count == 0 {
		dfs = append(dfs, tmpDfs...)
	}

	Html(pContext, "index.html", gin.H{"dfs": dfs}, err)
	return
}
