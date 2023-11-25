// @author xiangqian
// @date 22:40 2023/06/13
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
	"time"
)

// SignIn 登录
func SignIn(request *http.Request, session *session.Session) (string, model.Response) {
	// 请求方法
	method := request.Method

	// 登录页
	if method != http.MethodPost {
		return "system/signin", model.Response{
			Msg: session.GetMsg(),
		}
	}

	// 错误重定向到登录页
	errRedirect := func(err any) (string, model.Response) {
		return "redirect:/signin", model.Response{
			Msg: util_string.String(err),
		}
	}

	// 密码
	passwd := request.PostFormValue("passwd")

	// 校验密码
	err := util_validate.Passwd(passwd, session.GetLanguage())
	if err != nil {
		return errRedirect(err)
	}

	// 查询密钥信息表
	db := db.Get()
	result, err := db.Get("SELECT `passwd`, `try`, `last_sign_in_ip`, `last_sign_in_time`, `current_sign_in_ip`, `current_sign_in_time`, `upd_time` FROM `system` LIMIT 1")
	if err != nil {
		return errRedirect(err)
	}

	// 映射密钥信息表
	var system model.System
	err = result.Scan(&system)
	if err != nil {
		return errRedirect(err)
	}

	// 判断账号是否被锁定
	try := system.Try
	if try >= 3 {
		// 获取系统锁定时间
		lockTime := util_time.ParseUnix(system.UpdTime)
		// 系统锁定持续时间
		duration := time.Now().Sub(lockTime)
		hour := int64(duration.Hours())
		// 如果系统锁定超过24h，则自动解除锁定
		if hour > 24 {
			try = 0
			updTry(try)
		} else
		// 系统已锁定
		{
			return errRedirect(util_i18n.GetMessage("i18n.systemHasBeenLocked", session.GetLanguage()))
		}
	}

	// 校验密码
	err = util_crypto_bcrypt.CompareHash(passwd, system.Passwd)
	// 密码错误
	if err != nil {
		updTry(try + 1)
		if try == 1 {
			return errRedirect(util_i18n.GetMessage("i18n.systemWillBeLocked", session.GetLanguage()))
		} else {
			return errRedirect(util_i18n.GetMessage("i18n.passwdIncorrect", session.GetLanguage()))
		}
	}

	// 重置try值
	if try != 0 {
		updTry(0)
	}

	//	保存系统信息到会话
	session.SetSystem(system)

	// 重定向到首页
	return "redirect:index", model.Response{}
}

func updTry(try byte) {
	db := db.Get()
	db.Upd("UPDATE `System` SET `try` = ?, `upd_time` = ?", try, util_time.NowUnix())
}

// --------------------

//_method := strings.TrimSpace(request.URL.Query().Get("_method"))
//if _method == "" {
//_method = strings.TrimSpace(request.PostForm.Get("_method"))
//}
