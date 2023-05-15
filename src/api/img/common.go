// img common
// @author xiangqian
// @date 20:24 2023/04/27
package img

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	api_common_context "note/src/api/common/context"
	"note/src/api/common/db"
	"note/src/api/common/session"
	"note/src/typ"
	"note/src/util/json"
	util_os "note/src/util/os"
	"note/src/util/str"
	"os"
	"sort"
)

const ImgSessionKey = "img"

// DeserializeHist 反序列化历史记录
func DeserializeHist(hist string) ([]typ.Img, error) {
	if hist == "" {
		return nil, nil
	}

	// hists
	hists := make([]typ.Img, 0, 1) // len 0, cap ?
	err := json.Deserialize(hist, &hists)
	if err != nil {
		return nil, err
	}

	// sort
	Sort(&hists)

	return hists, nil
}

// SerializeHist 序列化历史记录
func SerializeHist(hists []typ.Img) (string, error) {
	return json.Serialize(hists)
}

// Sort 对img进行排序
func Sort(imgs *[]typ.Img) {
	sort.Slice(*imgs, func(i, j int) bool {
		return (*imgs)[i].UpdTime > (*imgs)[j].UpdTime
	})
}

func RedirectToList(context *gin.Context, err any) {
	resp := typ.Resp[any]{
		Msg: str.ConvTypeToStr(err),
	}

	// 记录查询参数
	img, err := session.Get[typ.Img](context, ImgSessionKey, false)
	if err != nil {
		api_common_context.Redirect(context, "/img/list", resp)
		return
	}

	api_common_context.Redirect(context, fmt.Sprintf("/img/list?id=%d&name=%s&type=%s&del=%d", img.Id, img.Name, img.Type, img.Del), resp)
}

// DbQry 查询图片信息
func DbQry(context *gin.Context, id int64, del int) (typ.Img, int64, error) {
	img, count, err := db.Qry[typ.Img](context, "SELECT `id`, `name`, `type`, `size`, `hist`, `hist_size`, `del`, `add_time`, `upd_time` FROM `img` WHERE `del` = ? AND `id` = ?", del, id)
	return img, count, err
}

// HistPath 获取图片历史记录物理路径
func HistPath(context *gin.Context, img typ.Img) (string, error) {
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

// DelImg 删除图片
func DelImg(context *gin.Context, img typ.Img) (string, error) {
	// path
	path, err := Path(context, img)
	if err != nil {
		return path, err
	}

	// del
	return path, util_os.DelFile(path)
}

// ClearImg 清空图片
func ClearImg(context *gin.Context, img typ.Img) (string, error) {
	// path
	path, err := Path(context, img)
	if err != nil {
		return path, err
	}

	// exist ?
	if !util_os.IsExist(path) {
		return path, nil
	}

	// open
	file, err := os.OpenFile(path,
		os.O_WRONLY|os.O_TRUNC, // 只写（O_WRONLY） & 清空文件（O_TRUNC）
		0666)

	// close
	file.Close()

	return path, err
}

// Path 获取图片物理路径
func Path(context *gin.Context, img typ.Img) (string, error) {
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
