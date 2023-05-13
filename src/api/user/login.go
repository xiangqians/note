// user login
// @author xiangqian
// @date 20:11 2023/05/11
package user

import (
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	api_common_context "note/src/api/common/context"
	"note/src/api/common/db"
	"note/src/api/common/session"
	"note/src/typ"
	"note/src/util/crypto/bcrypt"
	"note/src/util/str"
	"note/src/util/time"
	"note/src/util/validate"
	"strings"
)

// Login0 用户登录
func Login0(context *gin.Context) {
	// name
	name, _ := api_common_context.PostForm[string](context, "name")
	name = strings.TrimSpace(name)

	// redirect
	redirect := func(err any) {
		resp := typ.Resp[typ.User]{
			Msg:  str.ConvTypeToStr(err),
			Data: typ.User{Name: name},
		}
		api_common_context.Redirect(context, "/user/login", resp)
	}

	// validate name
	err := validate.UserName(name)
	if err != nil {
		redirect(err)
		return
	}

	// passwd
	passwd, _ := api_common_context.PostForm[string](context, "passwd")
	passwd = strings.TrimSpace(passwd)
	// validate passwd
	err = validate.Passwd(passwd)
	if err != nil {
		redirect(err)
		return
	}

	// query
	user, count, err := db.Qry[typ.User](nil, "SELECT `id`, `name`, `nickname`, `passwd`, `rem`, `try`, `add_time`, `upd_time` FROM `user` WHERE `del` = 0 AND `name` = ? LIMIT 1", name)
	if err != nil {
		redirect(err)
		return
	}

	// 校验用户信息是否存在
	if count == 0 {
		redirect(i18n.MustGetMessage("i18n.userOrPasswdIncorrect"))
		return
	}

	// 重置try
	resetTry := func() {
		db.Upd(nil, "UPDATE `user` SET `try` = 0, `upd_time` = ? WHERE `del` = 0 AND `id` = ?", time.NowUnix(), user.Id)
		user.Try = 0
	}

	// lock ?
	if user.Try >= 3 {
		// lock time
		lockTime := time.ParseUnix(user.UpdTime)
		// Duration
		duration := time.Now().Sub(lockTime)
		// hour
		hour := int64(duration.Hours())
		// 如果超过24h则自动解除锁定
		if hour > 24 {
			resetTry()
		} else {
			redirect(i18n.MustGetMessage("i18n.accountHasBeenLocked"))
			return
		}
	}

	// 密码是否正确
	if !bcrypt.CompareHash(user.Passwd, passwd) {
		if user.Try == 1 {
			redirect(i18n.MustGetMessage("i18n.accountWillBeLocked"))
		} else {
			redirect(i18n.MustGetMessage("i18n.userOrPasswdIncorrect"))
		}
		db.Upd(nil, "UPDATE `user` SET `try` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ?", user.Try+1, time.NowUnix(), user.Id)
		return
	}

	// 重置 try 状态
	if user.Try != 0 {
		resetTry()
	}

	// 密码不存于session
	user.Passwd = ""

	// 保存用户信息到session
	session.SetUser(context, user)

	// 重定向到首页
	api_common_context.Redirect(context, "/", typ.Resp[any]{})
}

// Login 用户登录页
func Login(context *gin.Context) {
	resp, _ := session.Get[typ.Resp[typ.User]](context, api_common_context.RespSessionKey, true)
	api_common_context.HtmlOk(context, "user/login.html", resp)
}
