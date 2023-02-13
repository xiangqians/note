// img
// @author xiangqian
// @date 11:34 2023/02/12
package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
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
	Html(pContext, "img/list.html", gin.H{"page": page}, err)
}

func ImgView(pContext *gin.Context) {
	// img
	img, err := ImgQry(pContext)
	if err != nil {
		log.Println(err)
		return
	}

	path, err := ImgPath(pContext, img.Id)
	if err != nil {
		return
	}

	buf, err := os.ReadFile(path)
	if err != nil {
		return
	}

	pContext.Writer.Write(buf)
	return
}

func ImgViewPage(pContext *gin.Context) {
	html := func(img typ.Img, msg any) {
		Html(pContext, "img/view.html", gin.H{"img": img}, msg)
	}

	// id
	id, err := Param[int64](pContext, "id")
	if err != nil {
		html(typ.Img{}, err)
		return
	}

	// img
	img, _, err := DbQry[typ.Img](pContext, "SELECT i.`id`, i.`name`, i.`type`, i.`size`, i.`add_time`, i.`upd_time` FROM `img` i WHERE i.`del` = 0 AND i.`id` = ?", id)
	img.Url = fmt.Sprintf("/img/%v/view", id)
	html(img, err)
	return
}

// ImgUpload 图片上传
func ImgUpload(pContext *gin.Context) {
	redirect := func(msg any) {
		Redirect(pContext, fmt.Sprintf("/img/listpage"), nil, msg)
	}

	// file
	file, _ := pContext.FormFile("file")
	if file == nil {
		redirect(nil)
		return
	}

	//log.Printf("Filename: %v\n", file.Filename)
	//log.Printf("Size: %v\n", file.Size)
	//log.Printf("Header: %v\n", file.Header)

	// img
	fName := file.Filename
	index := strings.LastIndex(fName, ".")
	fTypeStr := ""
	if index > 0 {
		fTypeStr = fName[index+1:]
	}
	fType := string(typ.FileTypeOf(fTypeStr))
	if !(fType == typ.FileTypeIco ||
		fType == typ.FileTypeJpg ||
		fType == typ.FileTypeJpeg ||
		fType == typ.FileTypePng ||
		fType == typ.FileTypeWebp) {
		redirect(fmt.Sprintf("不支持此类型文件上传：%s", fTypeStr))
		return
	}
	fSize := file.Size

	// add
	id, err := DbAdd(pContext, "INSERT INTO `img` (`name`, `type`, `size`, `add_time`) VALUES (?, ?, ?, ?)", fName, fType, fSize, time.Now().Unix())
	if err != nil {
		redirect(err)
		return
	}

	// img path
	fPath, err := ImgPath(pContext, id)
	if err != nil {
		redirect(err)
		return
	}

	// 保存文件
	err = pContext.SaveUploadedFile(file, fPath)

	redirect(err)
	return
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

	// update
	_, err = DbUpd(pContext, "UPDATE `img` SET `name` = ?, `upd_time` = ? WHERE id = ?", img.Name, time.Now().Unix(), img.Id)
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

// ImgPath 图片物理路径
// id: 图片主键id
func ImgPath(pContext *gin.Context, id int64) (string, error) {
	// dir
	dataDir := DataDir(pContext)
	imgDir := fmt.Sprintf("%s%simg", dataDir, util.FileSeparator)
	if !util.IsExistOfPath(imgDir) {
		err := util.Mkdir(imgDir)
		if err != nil {
			return "", err
		}
	}

	// path
	return fmt.Sprintf("%s%s%v", imgDir, util.FileSeparator, id), nil
}

func ImgQry(pContext *gin.Context) (typ.Img, error) {
	// id
	id, err := Param[int64](pContext, "id")
	if err != nil {
		return typ.Img{}, err
	}

	// img
	img, _, err := DbQry[typ.Img](pContext, "SELECT i.`id`, i.`name`, i.`type`, i.`size`, i.`add_time`, i.`upd_time` FROM `img` i WHERE i.`del` = 0 AND i.`id` = ?", id)
	return img, err
}
