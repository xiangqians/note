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
	"strings"
	"time"
)

// SignIn 登录
func SignIn(request *http.Request, writer http.ResponseWriter, session *session.Session) (string, model.Response) {
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

	// 获取系统信息
	system, err := getSystem()
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

	// 获取 X-Forwarded-For 头信息
	remoteAddr := request.Header.Get("X-Forwarded-For")

	// 如果 X-Forwarded-For 头信息为空，则尝试获取 X-Real-IP 头信息
	if remoteAddr == "" {
		remoteAddr = request.Header.Get("X-Real-IP")
	}

	// 如果仍然无法获取到客户端实际IP，则使用 RemoteAddr 字段的值作为备选方案
	if remoteAddr == "" {
		remoteAddr = request.RemoteAddr
		index := strings.LastIndex(remoteAddr, ":")
		if index > 0 {
			remoteAddr = remoteAddr[:index]
		}
		if strings.HasPrefix(remoteAddr, "[") && strings.HasSuffix(remoteAddr, "]") {
			remoteAddr = remoteAddr[1 : len(remoteAddr)-1]
		}
	}

	nowUnix := util_time.NowUnix()
	db := db.Get()
	db.Upd("UPDATE `system` SET `last_sign_in_ip` = `current_sign_in_ip`, `last_sign_in_time` = `current_sign_in_time`")
	db.Upd("UPDATE `system` SET `try` = ?, `current_sign_in_ip` = ?, `current_sign_in_time` = ?, `upd_time` = ?", 0, remoteAddr, nowUnix, nowUnix)

	system.Try = 0
	system.LastSignInIp = system.CurrentSignInIp
	system.LastSignInTime = system.CurrentSignInTime
	system.CurrentSignInIp = remoteAddr
	system.CurrentSignInTime = nowUnix
	system.UpdTime = nowUnix

	//	保存系统信息到会话
	session.SetSystem(system)

	// 重定向到首页
	return "redirect:/", model.Response{}
}

func updTry(try byte) error {
	db := db.Get()
	_, err := db.Upd("UPDATE `system` SET `try` = ?, `upd_time` = ?", try, util_time.NowUnix())
	return err
}

func getSystem() (model.System, error) {
	var system model.System

	db := db.Get()
	result, err := db.Get("SELECT `passwd`, `try`, `last_sign_in_ip`, `last_sign_in_time`, `current_sign_in_ip`, `current_sign_in_time`, `upd_time` FROM `system` LIMIT 1")
	if err != nil {
		return system, err
	}

	err = result.Scan(&system)
	return system, err
}
