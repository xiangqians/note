// context
// @author xiangqian
// @date 20:04 2023/03/22
package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"note/src/typ"
	"note/src/util"
)

func HtmlNotFound[T any](context *gin.Context, name string, resp typ.Resp[T]) {
	HtmlNew(context, http.StatusNotFound, name, resp)
}

func HtmlOkNew[T any](context *gin.Context, name string, resp typ.Resp[T]) {
	HtmlNew(context, http.StatusOK, name, resp)
}

// HtmlNew
// name: templateName
func HtmlNew[T any](context *gin.Context, code int, name string, resp typ.Resp[T]) {
	// user
	user, _ := GetSessionUser(context)
	// uri
	uri := context.Request.RequestURI
	// url
	url := context.Request.URL.Path
	// html
	context.HTML(code, name, gin.H{
		"resp": resp,
		"user": user,
		"uri":  uri,
		"url":  url,
	})
}

func HtmlOk(context *gin.Context, templateName string, h gin.H, msg any) {
	Html(context, http.StatusOK, templateName, h, msg)
}

func Html(context *gin.Context, code int, templateName string, h gin.H, msg any) {
	if h == nil {
		h = gin.H{}
	}

	// 获取 user
	_, r := h["user"]
	if !r {
		user, _ := GetSessionUser(context)
		h["user"] = user
	}

	// uri
	h["uri"] = context.Request.RequestURI
	// url
	h["url"] = context.Request.URL.Path

	msgStr := util.TypeAsStr(msg)
	sessionMsg, err := GetSessionV[string](context, "msg", true)
	if err == nil {
		if msgStr != "" {
			msgStr += ", "
		}
		msgStr += sessionMsg
	}
	h["msg"] = msgStr

	context.HTML(code, templateName, h)
}

func RedirectNew(context *gin.Context, location string, resp typ.Resp[any]) {
	SetSessionKv(context, "resp", resp)
	context.Redirect(http.StatusMovedPermanently, location)
}

func Redirect(context *gin.Context, location string, h gin.H, msg any) {
	session := Session(context)
	if h != nil {
		for k, v := range h {
			session.Set(k, v)
		}
	}
	session.Set("msg", util.TypeAsStr(msg))
	session.Save()
	context.Redirect(http.StatusMovedPermanently, location)
}

func PostForm[T any](context *gin.Context, key string) (T, error) {
	value := context.PostForm(key)
	return util.StrAsType[T](value)
}

func Param[T any](context *gin.Context, key string) (T, error) {
	value := context.Param(key)
	return util.StrAsType[T](value)
}

func Query[T any](context *gin.Context, key string) (T, error) {
	value := context.Query(key)
	return util.StrAsType[T](value)
}

func PageReq(context *gin.Context) (typ.PageReq, error) {
	req := typ.PageReq{Size: 10}
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
