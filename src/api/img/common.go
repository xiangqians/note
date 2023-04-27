// img common
// @author xiangqian
// @date 20:24 2023/04/27
package img

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	typ_api "note/src/typ/api"
	typ_page "note/src/typ/page"
	typ_resp "note/src/typ/resp"
	util_os "note/src/util/os"
	util_str "note/src/util/str"
)

func RedirectToList(context *gin.Context, err any) {
	resp := typ_resp.Resp[any]{
		Msg: util_str.TypeToStr(err),
	}

	// 记录查询参数
	img, err := common.GetSessionV[typ_api.Img](context, "img", false)
	if err != nil {
		common.Redirect(context, "/img/list", resp)
		return
	}

	common.Redirect(context, fmt.Sprintf("/img/list?id=%d&name=%s&type=%s&del=%d", img.Id, img.Name, img.Type, img.Del), resp)
}

// HistPath 获取图片历史记录物理路径
func HistPath(context *gin.Context, img typ_api.Img) (string, error) {
	// dir
	dataDir := common.DataDir(context)
	imgDir := fmt.Sprintf("%s%s%s%s%s%s%s", dataDir,
		util_os.FileSeparator(), "hist",
		util_os.FileSeparator(), "img",
		util_os.FileSeparator(), img.Type)
	if !util_os.IsExist(imgDir) {
		err := util_os.MkDir(imgDir)
		if err != nil {
			return "", err
		}
	}

	// file name
	time := img.UpdTime
	name := fmt.Sprintf("%d_%d", img.Id, time)

	// path
	return fmt.Sprintf("%s%s%s", imgDir, util_os.FileSeparator(), name), nil
}

// Path 获取图片物理路径
func Path(context *gin.Context, img typ_api.Img) (string, error) {
	// dir
	dataDir := common.DataDir(context)
	imgDir := fmt.Sprintf("%s%s%s%s%s", dataDir,
		util_os.FileSeparator(), "img",
		util_os.FileSeparator(), img.Type)
	if !util_os.IsExist(imgDir) {
		err := util_os.MkDir(imgDir)
		if err != nil {
			return "", err
		}
	}

	// file name
	name := fmt.Sprintf("%d", img.Id)

	// path
	return fmt.Sprintf("%s%s%s", imgDir, util_os.FileSeparator(), name), nil
}

// DbQryPermlyDelId 查询永久删除的图片记录id，以便复用此图片记录
func DbQryPermlyDelId(context *gin.Context) (int64, int64, error) {
	id, count, err := common.DbQry[int64](context, "SELECT `id` FROM `img` WHERE `del` = 2 LIMIT 1")
	return id, count, err
}

// DbPage 分页查询图片
func DbPage(context *gin.Context, img typ_api.Img) (typ_page.Page[typ_api.Img], error) {
	req, _ := common.PageReq(context)

	args := make([]any, 0, 1)
	sql := "SELECT i.`id`, i.`name`, i.`type`, i.`size`, i.`del`, i.`add_time`, i.`upd_time` FROM `img` i WHERE i.`del` = ? "
	args = append(args, img.Del)

	// id
	if img.Id > 0 {
		sql += "AND i.`id` = ? "
		args = append(args, img.Id)
	}

	// name
	if img.Name != "" {
		sql += "AND i.`name` LIKE '%' || ? || '%' "
		args = append(args, img.Name)
	}

	// type
	if img.Type != "" {
		sql += "AND i.`type` = ? "
		args = append(args, img.Type)
	}

	sql += "ORDER BY (CASE WHEN `upd_time` > `add_time` THEN `upd_time` ELSE `add_time` END) DESC"

	return common.DbPage[typ_api.Img](context, req, sql, args...)
}

// DbTypes 获取图片类型集合
func DbTypes(context *gin.Context) []string {
	// types
	types, count, err := common.DbQry[[]string](context, "SELECT DISTINCT(`type`) FROM `img` WHERE `del` = 0")
	if err != nil || count == 0 {
		types = nil
	}

	return types
}

func DbQry(context *gin.Context, id int64, del int) (typ_api.Img, int64, error) {
	img, count, err := common.DbQry[typ_api.Img](context, "SELECT `id`, `name`, `type`, `size`, `hist`, `hist_size`, `add_time`, `upd_time` FROM `img` WHERE `del` = ? AND `id` = ?", del, id)
	return img, count, err
}
