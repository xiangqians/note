// img
// @author xiangqian
// @date 11:34 2023/02/12
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
	util_time "note/src/util/time"
	"strings"
)

func Del(context *gin.Context) {
	redirect := func(err any) {
		resp := typ_resp.Resp[any]{Msg: util_str.TypeToStr(err)}
		common.Redirect(context, fmt.Sprintf("/img/list"), resp)
	}

	// id
	id, err := common.Param[int64](context, "id")
	if err != nil {
		redirect(err)
		return
	}

	// delete
	_, err = common.DbDel(context, "UPDATE `img` SET `del` = 1, `upd_time` = ? WHERE `id` = ?", util_time.NowUnix(), id)
	redirect(err)
	return
}

// UpdName 图片重命名
func UpdName(context *gin.Context) {
	redirect := func(err any) {
		resp := typ_resp.Resp[any]{Msg: util_str.TypeToStr(err)}
		common.Redirect(context, fmt.Sprintf("/img/list"), resp)
	}

	// img
	img := typ_api.Img{}
	err := common.ShouldBind(context, &img)
	if err != nil {
		redirect(err)
		return
	}

	// name
	img.Name = strings.TrimSpace(img.Name)
	//err = util_os.VerifyFileName(img.Name)
	//if err != nil {
	//	redirect(err)
	//	return
	//}

	// update
	_, err = common.DbUpd(context, "UPDATE `img` SET `name` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `name` <> ?", img.Name, util_time.NowUnix(), img.Id, img.Name)

	redirect(err)
	return
}

// List 图片列表页面
func List(context *gin.Context) {
	page, err := DbPage(context, 0)
	resp := typ_resp.Resp[typ_page.Page[typ_api.Img]]{
		Msg:  util_str.TypeToStr(err),
		Data: page,
	}
	common.HtmlOk(context, "img/list.html", resp)
}

// HistPath 获取图片历史记录物理路径
func HistPath(context *gin.Context, img typ_api.Img) (string, error) {
	// dir
	dataDir := common.DataDir(context)
	imgDir := fmt.Sprintf("%s%s%s%s%s%s%s", dataDir,
		util_os.FileSeparator, "hist",
		util_os.FileSeparator, "img",
		util_os.FileSeparator, img.Type)
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
	return fmt.Sprintf("%s%s%s", imgDir, util_os.FileSeparator, name), nil
}

// Path 获取图片物理路径
func Path(context *gin.Context, img typ_api.Img) (string, error) {
	// dir
	dataDir := common.DataDir(context)
	imgDir := fmt.Sprintf("%s%s%s%s%s", dataDir,
		util_os.FileSeparator, "img",
		util_os.FileSeparator, img.Type)
	if !util_os.IsExist(imgDir) {
		err := util_os.MkDir(imgDir)
		if err != nil {
			return "", err
		}
	}

	// file name
	name := fmt.Sprintf("%d", img.Id)

	// path
	return fmt.Sprintf("%s%s%s", imgDir, util_os.FileSeparator, name), nil
}

func DbPage(context *gin.Context, del int) (typ_page.Page[typ_api.Img], error) {
	req, _ := common.PageReq(context)
	return common.DbPage[typ_api.Img](context, req, "SELECT `id`, `name`, `type`, `size`, `add_time`, `upd_time` FROM `img` WHERE `del` = ? ORDER BY (CASE WHEN `upd_time` > `add_time` THEN `upd_time` ELSE `add_time` END) DESC", del)
}

// DbQryPermlyDelId 查询永久删除的图片记录id，以便复用此图片记录
func DbQryPermlyDelId(context *gin.Context) (int64, int64, error) {
	id, count, err := common.DbQry[int64](context, "SELECT `id` FROM `img` WHERE `del` = 2 LIMIT 1")
	return id, count, err
}

func DbQry(context *gin.Context, id int64, del int) (typ_api.Img, int64, error) {
	img, count, err := common.DbQry[typ_api.Img](context, "SELECT `id`, `name`, `type`, `size`, `hist`, `hist_size`, `add_time`, `upd_time` FROM `img` WHERE `del` = ? AND `id` = ?", del, id)
	return img, count, err
}
