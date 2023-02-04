// api
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
	"note/src/typ"
	"strings"
)

var (
	zhTrans ut.Translator
	enTrans ut.Translator
)

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

func ShouldBind(pContext *gin.Context, i any) error {
	err := pContext.ShouldBind(i)
	if err != nil {
		err = TransErr(pContext, err)
	}
	return err
}

func Html(pContext *gin.Context, templateName string, h gin.H, err error) {
	message, _ := SessionV[string](pContext, "message", true)

	if h == nil {
		h = gin.H{}
	}

	if err != nil {
		if message != "" {
			message += ", "
		}
		message += err.Error()
	}

	_, r := h["user"]
	if !r {
		user, _ := SessionUser(pContext)
		h["user"] = user
	}

	// 没有消息就是最好的消息
	h["message"] = message

	pContext.HTML(http.StatusOK, templateName, h)
}

func Redirect(pContext *gin.Context, location string, message any, m map[string]any) {
	if message != nil {
		if v, r := message.(error); r {
			message = v.Error()
		}
		if m == nil {
			m = map[string]any{}
		}
		m["message"] = message
	}

	if m != nil {
		session := sessions.Default(pContext)
		for k, v := range m {
			session.Set(k, v)
		}
		session.Save()
	}

	pContext.Redirect(http.StatusMovedPermanently, location)
}

func dsn(pContext *gin.Context) string {
	if pContext == nil {
		return fmt.Sprintf("%s/database.db", arg.DataDir)
	}

	user, _ := SessionUser(pContext)
	return fmt.Sprintf("%s/%v/database.db", arg.DataDir, user.Id)
}

func DbQry[T any](pContext *gin.Context, sql string, args ...any) (T, int64, error) {
	return db.Qry[T](dsn(pContext), sql, args...)
}

func DbAdd(pContext *gin.Context, sql string, args ...any) (int64, error) {
	return db.Add(dsn(pContext), sql, args...)
}
