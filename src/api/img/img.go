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
	"sort"
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
	err = util_os.VerifyFileName(img.Name)
	if err != nil {
		redirect(err)
		return
	}

	// update
	_, err = common.DbUpd(context, "UPDATE `img` SET `name` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ? AND `name` <> ?", img.Name, util_time.NowUnix(), img.Id, img.Name)

	redirect(err)
	return
}

// HistView 查看图片历史页面
func HistView(context *gin.Context) {
	html := func(img typ_api.Img, err any) {
		resp := typ_resp.Resp[typ_api.Img]{
			Msg:  util_str.TypeToStr(err),
			Data: img,
		}
		common.HtmlOk(context, "img/view.html", resp)
	}

	// id
	id, err := common.Param[int64](context, "id")
	if err != nil || id <= 0 {
		html(typ_api.Img{}, err)
		return
	}

	// idx
	idx, err := common.Param[int](context, "idx")
	if err != nil || idx < 0 {
		html(typ_api.Img{}, err)
		return
	}

	// img
	img, count, err := DbQry(context, id, 0)
	if err != nil || count == 0 {
		html(img, err)
		return
	}

	// 图片历史记录
	hist := img.Hist
	if hist == "" {
		html(img, err)
		return
	}

	// hists
	hists := make([]typ_api.Img, 0, 1) // len 0, cap ?
	err = util_json.Deserialize(hist, &hists)
	if err != nil {
		html(img, err)
		return
	}

	// 校验idx是否合法
	if idx >= len(hists) {
		html(img, err)
		return
	}

	// sort
	sort.Slice(hists, func(i, j int) bool {
		return hists[i].UpdTime > hists[j].UpdTime
	})

	// hist img
	histImg := hists[idx]
	histImg.Url = fmt.Sprintf("/img/%d/hist/%d?t=%d", id, idx, util_time.NowUnix())
	histImg.Hists = hists
	img = histImg

	// html
	html(img, err)
	return
}

// View 查看图片页面
func View(context *gin.Context) {
	html := func(img typ_api.Img, err any) {
		resp := typ_resp.Resp[typ_api.Img]{
			Msg:  util_str.TypeToStr(err),
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
	img, count, err := DbQry(context, id, 0)
	if err != nil || count == 0 {
		html(img, err)
		return
	}

	// url
	img.Url = fmt.Sprintf("/img/%d?t=%d", id, util_time.NowUnix())

	// 图片历史记录
	hist := img.Hist
	if hist != "" {
		// hists
		hists := make([]typ_api.Img, 0, 1) // len 0, cap ?
		err = util_json.Deserialize(hist, &hists)
		if err != nil {
			html(img, err)
			return
		}

		// sort
		sort.Slice(hists, func(i, j int) bool {
			return hists[i].UpdTime > hists[j].UpdTime
		})

		img.Hists = hists
	}

	// html
	html(img, err)
	return
}

// GetHist 获取历史图片
func GetHist(context *gin.Context) {
	// id
	id, err := common.Param[int64](context, "id")
	if err != nil {
		log.Println(err)
		return
	}

	// idx
	idx, err := common.Param[int](context, "idx")
	if err != nil || idx < 0 {
		log.Println(err)
		return
	}

	// img
	img, count, err := DbQry(context, id, 0)
	if err != nil || count == 0 {
		log.Println(err)
		return
	}

	// hist
	hist := img.Hist
	if hist == "" {
		log.Println(err)
		return
	}

	// hists
	hists := make([]typ_api.Img, 0, 1) // len 0, cap ?
	err = util_json.Deserialize(hist, &hists)
	if err != nil {
		log.Println(err)
		return
	}

	// 校验idx是否合法
	if idx >= len(hists) {
		log.Println(err)
		return
	}

	// sort
	sort.Slice(hists, func(i, j int) bool {
		return hists[i].UpdTime > hists[j].UpdTime
	})

	// hist img
	histImg := hists[idx]
	path, err := HistPath(context, histImg)
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
	log.Println(path, n, err)
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
	img, count, err := DbQry(context, id, 0)
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
	log.Println(path, n, err)
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
