// common
// @author xiangqian
// @date 14:52 2023/02/04
package api

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_trans "github.com/go-playground/validator/v10/translations/en"
	zh_trans "github.com/go-playground/validator/v10/translations/zh"
	"net/http"
	"note/src/arg"
	"note/src/db"
	"note/src/page"
	"note/src/typ"
	"note/src/util"
	"strings"
)

var (
	zhTrans ut.Translator
	enTrans ut.Translator
)

// ValidateTrans 检验器翻译
func ValidateTrans() {
	if v, r := binding.Validator.Engine().(*validator.Validate); r {
		uni := ut.New(zh.New(), // 备用语言
			// 支持的语言
			zh.New(),
			en.New())
		if trans, r := uni.GetTranslator(typ.LocaleZh); r {
			zh_trans.RegisterDefaultTranslations(v, trans)
			zhTrans = trans
		}
		if trans, r := uni.GetTranslator(typ.LocaleEn); r {
			en_trans.RegisterDefaultTranslations(v, trans)
			enTrans = trans
		}
	}
}

// SessionV 根据 key 获取 session value
// key: key
// del: 是否删除 session 中的key
func SessionV[T any](pContext *gin.Context, key any, del bool) (T, error) {
	session := sessions.Default(pContext)
	value := session.Get(key)
	if del {
		session.Delete(key)
		session.Save()
	}

	// t
	if t, r := value.(T); r {
		return t, nil
	}

	// default
	var t T
	return t, errors.New("unknown")
}

// SessionKv 设置 session kv
func SessionKv(pContext *gin.Context, key string, value any) {
	session := sessions.Default(pContext)
	session.Set(key, value)
	session.Save()
}

// SessionClear 清理 session
func SessionClear(pContext *gin.Context) {
	// 解析session
	session := sessions.Default(pContext)
	// 清除session
	session.Clear()
	// 保存session数据
	session.Save()
}

// TransErr 翻译异常
func TransErr(pContext *gin.Context, err error) error {
	if errs, r := err.(validator.ValidationErrors); r {
		session := sessions.Default(pContext)
		lang := ""
		if v, r := session.Get("lang").(string); r {
			lang = v
		}
		var validationErrTrans validator.ValidationErrorsTranslations
		switch lang {
		//case com.LocaleZh:
		//	validationErrTrans = errs.Translate(zhTrans)
		case typ.LocaleEn:
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
				case typ.LocaleEn:
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

// ShouldBind 应该绑定参数
func ShouldBind(pContext *gin.Context, i any) error {
	err := pContext.ShouldBind(i)
	if err != nil {
		err = TransErr(pContext, err)
	}
	return err
}

func ConvAnyToStr(i any) string {
	if i == nil {
		return ""
	}

	if v, r := i.(error); r {
		return v.Error()
	}

	return fmt.Sprintf("%v", i)
}

func Html(pContext *gin.Context, templateName string, h gin.H, msg any) {
	if h == nil {
		h = gin.H{}
	}

	// 获取 user
	_, r := h["user"]
	if !r {
		user, _ := SessionUser(pContext)
		h["user"] = user
	}

	h["uri"] = pContext.Request.RequestURI
	h["url"] = pContext.Request.URL.Path

	// 没有消息就是最好的消息
	msgStr := ConvAnyToStr(msg)
	sessionMsg, err := SessionV[string](pContext, "msg", true)
	if err == nil {
		if msgStr != "" {
			msgStr += ", "
		}
		msgStr += sessionMsg
	}
	h["msg"] = msgStr

	pContext.HTML(http.StatusOK, templateName, h)
}

func Redirect(pContext *gin.Context, location string, h gin.H, msg any) {
	session := sessions.Default(pContext)
	if h != nil {
		for k, v := range h {
			session.Set(k, v)
		}
	}
	session.Set("msg", ConvAnyToStr(msg))
	session.Save()
	pContext.Redirect(http.StatusMovedPermanently, location)
}

func PostForm[T any](pContext *gin.Context, key string) (T, error) {
	value := pContext.PostForm(key)
	return util.ConvStrToT[T](value)
}

func Param[T any](pContext *gin.Context, key string) (T, error) {
	value := pContext.Param(key)
	return util.ConvStrToT[T](value)
}

func Query[T any](pContext *gin.Context, key string) (T, error) {
	value := pContext.Query(key)
	return util.ConvStrToT[T](value)
}

func DataDir(pContext *gin.Context) string {
	if pContext == nil {
		return arg.DataDir
	}

	user, _ := SessionUser(pContext)
	return fmt.Sprintf("%s%s%d", arg.DataDir, util.FileSeparator, user.Id)
}

func dsn(pContext *gin.Context) string {
	dataDir := DataDir(pContext)
	return fmt.Sprintf("%s%sdatabase.db", dataDir, util.FileSeparator)
}

func DbQry[T any](pContext *gin.Context, sql string, args ...any) (T, int64, error) {
	return db.Qry[T](dsn(pContext), sql, args...)
}

func DbAdd(pContext *gin.Context, sql string, args ...any) (int64, error) {
	return db.Add(dsn(pContext), sql, args...)
}

func DbUpd(pContext *gin.Context, sql string, args ...any) (int64, error) {
	return db.Upd(dsn(pContext), sql, args...)
}

func DbDel(pContext *gin.Context, sql string, args ...any) (int64, error) {
	return db.Del(dsn(pContext), sql, args...)
}

func DbPage[T any](pContext *gin.Context, req page.Req, sql string, args ...any) (page.Page[T], error) {
	return db.Page[T](dsn(pContext), req, sql, args...)
}

func PageReq(pContext *gin.Context) (page.Req, error) {
	req := page.Req{Size: 10}
	err := ShouldBind(pContext, &req)
	if req.Current <= 0 {
		req.Current = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	return req, err
}
