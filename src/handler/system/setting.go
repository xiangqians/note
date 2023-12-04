// @author xiangqian
// @date 22:33 2023/05/11
package system

import (
	"net/http"
	"note/src/db"
	"note/src/model"
	"note/src/session"
	util_crypto_bcrypt "note/src/util/crypto/bcrypt"
	util_i18n "note/src/util/i18n"
	util_string "note/src/util/string"
	util_time "note/src/util/time"
	util_validate "note/src/util/validate"
)

// Setting 设置
func Setting(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
	// 请求方法
	method := request.Method

	// 密码页
	if method != http.MethodPost {
		return "system/setting", model.Response{
			Msg: session.GetMsg(),
		}
	}

	// 错误重定向到密码页
	errRedirect := func(err any) (string, model.Response) {
		return "redirect:/setting", model.Response{
			Msg: util_string.String(err),
		}
	}

	// 原密码
	origPasswd := request.PostFormValue("origPasswd")
	// 校验原密码
	err := util_validate.Passwd(origPasswd, session.GetLanguage())
	if err != nil {
		return errRedirect(util_i18n.GetMessage("i18n.origPasswdIncorrect", session.GetLanguage()))
	}

	// 新密码
	newPasswd := request.PostFormValue("newPasswd")
	// 再次输入新密码
	reNewPasswd := request.PostFormValue("reNewPasswd")

	if newPasswd != reNewPasswd {
		return errRedirect(util_i18n.GetMessage("i18n.newPasswdEnteredTwiceInconsistent", session.GetLanguage()))
	}

	// 校验新密码
	err = util_validate.Passwd(newPasswd, session.GetLanguage())
	if err != nil {
		return errRedirect(err)
	}

	// 获取系统信息
	system, err := getSystem()
	if err != nil {
		return errRedirect(err)
	}

	// 校验原密码
	err = util_crypto_bcrypt.CompareHash(origPasswd, system.Passwd)
	// 原密码错误
	if err != nil {
		return errRedirect(util_i18n.GetMessage("i18n.origPasswdIncorrect", session.GetLanguage()))
	}

	newHash, err := util_crypto_bcrypt.Generate(newPasswd)
	if err != nil {
		return errRedirect(err)
	}

	db := db.Get()
	db.Upd("UPDATE `system` SET `passwd` = ?, `upd_time` = ?", newHash, util_time.NowUnix())

	// 修改密码
	return "redirect:/setting", model.Response{}
}
