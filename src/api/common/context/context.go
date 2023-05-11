// context
// @author xiangqian
// @date 20:04 2023/03/22
package context

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"note/src/api/common/session"
	"note/src/api/common/trans"
	"note/src/typ"
	"note/src/util/reflect"
	"note/src/util/str"
	"note/src/util/time"
	"os"
	"strings"
)

const UserSessionKey = "user"
const RespSessionKey = "resp"
const UrlSessionKey = "url"

func HtmlNotFound[T any](context *gin.Context, name string, resp typ.Resp[T]) {
	Html[T](context, http.StatusNotFound, name, resp)
}

func HtmlOk[T any](context *gin.Context, name string, resp typ.Resp[T]) {
	Html[T](context, http.StatusOK, name, resp)
}

// Html
// name: templateName
func Html[T any](context *gin.Context, code int, name string, resp typ.Resp[T]) {
	// resp msg
	storedResp, err := session.Get[any](context, RespSessionKey, true)
	if err == nil {
		msg := reflect.CallField[string](storedResp, "Msg")
		if msg != "" {
			resp.Msg = fmt.Sprintf("%s\n%s", msg, resp.Msg)
		}
	}

	// user
	user, _ := session.GetUser(context)

	// url
	url := context.Request.RequestURI

	// html
	context.HTML(code, name, gin.H{
		RespSessionKey: resp,
		UserSessionKey: user,
		UrlSessionKey:  url,
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
	session.Set(context, RespSessionKey, resp)
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

// WriteFile 写入文件
// path: 文件路径
func WriteFile(context *gin.Context, path string) (written int, err error) {
	// read
	buf, err := os.ReadFile(path)
	if err != nil {
		return
	}

	// write
	written, err = context.Writer.Write(buf)
	return
}
