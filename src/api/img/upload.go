// img upload
// @author xiangqian
// @date 21:26 2023/04/10
package img

import (
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	api_common_context "note/src/api/common/context"
	"note/src/api/common/db"
	"note/src/typ"
	"note/src/util/os"
	"note/src/util/str"
	"note/src/util/time"
	"note/src/util/validate"
	"strings"
)

// ReUpload 重新上传图片
func ReUpload(context *gin.Context) {
	// redirect
	redirect := func(id int64, err any) {
		resp := typ.Resp[any]{Msg: str.ConvTypeToStr(err)}
		api_common_context.Redirect(context, fmt.Sprintf("/img/%d/view", id), resp)
	}

	// id
	id, err := api_common_context.PostForm[int64](context, "id")
	if err != nil || id <= 0 {
		redirect(id, err)
		return
	}

	// file header
	fh, err := context.FormFile("file")
	if err != nil || fh == nil {
		redirect(id, err)
		return
	}

	// name
	name := strings.TrimSpace(fh.Filename)
	// validate name
	err = validate.FileName(name)
	if err != nil {
		redirect(id, err)
		return
	}

	// type
	contentType := fh.Header.Get("Content-Type")
	ft := typ.ContentTypeOf(contentType)
	if !typ.IsImg(ft) {
		redirect(id, fmt.Sprintf("%s, %s", i18n.MustGetMessage("i18n.fileTypeUnsupportedUpload"), contentType))
		return
	}
	_type := string(ft)

	// file size
	size := fh.Size

	// img
	var count int64
	img, count, err := DbQry(context, id, 0)
	if err != nil || count == 0 {
		redirect(id, err)
		return
	}

	// 图片历史记录
	histImgs, err := DeserializeHist(img.Hist)
	if err != nil {
		redirect(id, err)
		return
	}
	if histImgs == nil {
		histImgs = make([]typ.Img, 0, 1)
	}

	// 将原图片添加到历史记录
	histImg := typ.Img{
		Abs: typ.Abs{
			Id:      img.Id,
			AddTime: img.AddTime,
			UpdTime: img.UpdTime,
		},
		Name: img.Name,
		Type: img.Type,
		Size: img.Size,
	}
	histImgs = append(histImgs, histImg)
	Sort(&histImgs)

	// 备份最近一条历史记录
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
	_, err = os.CopyFile(dstPath, srcPath)
	if err != nil {
		redirect(id, err)
		return
	}

	// 图片历史记录至多保存15张，超过15张则删除最早地历史图片
	max := 15
	l := len(histImgs)
	if l > max {
		for i := max; i < l; i++ {
			path, err := HistPath(context, histImgs[i])
			if err == nil {
				os.DelFile(path)
			}
		}
		histImgs = histImgs[:max]
	}

	// hist size
	var histSize int64 = 0
	for _, imgHist := range histImgs {
		histSize += imgHist.Size
	}

	// serialize
	hist, err := SerializeHist(histImgs)
	if err != nil {
		redirect(id, err)
		return
	}

	// new img
	newImg := typ.Img{
		Abs: typ.Abs{
			Id:      id,
			UpdTime: time.NowUnix(),
		},
		Name:     name,
		Type:     _type,
		Size:     size,
		Hist:     hist,
		HistSize: histSize,
	}

	// update
	_, err = db.Upd(context, "UPDATE `img` SET `name` = ?, `type` = ?, `size` = ?, `hist` = ?, `hist_size` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ?",
		newImg.Name, newImg.Type, newImg.Size, newImg.Hist, newImg.HistSize, newImg.UpdTime, newImg.Id)
	if err != nil {
		redirect(id, err)
		return
	}

	// 清空文件
	path, err := ClearImg(context, newImg)

	// 保存文件
	err = context.SaveUploadedFile(fh, path)

	// 保存文件成功时，判断如果重传不是同一个文件类型，则删除之前文件
	if img.Type != newImg.Type {
		path, err = Path(context, img)
		if err == nil {
			os.DelFile(path)
		}
	}

	// redirect
	redirect(id, err)
}

// Upload 上传图片
func Upload(context *gin.Context) {
	// redirect
	redirect := func(img typ.Img, err any) {
		resp := typ.Resp[typ.Img]{Msg: str.ConvTypeToStr(err), Data: img}
		dataType, _ := api_common_context.PostForm[string](context, "dataType")
		if dataType == "json" {
			if err != nil {
				api_common_context.JsonBadRequest(context, resp)
			} else {
				api_common_context.JsonOk(context, resp)
			}
		} else {
			api_common_context.Redirect(context, fmt.Sprintf("/img/list"), resp)
		}
	}

	// file header
	fh, err := context.FormFile("file")
	if err != nil || fh == nil {
		redirect(typ.Img{}, err)
		return
	}

	// name
	name := strings.TrimSpace(fh.Filename)
	// validate name
	err = validate.FileName(name)
	if err != nil {
		redirect(typ.Img{}, err)
		return
	}

	// type
	contentType := fh.Header.Get("Content-Type")
	ft := typ.ContentTypeOf(contentType)
	if !typ.IsImg(ft) {
		redirect(typ.Img{}, fmt.Sprintf("%s, %s", i18n.MustGetMessage("i18n.fileTypeUnsupportedUpload"), contentType))
		return
	}
	_type := string(ft)

	// file size
	size := fh.Size

	// 查询是否有永久删除的图片记录id，以复用
	id, count, err := DbQryPermlyDelId(context)
	// 新id
	if err != nil || count == 0 {
		id, err = db.Add(context, "INSERT INTO `img` (`name`, `type`, `size`, `add_time`) VALUES (?, ?, ?, ?)", name, _type, size, time.NowUnix())
	} else
	// 复用id
	{
		_, err = db.Upd(context, "UPDATE `img` SET `name` = ?, `type` = ?, `size` = ?, `hist` = '', `hist_size` = 0, `del` = 0, `add_time` = ?, `upd_time` = 0 WHERE `id` = ?", name, _type, size, time.NowUnix(), id)
	}
	img := typ.Img{Abs: typ.Abs{Id: id}, Name: name, Type: _type}
	if err != nil {
		redirect(img, err)
		return
	}

	// 清空图片
	path, err := ClearImg(context, img)
	if err != nil {
		redirect(img, err)
		return
	}

	// 保存文件
	err = context.SaveUploadedFile(fh, path)

	// redirect
	redirect(img, err)
}

// DbQryPermlyDelId 查询永久删除的图片记录id，以复用
func DbQryPermlyDelId(context *gin.Context) (int64, int64, error) {
	id, count, err := db.Qry[int64](context, "SELECT `id` FROM `img` WHERE `del` = 2 LIMIT 1")
	return id, count, err
}
