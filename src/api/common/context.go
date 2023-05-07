// context
// @author xiangqian
// @date 20:04 2023/03/22
package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	typ_resp "note/src/typ"
	util_reflect "note/src/util/reflect"
	util_str "note/src/util/str"
	util_time "note/src/util/time"
	"os"
	"strings"
)

func HtmlNotFound[T any](context *gin.Context, name string, resp typ_resp.Resp[T]) {
	Html[T](context, http.StatusNotFound, name, resp)
}

func HtmlOk[T any](context *gin.Context, name string, resp typ_resp.Resp[T]) {
	Html[T](context, http.StatusOK, name, resp)
}

// Html
// name: templateName
func Html[T any](context *gin.Context, code int, name string, resp typ_resp.Resp[T]) {
	// resp msg
	resp0, err := GetSessionV[any](context, RespSessionKey, true)
	if err == nil {
		msg := util_reflect.CallField[string](resp0, "Msg")
		if msg != "" {
			resp.Msg = fmt.Sprintf("%s\n%s", msg, resp.Msg)
		}
	}

	// user
	user, _ := GetSessionUser(context)

	// url
	url := context.Request.RequestURI

	// html
	context.HTML(code, name, gin.H{
		RespSessionKey: resp,
		UserSessionKey: user,
		UrlSessionKey:  url,
	})
}

func JsonBadRequest[T any](context *gin.Context, resp typ_resp.Resp[T]) {
	Json(context, http.StatusBadRequest, resp)
}

func JsonOk[T any](context *gin.Context, resp typ_resp.Resp[T]) {
	Json(context, http.StatusOK, resp)
}

func Json[T any](context *gin.Context, code int, resp typ_resp.Resp[T]) {
	context.JSON(code, resp)
}

func Redirect[T any](context *gin.Context, location string, resp typ_resp.Resp[T]) {
	SetSessionKv(context, RespSessionKey, resp)
	if strings.Contains(location, "?") {
		location = fmt.Sprintf("%s&t=%d", location, util_time.NowUnix())
	} else {
		location = fmt.Sprintf("%s?t=%d", location, util_time.NowUnix())
	}
	context.Redirect(http.StatusMovedPermanently, location)
}

func PostForm[T any](context *gin.Context, key string) (T, error) {
	value := context.PostForm(key)
	return util_str.ConvStrToType[T](value)
}

func Param[T any](context *gin.Context, key string) (T, error) {
	value := context.Param(key)
	return util_str.ConvStrToType[T](value)
}

func Query[T any](context *gin.Context, key string) (T, error) {
	value := context.Query(key)
	return util_str.ConvStrToType[T](value)
}

func PageReq(context *gin.Context) (typ_resp.Req, error) {
	req := typ_resp.Req{Size: 10}
	err := ShouldBind(context, &req)
	if req.Current <= 0 {
		req.Current = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	return req, err
}

func ShouldBindQuery(context *gin.Context, i any) error {
	err := context.ShouldBindQuery(i)
	if err != nil {
		err = TransErr(context, err)
	}
	return err
}

// ShouldBind 应该绑定参数
func ShouldBind(context *gin.Context, i any) error {
	err := context.ShouldBind(i)
	if err != nil {
		err = TransErr(context, err)
	}
	return err
}

func Write(context *gin.Context, path string) {
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
