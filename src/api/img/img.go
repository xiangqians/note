// img
// @author xiangqian
// @date 11:34 2023/02/12
package img

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"note/src/api/common"
	"note/src/typ"
	"note/src/util"
	"os"
	"strings"
	"time"
)

// List 图片列表页面
func List(context *gin.Context) {
	req, _ := common.PageReq(context)
	page, err := common.DbPage[typ.Img](context, req, "SELECT i.`id`, i.`name`, i.`type`, i.`size`, i.`add_time`, i.`upd_time` FROM `img` i WHERE i.`del` = 0 ORDER BY i.`add_time`, i.`upd_time` DESC")
	common.HtmlOkNew(context, "img/list.html", typ.Resp[typ.Page[typ.Img]]{
		Msg:  util.TypeAsStr(err),
		Data: page,
	})
}

// Upload 图片上传
func Upload(context *gin.Context) {
	method := context.Request.Method
	redirect := func(id int64, msg any) {
		resp := typ.Resp[any]{Msg: util.TypeAsStr(msg)}
		switch method {
		// 上传
		case http.MethodPost:
			common.RedirectNew(context, fmt.Sprintf("/img/list"), resp)

		// 重传
		case http.MethodPut:
			common.RedirectNew(context, fmt.Sprintf("/img/%d/edit", id), resp)
		}
	}

	// id
	var id int64
	var err error
	if method == http.MethodPut {
		id, err = common.PostForm[int64](context, "id")
		if err != nil {
			redirect(id, err)
			return
		}
	}

	// fh
	fh, err := context.FormFile("file")
	if err != nil || fh == nil {
		redirect(id, err)
		return
	}

	// name
	fn := fh.Filename
	fn = strings.TrimSpace(fn)
	err = util.VerifyFileName(fn)
	if err != nil {
		redirect(id, err)
		return
	}

	// type
	index := strings.LastIndex(fn, ".")
	ftStr := ""
	if index > 0 {
		ftStr = fn[index+1:]
	}
	ft := typ.FileTypeImgOf(ftStr)
	if ft == typ.FileTypeUnk {
		redirect(id, fmt.Sprintf("%s, %s", errors.New(i18n.MustGetMessage("i18n.fileTypeUnsupported")), ftStr))
		return
	}

	// size
	fs := fh.Size

	// 限制只能重传相同文件类型
	/*
		if method == http.MethodPut {
			img, count, err := ImgQry(context, id)
			if err != nil || count == 0 {
				redirect(id, err)
				return
			}
			if ft != FileTypeImgOf(img.Type) {
				redirect(id, "重传必须是相同文件类型")
				return
			}
		}
	*/

	// 操作数据库
	switch method {
	case http.MethodPost:
		id, err = common.DbAdd(context, "INSERT INTO `img` (`name`, `type`, `size`, `add_time`) VALUES (?, ?, ?, ?)",
			fn, ft, fs, time.Now().Unix())

	case http.MethodPut:
		_, err = common.DbUpd(context, "UPDATE `img` SET `name` = ?, `type` = ?, `size` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ?",
			fn, ft, fs, time.Now().Unix(), id)
	}
	if err != nil {
		redirect(id, err)
		return
	}

	// path
	img := typ.Img{}
	img.Id = id
	img.Type = string(ft)
	fp, err := Path(context, img)
	if err != nil {
		redirect(id, err)
		return
	}

	// 清空文件
	if method == http.MethodPut && util.IsExistOfPath(fp) {
		file, err := os.OpenFile(fp,
			os.O_WRONLY|os.O_TRUNC, // 只写（O_WRONLY） & 清空文件（O_TRUNC）
			0666)
		if err != nil {
			redirect(id, err)
			return
		}
		file.Close()
	}

	// 保存文件
	err = context.SaveUploadedFile(fh, fp)

	redirect(id, err)
	return
}

func ReUpload(context *gin.Context) {
	context.Request.Method = http.MethodPut
	Upload(context)
}

// UpdName 图片重命名
func UpdName(context *gin.Context) {
	redirect := func(msg any) {
		common.Redirect(context, fmt.Sprintf("/img/listpage"), nil, msg)
	}

	// img
	img := typ.Img{}
	err := common.ShouldBind(context, &img)
	if err != nil {
		redirect(err)
		return
	}

	// name
	img.Name = strings.TrimSpace(img.Name)
	err = util.VerifyFileName(img.Name)
	if err != nil {
		redirect(err)
		return
	}

	imgType, count, err := common.DbQry[string](context, "SELECT `type` FROM `img` WHERE `del` = 0 AND `id` = ?", img.Id)
	if count > 0 {
		name := img.Name
		ft := typ.FileTypeImgOf(imgType)
		if ft != typ.FileTypeUnk && !strings.HasSuffix(name, string(ft)) {
			name = fmt.Sprintf("%s.%s", name, string(ft))
		}

		// update
		_, err = common.DbUpd(context, "UPDATE `img` SET `name` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `name` <> ?", name, time.Now().Unix(), img.Id, name)
	}

	redirect(err)
	return
}

func Del(context *gin.Context) {
	redirect := func(msg any) {
		common.Redirect(context, fmt.Sprintf("/img/listpage"), nil, msg)
	}

	// id
	id, err := common.Param[int64](context, "id")
	if err != nil {
		redirect(err)
		return
	}

	// delete
	_, err = common.DbDel(context, "UPDATE `img` SET `del` = 1, `upd_time` = ? WHERE `id` = ?", time.Now().Unix(), id)
	redirect(err)
	return
}

// Get 获取图片
func Get(context *gin.Context) {
	// id
	id, err := common.Param[int64](context, "id")
	if err != nil {
		log.Println(err)
		return
	}

	// img
	img, count, err := ImgQry(context, id)
	if err != nil || count == 0 {
		log.Println(err)
		return
	}

	// path
	path, err := Path(context, img)
	if err != nil {
		log.Println(err)
		return
	}

	// read
	buf, err := os.ReadFile(path)
	if err != nil {
		log.Println(err)
		return
	}

	context.Writer.Write(buf)
	return
}

func View(context *gin.Context) {
	html := func(img typ.Img, msg any) {
		common.HtmlOk(context, "img/view.html", gin.H{"img": img}, msg)
	}

	// id
	id, err := common.Param[int64](context, "id")
	if err != nil {
		html(typ.Img{}, err)
		return
	}

	// img
	img, _, err := ImgQry(context, id)
	img.Url = fmt.Sprintf("/img/%v/view", id)

	html(img, err)
	return
}

func Edit(context *gin.Context) {
	html := func(img typ.Img, msg any) {
		common.HtmlOk(context, "img/edit.html", gin.H{"img": img}, msg)
	}

	// id
	id, err := common.Param[int64](context, "id")
	if err != nil {
		html(typ.Img{}, err)
		return
	}

	// img
	img, _, err := ImgQry(context, id)
	html(img, err)
	return
}

// Path 图片物理路径
// id: 图片主键id
func Path(context *gin.Context, img typ.Img) (string, error) {
	// dir
	dataDir := common.DataDir(context)
	imgDir := fmt.Sprintf("%s%s%s%s%s", dataDir, util.FileSeparator, "img", util.FileSeparator, img.Type)
	if !util.IsExistOfPath(imgDir) {
		err := util.Mkdir(imgDir)
		if err != nil {
			return "", err
		}
	}

	// path
	return fmt.Sprintf("%s%s%d", imgDir, util.FileSeparator, img.Id), nil
}

func ImgQry(context *gin.Context, id int64) (typ.Img, int64, error) {
	img, count, err := common.DbQry[typ.Img](context, "SELECT i.`id`, i.`name`, i.`type`, i.`size`, i.`add_time`, i.`upd_time` FROM `img` i WHERE i.`del` = 0 AND i.`id` = ?", id)
	return img, count, err
}
