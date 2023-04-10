// img upload
// @author xiangqian
// @date 21:26 2023/04/10
package img

import (
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"net/http"
	"note/src/api/common"
	typ_api "note/src/typ/api"
	typ_ft "note/src/typ/ft"
	typ_resp "note/src/typ/resp"
	util_json "note/src/util/json"
	util_os "note/src/util/os"
	util_str "note/src/util/str"
	util_time "note/src/util/time"
	"os"
	"strings"
)

// Upload 图片上传
func Upload(context *gin.Context) {
	// method
	method := context.Request.Method

	// 响应数据类型
	dataType, _ := common.PostForm[string](context, "dataType")

	// redirect
	redirect := func(id int64, err any) {
		resp := typ_resp.Resp[int64]{Msg: util_str.TypeToStr(err), Data: id}

		// 响应json格式
		if dataType == "json" {
			common.JsonOk(context, resp)
			return
		}

		switch method {
		// 上传
		case http.MethodPost:
			common.Redirect(context, fmt.Sprintf("/img/list"), resp)

		// 重传
		case http.MethodPut:
			common.Redirect(context, fmt.Sprintf("/img/%v/view", id), resp)
		}
	}

	// id
	var id int64
	var err error

	// 重传文件必须有id
	if method == http.MethodPut {
		id, err = common.PostForm[int64](context, "id")
		if err != nil || id <= 0 {
			redirect(id, err)
			return
		}
	}

	// file header
	fh, err := context.FormFile("file")
	if err != nil || fh == nil {
		redirect(id, err)
		return
	}

	// file name
	fn := strings.TrimSpace(fh.Filename)

	// file type
	contentType := fh.Header.Get("Content-Type")
	ft := typ_ft.ContentTypeOf(contentType)
	if !typ_ft.IsImg(ft) {
		redirect(id, fmt.Sprintf("%s, %s", i18n.MustGetMessage("i18n.fileTypeUnsupported"), contentType))
		return
	}

	// size
	fs := fh.Size

	// 原图片信息
	var oldImg typ_api.Img

	// 操作数据库
	switch method {
	case http.MethodPost:
		// 查询是否有永久删除的图片记录id，以便复用此图片记录
		var count int64
		id, count, err = DbQryPermlyDelId(context)
		// 新id
		if err != nil || count == 0 {
			id, err = common.DbAdd(context, "INSERT INTO `img` (`name`, `type`, `size`, `add_time`) VALUES (?, ?, ?, ?)", fn, ft, fs, util_time.NowUnix())
		} else
		// 复用id
		{
			_, err = common.DbUpd(context, "UPDATE `img` SET `name` = ?, `type` = ?, `size` = ?, `del` = 0, `add_time` = ? WHERE `id` = ?", fn, ft, fs, util_time.NowUnix(), id)
		}

	case http.MethodPut:
		// 原图片信息
		var count int64
		oldImg, count, err = DbQry(context, id, 0)
		if err != nil || count == 0 {
			redirect(id, err)
			return
		}

		// 图片历史记录
		hist := oldImg.Hist
		histSize := oldImg.HistSize
		histImgs := make([]typ_api.Img, 0, 1) // len 0, cap ?
		if hist != "" {
			err = util_json.Deserialize(hist, &histImgs)
			if err != nil {
				redirect(id, err)
				return
			}
		}

		// 将原图片添加到历史记录
		histImg := typ_api.Img{
			Abs: typ_api.Abs{
				Id:      oldImg.Id,
				AddTime: oldImg.AddTime,
				UpdTime: oldImg.UpdTime,
			},
			Name: oldImg.Name,
			Type: oldImg.Type,
			Size: oldImg.Size,
		}
		histImgs = append(histImgs, histImg)
		// src
		var srcPath string
		srcPath, err = Path(context, histImg)
		if err != nil {
			redirect(id, err)
			return
		}
		// dst
		var dstPath string
		dstPath, err = HistPath(context, histImg)
		if err != nil {
			redirect(id, err)
			return
		}
		// copy
		err = util_os.CopyFile(srcPath, dstPath)
		if err != nil {
			redirect(id, err)
			return
		}

		// 图片历史记录至多保存15张，超过15张则删除最早地历史图片
		max := 15
		if len(histImgs) > max {
			l := len(histImgs) - max
			for i := 0; i < l; i++ {
				path, err := HistPath(context, histImgs[i])
				if err == nil {
					util_os.DelFile(path)
				}
			}
			histImgs = histImgs[l:]
		}

		// hist size
		for _, imgHist := range histImgs {
			histSize += imgHist.Size
		}

		// serialize
		hist, err = util_json.Serialize(histImgs)
		if err != nil {
			redirect(id, err)
			return
		}

		// update
		_, err = common.DbUpd(context, "UPDATE `img` SET `name` = ?, `type` = ?, `size` = ?, `hist` = ?, `hist_size` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ?",
			fn, ft, fs, hist, histSize, util_time.NowUnix(), id)
	}
	if err != nil {
		redirect(id, err)
		return
	}

	// path
	img := typ_api.Img{}
	img.Id = id
	img.Type = string(ft)
	path, err := Path(context, img)
	if err != nil {
		redirect(id, err)
		return
	}

	// 清空文件
	if util_os.IsExist(path) {
		var file *os.File
		file, err = os.OpenFile(path,
			os.O_WRONLY|os.O_TRUNC, // 只写（O_WRONLY） & 清空文件（O_TRUNC）
			0666)
		if err != nil {
			redirect(id, err)
			return
		}
		file.Close()
	}

	// 保存文件
	err = context.SaveUploadedFile(fh, path)

	// 保存文件成功时，判断如果重传不是同一个文件类型，则删除之前文件
	if err == nil && method == http.MethodPut && oldImg.Type != img.Type {
		var oldPath string
		oldPath, err = Path(context, oldImg)
		if err == nil {
			util_os.DelFile(oldPath)
		}
	}

	// redirect
	redirect(id, err)
	return
}
