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
	"note/src/model"
	"note/src/session"
	util_string "note/src/util/string"
	util_time "note/src/util/time"
	"reflect"
	"strconv"
	"strings"
)

var (
	zhTrans ut.Translator
	enTrans ut.Translator
)

const redirectMsgKey = "redirectMsg"

func init() {
	// 初始化翻译器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		uni := ut.New(zh.New(), // 备用语言
			// 支持语言
			zh.New(),
			en.New())
		if trans, ok := uni.GetTranslator(model.Zh); ok {
			zh_trans.RegisterDefaultTranslations(v, trans)
			zhTrans = trans
		}
		if trans, ok := uni.GetTranslator(model.En); ok {
			en_trans.RegisterDefaultTranslations(v, trans)
			enTrans = trans
		}
	}
}

func HtmlNotFound[T any](ctx *gin.Context, name string, resp model.Resp[T]) {
	Html[T](ctx, http.StatusNotFound, name, resp)
}

func HtmlOk[T any](ctx *gin.Context, name string, resp model.Resp[T]) {
	Html[T](ctx, http.StatusOK, name, resp)
}

// Html html模板
// ctx  : *gin.Context
// code : http状态码
// name : html模板名称
// resp : 响应数据
func Html[T any](ctx *gin.Context, code int, name string, resp model.Resp[T]) {
	// 从session中获取重定向消息
	redirectMsg, _ := session.Get[string](ctx, redirectMsgKey, true)
	if redirectMsg != "" {
		resp.Msg = redirectMsg + " " + resp.Msg
	}

	// 当前登录用户信息
	user, _ := session.GetUser(ctx)

	// 请求信息
	request := ctx.Request

	// html模板
	ctx.HTML(code, name, gin.H{
		"contextPath": model.GetArg().ContextPath, // 上下文路径
		"path":        request.URL.Path,           // 请求路径
		"uri":         request.RequestURI,         // 请求uri地址
		"user":        user,                       // 登录用户信息
		"resp":        resp,                       // 响应数据
	})
}

func JsonBadRequest[T any](ctx *gin.Context, resp model.Resp[T]) {
	Json(ctx, http.StatusBadRequest, resp)
}

func JsonOk[T any](ctx *gin.Context, resp model.Resp[T]) {
	Json(ctx, http.StatusOK, resp)
}

func Json[T any](ctx *gin.Context, code int, resp model.Resp[T]) {
	ctx.JSON(code, resp)
}

// Redirect 重定向
// ctx 		: *gin.Context
// location : 重定向地址
// paramMap : 重定向参数映射
// msg		: 重定向消息（没有消息就是最好的消息）
func Redirect(ctx *gin.Context, location string, paramMap map[string]any, msg any) {
	// 存储重定向消息到session中
	session.Set(ctx, redirectMsgKey, util_string.String(msg))

	// 数组容量
	var cap int = 1
	if paramMap != nil {
		cap += len(paramMap)
	}

	// len 0, cap ?
	arr := make([]string, 0, cap)

	// 【请求参数】more
	if paramMap != nil && len(paramMap) > 0 {
		for k, v := range paramMap {
			arr = append(arr, fmt.Sprintf("%v=%v", k, v))
		}
	}
	// 【请求参数】时间戳
	arr = append(arr, fmt.Sprintf("t=%v", util_time.NowUnix()))

	// 重定向
	ctx.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s%s?%s", model.GetArg().ContextPath, location, strings.Join(arr, "&")))
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
	if errs, ok := err.(validator.ValidationErrors); ok {
		session := session.Session(ctx)
		lang := ""
		if v, ok := session.Get("lang").(string); ok {
			lang = v
		}
		var validationErrTrans validator.ValidationErrorsTranslations
		switch lang {
		case model.En:
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
				case model.En:
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
