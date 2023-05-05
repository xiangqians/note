// img common
// @author xiangqian
// @date 20:24 2023/04/27
package img

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	typ_api "note/src/typ/api"
	typ_resp "note/src/typ/resp"
	util_json "note/src/util/json"
	util_os "note/src/util/os"
	util_str "note/src/util/str"
	"sort"
)

// DeserializeHist 反序列化历史记录
func DeserializeHist(hist string) ([]typ_api.Img, error) {
	if hist == "" {
		return nil, nil
	}

	// hists
	hists := make([]typ_api.Img, 0, 1) // len 0, cap ?
	err := util_json.Deserialize(hist, &hists)
	if err != nil {
		return nil, err
	}

	// sort
	Sort(&hists)

	return hists, nil
}

// SerializeHist 序列化历史记录
func SerializeHist(hists []typ_api.Img) (string, error) {
	return util_json.Serialize(hists)
}

// Sort 对img进行排序
func Sort(imgs *[]typ_api.Img) {
	sort.Slice(*imgs, func(i, j int) bool {
		return (*imgs)[i].UpdTime > (*imgs)[j].UpdTime
	})
}

func RedirectToList(context *gin.Context, err any) {
	resp := typ_resp.Resp[any]{
		Msg: util_str.ConvTypeToStr(err),
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

// DbQry 查询图片信息
func DbQry(context *gin.Context, id int64, del int) (typ_api.Img, int64, error) {
	img, count, err := common.DbQry[typ_api.Img](context, "SELECT `id`, `name`, `type`, `size`, `hist`, `hist_size`, `del`, `add_time`, `upd_time` FROM `img` WHERE `del` = ? AND `id` = ?", del, id)
	return img, count, err
}
