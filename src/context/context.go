// context
// @author xiangqian
// @date 20:04 2023/03/22
package context

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"note/src/session"
	"note/src/trans"
	"note/src/typ"
	"note/src/util/str"
	"note/src/util/time"
	"strings"
)

const RespSessionKey = "resp"

func HtmlNotFound[T any](context *gin.Context, name string, resp typ.Resp[T]) {
	Html[T](context, http.StatusNotFound, name, resp)
}

func HtmlOk[T any](context *gin.Context, name string, resp typ.Resp[T]) {
	Html[T](context, http.StatusOK, name, resp)
}

// Html
// context: Context
// code: http code
// name: templateName
// resp: response
func Html[T any](context *gin.Context, code int, name string, resp typ.Resp[T]) {
	user, _ := session.GetUser(context)
	context.HTML(code, name, gin.H{
		"url":          context.Request.RequestURI, // url
		"user":         user,                       // user
		RespSessionKey: resp,                       // resp
	})
}

func JsonBadRequest[T any](context *gin.Context, resp typ.Resp[T]) {
	Json(context, http.StatusBadRequest, resp)
}

func JsonOk[T any](context *gin.Context, resp typ.Resp[T]) {
	Json(context, http.StatusOK, resp)
}

func Json[T any](context *gin.Context, code int, resp typ.Resp[T]) {
	context.JSON(code, resp)
}

func Redirect[T any](context *gin.Context, location string, resp typ.Resp[T]) {
	if strings.Contains(location, "?") {
		location = fmt.Sprintf("%s&t=%d", location, time.NowUnix())
	} else {
		location = fmt.Sprintf("%s?t=%d", location, time.NowUnix())
	}
	context.Redirect(http.StatusMovedPermanently, location)
}

func PostForm[T any](context *gin.Context, key string) (T, error) {
	value := context.PostForm(key)
	return str.ConvStrToType[T](value)
}

func Param[T any](context *gin.Context, key string) (T, error) {
	value := context.Param(key)
	return str.ConvStrToType[T](value)
}

func Query[T any](context *gin.Context, key string) (T, error) {
	value := context.Query(key)
	return str.ConvStrToType[T](value)
}

func ShouldBindQuery(context *gin.Context, i any) error {
	err := context.ShouldBindQuery(i)
	if err != nil {
		err = trans.Err(context, err)
	}
	return err
}

// ShouldBind 应该绑定参数
func ShouldBind(context *gin.Context, i any) error {
	err := context.ShouldBind(i)
	if err != nil {
		err = trans.Err(context, err)
	}
	return err
}
