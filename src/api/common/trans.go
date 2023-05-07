// trans
// @author xiangqian
// @date 20:02 2023/03/22
package common

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
	"note/src/typ"
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

// TransErr 翻译异常
func TransErr(context *gin.Context, err error) error {
	if errs, r := err.(validator.ValidationErrors); r {
		session := Session(context)
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
