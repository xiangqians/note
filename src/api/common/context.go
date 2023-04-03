// context
// @author xiangqian
// @date 20:04 2023/03/22
package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	typ_page "note/src/typ/page"
	typ_resp "note/src/typ/resp"
	util_reflect "note/src/util/reflect"
	util_str "note/src/util/str"
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
	// msg
	resp0, err := GetSessionV[any](context, RespSessionKey, true)
	if err == nil {
		msg := util_reflect.CallField[string](resp0, "Msg")
		if msg != "" {
			resp.Msg = fmt.Sprintf("%s\n%s", msg, resp.Msg)
		}
	}

	// user
	user, _ := GetSessionUser(context)
	// uri
	uri := context.Request.RequestURI
	// url
	url := context.Request.URL.Path
	// html
	context.HTML(code, name, gin.H{
		RespSessionKey: resp,
		UserSessionKey: user,
		UriSessionKey:  uri,
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
	context.Redirect(http.StatusMovedPermanently, location)
}

func PostForm[T any](context *gin.Context, key string) (T, error) {
	value := context.PostForm(key)
	return util_str.StrToType[T](value)
}

func Param[T any](context *gin.Context, key string) (T, error) {
	value := context.Param(key)
	return util_str.StrToType[T](value)
}

func Query[T any](context *gin.Context, key string) (T, error) {
	value := context.Query(key)
	return util_str.StrToType[T](value)
}

func PageReq(context *gin.Context) (typ_page.PageReq, error) {
	req := typ_page.PageReq{Size: 10}
	err := ShouldBind(context, &req)
	if req.Current <= 0 {
		req.Current = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	return req, err
}

// ShouldBind 应该绑定参数
func ShouldBind(context *gin.Context, i any) error {
	err := context.ShouldBind(i)
	if err != nil {
		err = TransErr(context, err)
	}
	return err
}
