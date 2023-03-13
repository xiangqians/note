// img
// @author xiangqian
// @date 11:34 2023/02/12
package api

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"note/src/typ"
	"note/src/util"
	"os"
	"strings"
	"time"
)

// ImgListPage 图片列表页面
func ImgListPage(pContext *gin.Context) {
	req, _ := PageReq(pContext)
	page, err := DbPage[typ.Img](pContext, req, "SELECT i.`id`, i.`name`, i.`type`, i.`size`, i.`add_time`, i.`upd_time` FROM `img` i WHERE i.`del` = 0 ORDER BY i.`add_time` DESC")
	HtmlOk(pContext, "img/list.html", gin.H{"page": page}, err)
}

// ImgUpload 图片上传
func ImgUpload(pContext *gin.Context) {
	method := pContext.Request.Method
	redirect := func(id int64, msg any) {
		switch method {
		case http.MethodPost:
			Redirect(pContext, fmt.Sprintf("/img/listpage"), nil, msg)

		case http.MethodPut:
			Redirect(pContext, fmt.Sprintf("/img/%d/editpage", id), nil, msg)
		}
	}

	// id
	var id int64
	var err error
	if method == http.MethodPut {
		id, err = PostForm[int64](pContext, "id")
		if err != nil {
			redirect(id, err)
			return
		}
	}

	// fh
	fh, err := pContext.FormFile("file")
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
	ft := FileTypeImgOf(ftStr)
	if ft == typ.FileTypeUnk {
		redirect(id, fmt.Sprintf("%s, %s", errors.New(i18n.MustGetMessage("i18n.fileTypeUnsupported")), ftStr))
		return
	}

	// size
	fs := fh.Size

	// 限制只能重传相同文件类型
	if method == http.MethodPut {
		img, count, err := ImgQry(pContext, id)
		if err != nil || count == 0 {
			redirect(id, err)
			return
		}
		if ft != FileTypeImgOf(img.Type) {
			redirect(id, "重传必须是相同文件类型")
			return
		}
	}

	// 操作数据库
	switch method {
	case http.MethodPost:
		id, err = DbAdd(pContext, "INSERT INTO `img` (`name`, `type`, `size`, `add_time`) VALUES (?, ?, ?, ?)",
			fn, ft, fs, time.Now().Unix())

	case http.MethodPut:
		_, err = DbUpd(pContext, "UPDATE `img` SET `name` = ?, `type` = ?, `size` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ?",
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
	fp, err := ImgPath(pContext, img)
	if err != nil {
		redirect(id, err)
		return
	}

	// 清空文件
	if method == http.MethodPut && util.IsExistOfPath(fp) {
		pFile, err := os.OpenFile(fp,
			os.O_WRONLY|os.O_TRUNC, // 只写（O_WRONLY） & 清空文件（O_TRUNC）
			0666)
		if err != nil {
			redirect(id, err)
			return
		}
		pFile.Close()
	}

	// 保存文件
	err = pContext.SaveUploadedFile(fh, fp)

	redirect(id, err)
	return
}

func ImgReUpload(pContext *gin.Context) {
	pContext.Request.Method = http.MethodPut
	ImgUpload(pContext)
}

// ImgUpdName 图片重命名
func ImgUpdName(pContext *gin.Context) {
	redirect := func(msg any) {
		Redirect(pContext, fmt.Sprintf("/img/listpage"), nil, msg)
	}

	// img
	img := typ.Img{}
	err := ShouldBind(pContext, &img)
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

	imgType, count, err := DbQry[string](pContext, "SELECT `type` FROM `img` WHERE `del` = 0 AND `id` = ?", img.Id)
	if count > 0 {
		name := img.Name
		ft := FileTypeImgOf(imgType)
		if ft != typ.FileTypeUnk && !strings.HasSuffix(name, string(ft)) {
			name = fmt.Sprintf("%s.%s", name, string(ft))
		}

		// update
		_, err = DbUpd(pContext, "UPDATE `img` SET `name` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `name` <> ?", name, time.Now().Unix(), img.Id, name)
	}

	redirect(err)
	return
}

func ImgDel(pContext *gin.Context) {
	redirect := func(msg any) {
		Redirect(pContext, fmt.Sprintf("/img/listpage"), nil, msg)
	}

	// id
	id, err := Param[int64](pContext, "id")
	if err != nil {
		redirect(err)
		return
	}

	// delete
	_, err = DbDel(pContext, "UPDATE `img` SET `del` = 1, `upd_time` = ? WHERE `id` = ?", time.Now().Unix(), id)
	redirect(err)
	return
}

// ImgView 查看图片
func ImgView(pContext *gin.Context) {
	// id
	id, err := Param[int64](pContext, "id")
	if err != nil {
		log.Println(err)
		return
	}

	// img
	img, count, err := ImgQry(pContext, id)
	if err != nil || count == 0 {
		log.Println(err)
		return
	}

	// path
	path, err := ImgPath(pContext, img)
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

	pContext.Writer.Write(buf)
	return
}

func ImgViewPage(pContext *gin.Context) {
	html := func(img typ.Img, msg any) {
		HtmlOk(pContext, "img/view.html", gin.H{"img": img}, msg)
	}

	// id
	id, err := Param[int64](pContext, "id")
	if err != nil {
		html(typ.Img{}, err)
		return
	}

	// img
	img, _, err := ImgQry(pContext, id)
	img.Url = fmt.Sprintf("/img/%v/view", id)

	html(img, err)
	return
}

func ImgEditPage(pContext *gin.Context) {
	html := func(img typ.Img, msg any) {
		HtmlOk(pContext, "img/edit.html", gin.H{"img": img}, msg)
	}

	// id
	id, err := Param[int64](pContext, "id")
	if err != nil {
		html(typ.Img{}, err)
		return
	}

	// img
	img, _, err := ImgQry(pContext, id)
	html(img, err)
	return
}

// ImgPath 图片物理路径
// id: 图片主键id
func ImgPath(pContext *gin.Context, img typ.Img) (string, error) {
	// dir
	dataDir := DataDir(pContext)
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

func ImgQry(pContext *gin.Context, id int64) (typ.Img, int64, error) {
	img, count, err := DbQry[typ.Img](pContext, "SELECT i.`id`, i.`name`, i.`type`, i.`size`, i.`add_time`, i.`upd_time` FROM `img` i WHERE i.`del` = 0 AND i.`id` = ?", id)
	return img, count, err
}

var imgFileTypes = [...]typ.FileType{typ.FileTypeIco, typ.FileTypeGif, typ.FileTypeJpg, typ.FileTypeJpeg, typ.FileTypePng, typ.FileTypeWebp}

func FileTypeImgOf(value string) typ.FileType {
	for _, imgFileType := range imgFileTypes {
		if strings.ToLower(string(imgFileType)) == strings.ToLower(value) {
			return imgFileType
		}
	}

	return typ.FileTypeUnk
}
