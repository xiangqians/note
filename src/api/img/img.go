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
	typ_api "note/src/typ/api"
	typ_ft "note/src/typ/ft"
	typ_page "note/src/typ/page"
	typ_resp "note/src/typ/resp"
	util_json "note/src/util/json"
	util_os "note/src/util/os"
	util_str "note/src/util/str"
	util_time "note/src/util/time"
	"os"
	"strings"
	"time"
)

// UpdName 图片重命名
func UpdName(context *gin.Context) {
	redirect := func(msg any) {
		resp := typ_resp.Resp[any]{Msg: util_str.TypeToStr(msg)}
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
	err = util_os.VerifyFileName(img.Name)
	if err != nil {
		redirect(err)
		return
	}

	//imgType, count, err := common.DbQry[string](context, "SELECT `type` FROM `img` WHERE `del` = 0 AND `id` = ?", img.Id)
	//if count > 0 {
	//	name := img.Name
	//	ft := typ.FileTypeImgOf(imgType)
	//	if ft != typ.FileTypeUnk && !strings.HasSuffix(name, string(ft)) {
	//		name = fmt.Sprintf("%s.%s", name, string(ft))
	//	}
	//
	//	// update
	//	_, err = common.DbUpd(context, "UPDATE `img` SET `name` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `name` <> ?", name, time.Now().Unix(), img.Id, name)
	//}

	// update
	_, err = common.DbUpd(context, "UPDATE `img` SET `name` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `name` <> ?", img.Name, time.Now().Unix(), img.Id, img.Name)

	redirect(err)
	return
}

func Del(context *gin.Context) {
	redirect := func(msg any) {
		resp := typ_resp.Resp[any]{Msg: util_str.TypeToStr(msg)}
		common.Redirect(context, fmt.Sprintf("/img/list"), resp)
	}

	// Delete not supported
	redirect("Delete not supported")
	return

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
	img, count, err := DbQry(context, id)
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

	// write
	n, err := context.Writer.Write(buf)
	log.Println("view", path, n, err)
	return
}

// View 查看图片页面
func View(context *gin.Context) {
	html := func(img typ_api.Img, msg any) {
		resp := typ_resp.Resp[typ_api.Img]{
			Msg:  util_str.TypeToStr(msg),
			Data: img,
		}
		common.HtmlOk(context, "img/view.html", resp)
	}

	// id
	id, err := common.Param[int64](context, "id")
	if err != nil {
		html(typ_api.Img{}, err)
		return
	}

	// img
	img, _, err := DbQry(context, id)
	img.Url = fmt.Sprintf("/img/%v", id)

	html(img, err)
	return
}

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
		redirect(id, fmt.Sprintf("%s, %s", errors.New(i18n.MustGetMessage("i18n.fileTypeUnsupported")), contentType))
		return
	}

	// size
	fs := fh.Size

	// 原图片信息
	var oldImg typ_api.Img

	// 操作数据库
	switch method {
	case http.MethodPost:
		id, err = common.DbAdd(context, "INSERT INTO `img` (`name`, `type`, `size`, `add_time`) VALUES (?, ?, ?, ?)", fn, ft, fs, util_time.NowUnix())

	case http.MethodPut:
		// 原图片信息
		var count int64
		oldImg, count, err = DbQry(context, id)
		if err != nil || count == 0 {
			redirect(id, err)
			return
		}

		// 图片历史记录
		hist := oldImg.Hist
		histSize := oldImg.HistSize
		imgHists := make([]typ_api.Img, 0, 1) // len 0, cap ?
		if hist != "" {
			err = util_json.Deserialize(hist, &imgHists)
			if err != nil {
				redirect(id, err)
				return
			}
		}

		// 将原图片添加到历史记录
		imgHists = append(imgHists, typ_api.Img{
			Abs: typ_api.Abs{
				Id:      oldImg.Id,
				AddTime: oldImg.AddTime,
				UpdTime: oldImg.UpdTime,
			},
			Name: oldImg.Name,
			Type: oldImg.Type,
			Size: oldImg.Size,
		})
		// src
		var srcPath string
		srcPath, err = Path(context, oldImg)
		if err != nil {
			redirect(id, err)
			return
		}
		// dst
		var dstPath string
		dstPath, err = HistPath(context, oldImg)
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

		// 图片历史记录至多保存三张，超过则删除最早地历史图片
		maxHist := 3
		if len(imgHists) > maxHist {
			l := len(imgHists) - maxHist
			for i := 0; i < l; i++ {
				path, err := HistPath(context, imgHists[i])
				if err == nil {
					util_os.DelFile(path)
				}
			}
			imgHists = imgHists[l:]
		}

		// hist size
		for _, imgHist := range imgHists {
			histSize += imgHist.Size
		}

		// serialize
		hist, err = util_json.Serialize(imgHists)
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
		file, err := os.OpenFile(path,
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

// List 图片列表页面
func List(context *gin.Context) {
	req, _ := common.PageReq(context)
	page, err := common.DbPage[typ_api.Img](context, req, "SELECT `id`, `name`, `type`, `size`, `add_time`, `upd_time` FROM `img` WHERE `del` = 0 ORDER BY (CASE WHEN `upd_time` > `add_time` THEN `upd_time` ELSE `add_time` END) DESC")
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
	time := img.AddTime
	if img.AddTime < img.UpdTime {
		time = img.UpdTime
	}
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

func DbQry(context *gin.Context, id int64) (typ_api.Img, int64, error) {
	img, count, err := common.DbQry[typ_api.Img](context, "SELECT `id`, `name`, `type`, `size`, `hist`, `hist_size`, `add_time`, `upd_time` FROM `img` WHERE `del` = 0 AND `id` = ?", id)
	return img, count, err
}
