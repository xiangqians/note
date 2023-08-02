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
	"note/src/app"
	"note/src/session"
	"note/src/typ"
	"note/src/util/time"
	"reflect"
	"strconv"
	"strings"
)

var (
	zhTrans ut.Translator
	enTrans ut.Translator
)

func init() {
	// 初始化翻译器
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

func HtmlNotFound[T any](ctx *gin.Context, name string, resp typ.Resp[T]) {
	Html[T](ctx, http.StatusNotFound, name, resp)
}

func HtmlOk[T any](ctx *gin.Context, name string, resp typ.Resp[T]) {
	Html[T](ctx, http.StatusOK, name, resp)
}

// Html
// ctx  : Context
// code : http code
// name : templateName
// resp : response
func Html[T any](ctx *gin.Context, code int, name string, resp typ.Resp[T]) {
	user, _ := session.GetUser(ctx)
	ctx.HTML(code, name, gin.H{
		"resp": resp,                   // 响应数据
		"url":  ctx.Request.RequestURI, // 请求url地址
		"user": user,                   // 登录用户信息
		"arg":  app.GetArg(),           // 应用参数
	})
}

func JsonBadRequest[T any](ctx *gin.Context, resp typ.Resp[T]) {
	Json(ctx, http.StatusBadRequest, resp)
}

func JsonOk[T any](ctx *gin.Context, resp typ.Resp[T]) {
	Json(ctx, http.StatusOK, resp)
}

func Json[T any](ctx *gin.Context, code int, resp typ.Resp[T]) {
	ctx.JSON(code, resp)
}

func Redirect(ctx *gin.Context, location string) {
	if strings.Contains(location, "?") {
		location = fmt.Sprintf("%s&t=%d", location, time.NowUnix())
	} else {
		location = fmt.Sprintf("%s?t=%d", location, time.NowUnix())
	}
	ctx.Redirect(http.StatusMovedPermanently, location)
}

func PostForm[T any](ctx *gin.Context, key string) (T, error) {
	value := ctx.PostForm(key)
	return convStrToType[T](value)
}

func Param[T any](ctx *gin.Context, key string) (T, error) {
	value := ctx.Param(key)
	return convStrToType[T](value)
}

func Query[T any](ctx *gin.Context, key string) (T, error) {
	value := ctx.Query(key)
	return convStrToType[T](value)
}

// convStrToType string转类型（基本数据类型）
func convStrToType[T any](value string) (T, error) {
	var t T
	rflVal := reflect.ValueOf(t)
	//log.Println(rflVal)
	switch rflVal.Type().Kind() {
	case reflect.Int:
		id, err := strconv.ParseInt(value, 10, 64)
		return any(int(any(id).(int64))).(T), err

	case reflect.Int8:
		id, err := strconv.ParseInt(value, 10, 64)
		return any(int8(any(id).(int64))).(T), err

	case reflect.Uint8:
		id, err := strconv.ParseInt(value, 10, 64)
		return any(uint8(any(id).(int64))).(T), err

	case reflect.Int64:
		id, err := strconv.ParseInt(value, 10, 64)
		return any(id).(T), err

	case reflect.String:
		return any(strings.TrimSpace(value)).(T), nil
	}

	return t, errors.New(fmt.Sprintf("This type does not support conversion: %v", rflVal.Type().Kind()))
}

func ShouldBindQuery(ctx *gin.Context, i any) error {
	err := ctx.ShouldBindQuery(i)
	if err != nil {
		err = transErr(ctx, err)
	}
	return err
}

func ShouldBind(ctx *gin.Context, i any) error {
	err := ctx.ShouldBind(i)
	if err != nil {
		err = transErr(ctx, err)
	}
	return err
}

// transErr 翻译异常
func transErr(ctx *gin.Context, err error) error {
	if errs, r := err.(validator.ValidationErrors); r {
		session := session.Default(ctx)
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
