// context
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
	"note/src/session"
	"note/src/typ"
	"note/src/util/str"
	"note/src/util/time"
	"strings"
)

var (
	zhTrans ut.Translator
	enTrans ut.Translator
)

const RespSessionKey = "resp"

func init() {
	translator()
}

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
		err = transErr(context, err)
	}
	return err
}

// ShouldBind 应该绑定参数
func ShouldBind(context *gin.Context, i any) error {
	err := context.ShouldBind(i)
	if err != nil {
		err = transErr(context, err)
	}
	return err
}

// translator 检验器翻译
func translator() {
	if v, r := binding.Validator.Engine().(*validator.Validate); r {
		uni := ut.New(zh.New(), // 备用语言
			// 支持的语言
			zh.New(),
			en.New())
		if trans, r := uni.GetTranslator(typ.Zh); r {
			zh_trans.RegisterDefaultTranslations(v, trans)
			zhTrans = trans
		}
		if trans, r := uni.GetTranslator(typ.En); r {
			en_trans.RegisterDefaultTranslations(v, trans)
			enTrans = trans
		}
	}
}

// transErr 翻译异常
func transErr(context *gin.Context, err error) error {
	if errs, r := err.(validator.ValidationErrors); r {
		session := session.Session(context)
		lang := ""
		if v, r := session.Get("lang").(string); r {
			lang = v
		}
		var validationErrTrans validator.ValidationErrorsTranslations
		switch lang {
		case typ.En:
			validationErrTrans = errs.Translate(enTrans)
		default:
			validationErrTrans = errs.Translate(zhTrans)
		}

		errMsg := ""
		for key, value := range validationErrTrans {
			name := key[strings.Index(key, ".")+1:]
			msg, ierr := i18n.GetMessage(fmt.Sprintf("i18n.%s", strings.ToLower(name)))
			if ierr == nil {
				value = strings.Replace(value, name, msg, 1)
			}
			if errMsg != "" {
				switch lang {
				case typ.En:
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
