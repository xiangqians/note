// img upload
// @author xiangqian
// @date 21:26 2023/04/10
package img

import (
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"note/src/api/common"
	"note/src/typ"
	util_os "note/src/util/os"
	util_str "note/src/util/str"
	util_time "note/src/util/time"
	util_validate "note/src/util/validate"
	"os"
	"strings"
)

// ReUpload 重新上传图片
func ReUpload(context *gin.Context) {
	redirect := func(id int64, err any) {
		resp := typ.Resp[any]{Msg: util_str.ConvTypeToStr(err)}
		common.Redirect(context, fmt.Sprintf("/img/%d/view", id), resp)
	}

	// id
	id, err := common.PostForm[int64](context, "id")
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

	// file name
	name := strings.TrimSpace(fh.Filename)
	err = util_validate.FileName(name)
	if err != nil {
		redirect(id, err)
		return
	}

	// file type
	contentType := fh.Header.Get("Content-Type")
	ft := typ.ContentTypeOf(contentType)
	if !typ.IsImg(ft) {
		redirect(id, fmt.Sprintf("%s, %s", i18n.MustGetMessage("i18n.fileTypeUnsupported"), contentType))
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

	// 备份历史记录
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
	_, err = util_os.CopyFile(dstPath, srcPath)
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
				util_os.DelFile(path)
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
			UpdTime: util_time.NowUnix(),
		},
		Name:     name,
		Type:     _type,
		Size:     size,
		Hist:     hist,
		HistSize: histSize,
	}

	// update
	_, err = common.DbUpd(context, "UPDATE `img` SET `name` = ?, `type` = ?, `size` = ?, `hist` = ?, `hist_size` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ?",
		newImg.Name, newImg.Type, newImg.Size, newImg.Hist, newImg.HistSize, newImg.UpdTime, newImg.Id)
	if err != nil {
		redirect(id, err)
		return
	}

	// path
	path, err := Path(context, newImg)
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
	if img.Type != newImg.Type {
		path, err = Path(context, img)
		if err == nil {
			util_os.DelFile(path)
		}
	}

	// redirect
	redirect(id, err)
}

// Upload 上传图片
func Upload(context *gin.Context) {
	redirect := func(err any) {
		resp := typ.Resp[any]{Msg: util_str.ConvTypeToStr(err)}
		common.Redirect(context, fmt.Sprintf("/img/list"), resp)
	}

	// file header
	fh, err := context.FormFile("file")
	if err != nil || fh == nil {
		redirect(err)
		return
	}

	// file name
	name := strings.TrimSpace(fh.Filename)
	err = util_validate.FileName(name)
	if err != nil {
		redirect(err)
		return
	}

	// file type
	contentType := fh.Header.Get("Content-Type")
	ft := typ.ContentTypeOf(contentType)
	if !typ.IsImg(ft) {
		redirect(fmt.Sprintf("%s, %s", i18n.MustGetMessage("i18n.fileTypeUnsupported"), contentType))
		return
	}
	_type := string(ft)

	// file size
	size := fh.Size

	// 查询是否有永久删除的图片记录id，以复用
	id, count, err := DbQryPermlyDelId(context)
	// 新id
	if err != nil || count == 0 {
		id, err = common.DbAdd(context, "INSERT INTO `img` (`name`, `type`, `size`, `add_time`) VALUES (?, ?, ?, ?)", name, _type, size, util_time.NowUnix())
	} else
	// 复用id
	{
		_, err = common.DbUpd(context, "UPDATE `img` SET `name` = ?, `type` = ?, `size` = ?, `hist` = '', `hist_size` = 0, `del` = 0, `add_time` = ?, `upd_time` = 0 WHERE `id` = ?", name, _type, size, util_time.NowUnix(), id)
	}
	if err != nil {
		redirect(err)
		return
	}

	// path
	img := typ.Img{}
	img.Id = id
	img.Type = _type
	path, err := Path(context, img)
	if err != nil {
		redirect(err)
		return
	}

	// 清空文件
	if util_os.IsExist(path) {
		var file *os.File
		file, err = os.OpenFile(path,
			os.O_WRONLY|os.O_TRUNC, // 只写（O_WRONLY） & 清空文件（O_TRUNC）
			0666)
		if err != nil {
			redirect(err)
			return
		}
		file.Close()
	}

	// 保存文件
	err = context.SaveUploadedFile(fh, path)

	// redirect
	redirect(err)
}

// DbQryPermlyDelId 查询永久删除的图片记录id，以复用
func DbQryPermlyDelId(context *gin.Context) (int64, int64, error) {
	id, count, err := common.DbQry[int64](context, "SELECT `id` FROM `img` WHERE `del` = 2 LIMIT 1")
	return id, count, err
}
