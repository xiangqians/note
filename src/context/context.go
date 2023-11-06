// @author xiangqian
// @date 20:04 2023/03/22
package context

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_trans "github.com/go-playground/validator/v10/translations/en"
	zh_trans "github.com/go-playground/validator/v10/translations/zh"
	"net/http"
	"note/src/model"
	"note/src/session"
	util_time "note/src/util/time"
	"strings"
)

var (
	zhTrans ut.Translator
	enTrans ut.Translator
)

const (
	ZH = "zh"
	EN = "en"
)

const redirectMsgKey = "redirectMsg"

func init() {
	// 初始化翻译器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		uni := ut.New(zh.New(), // 备用语言
			// 支持语言
			zh.New(),
			en.New())
		if trans, ok := uni.GetTranslator(ZH); ok {
			zh_trans.RegisterDefaultTranslations(v, trans)
			zhTrans = trans
		}
		if trans, ok := uni.GetTranslator(EN); ok {
			en_trans.RegisterDefaultTranslations(v, trans)
			enTrans = trans
		}
	}
}

func HtmlNotFound(ctx *gin.Context, name string, resp model.Resp) {
	Html(ctx, http.StatusNotFound, name, resp)
}

func HtmlOk(ctx *gin.Context, name string, resp model.Resp) {
	Html(ctx, http.StatusOK, name, resp)
}

// Html html模板
// ctx  : *gin.Context
// code : http状态码
// name : html模板名称
// resp : 响应数据
func Html(ctx *gin.Context, code int, name string, resp model.Resp) {
	// 从session中获取重定向消息
	redirectMsg := session.GetString(ctx, redirectMsgKey, true)
	if redirectMsg != "" {
		resp.Msg = redirectMsg + " " + resp.Msg
	}

	// 请求信息
	request := ctx.Request

	// html模板
	ctx.HTML(code, name, gin.H{
		"contextPath": model.GetArg().ContextPath, // 上下文路径
		"path":        request.URL.Path,           // 请求路径
		"uri":         request.RequestURI,         // 请求uri地址
		"user":        session.GetUser(ctx),       // 当前登录用户信息
		"resp":        resp,                       // 响应数据
	})
}

func JsonBadRequest(ctx *gin.Context, resp model.Resp) {
	Json(ctx, http.StatusBadRequest, resp)
}

func JsonOk(ctx *gin.Context, resp model.Resp) {
	Json(ctx, http.StatusOK, resp)
}

func Json(ctx *gin.Context, code int, resp model.Resp) {
	ctx.JSON(code, resp)
}

// Redirect 重定向
// ctx 		: *gin.Context
// location : 重定向地址
// msg		: 重定向消息（没有消息就是最好的消息）
func Redirect(ctx *gin.Context, location string, msg any) {
	// 存储重定向消息到session中
	session.SetString(ctx, redirectMsgKey, msg)

	// 添加时间戳请求参数
	if strings.Contains(location, "?") {
		location += fmt.Sprintf("&t=%v", util_time.NowUnix())
	} else {
		location += fmt.Sprintf("?t=%v", util_time.NowUnix())
	}

	// 重定向
	ctx.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s%s", model.GetArg().ContextPath, location))
}

// PostForm 获取请求表单参数
func PostForm(ctx *gin.Context, name string) string {
	return ctx.PostForm(name)
}

// Param 获取请求占位符参数
func Param(ctx *gin.Context, name string) string {
	return ctx.Param(name)
}

// Query 获取请求参数
func Query(ctx *gin.Context, name string) string {
	return ctx.Query(name)
}

// Header 获取请求头参数
func Header(ctx *gin.Context, name string) string {
	return ctx.GetHeader(name)
}

// ShouldBindQuery 将请求参数绑定给变量
func ShouldBindQuery(ctx *gin.Context, i any) error {
	err := ctx.ShouldBindQuery(i)
	if err != nil {
		err = transErr(ctx, err)
	}
	return err
}

// ShouldBind 将请求报文体绑定给变量
func ShouldBind(ctx *gin.Context, i any) error {
	err := ctx.ShouldBind(i)
	if err != nil {
		err = transErr(ctx, err)
	}
	return err
}

// transErr 翻译异常
func transErr(ctx *gin.Context, err error) error {
	if errs, ok := err.(validator.ValidationErrors); ok {
		// 语言
		language := session.GetString(ctx, "language", false)

		// 校验异常翻译器
		var validationErrTrans validator.ValidationErrorsTranslations
		switch language {
		case EN:
			validationErrTrans = errs.Translate(enTrans)
		default:
			validationErrTrans = errs.Translate(zhTrans)
		}

		errMsg := ""
		for name, value := range validationErrTrans {
			name = name[strings.Index(name, ".")+1:]
			if msg, err := i18n.GetMessage(fmt.Sprintf("i18n.%s", strings.ToLower(name))); err == nil {
				value = strings.Replace(value, name, msg, 1)
			}
			if errMsg != "" {
				switch language {
				case EN:
					errMsg += ", "
				default:
					errMsg += "、"
				}
			}
			errMsg += value
		}
		return errors.New(errMsg)
	}
	return err
}
