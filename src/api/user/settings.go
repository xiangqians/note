// 用户设置
// @author xiangqian
// @date 22:33 2023/05/11
package user

import (
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"note/src/api"
	"note/src/context"
	"note/src/session"
	"note/src/typ"
	util_crypto_bcrypt "note/src/util/crypto/bcrypt"
	util_string "note/src/util/string"
	"note/src/util/time"
	"note/src/util/validate"
	"strings"
)

const settingsUserKey = "settingsUser"
const settingsErrKey = "settingsErr"

// Settings 用户设置
func Settings(ctx *gin.Context) {
	// 请求方法
	method, _ := context.PostForm[string](ctx, "_method")
	// 设置
	if method == "PUT" {
		settings(ctx)
	} else
	// 设置页
	{
		user, _ := session.Get[typ.UpdUser](ctx, settingsUserKey, true)
		sessionUser, _ := session.GetUser(ctx)
		if user.Name == "" {
			user.Name = sessionUser.Name
		}
		if user.Nickname == "" {
			user.Nickname = sessionUser.Nickname
		}
		if user.Rem == "" {
			user.Rem = sessionUser.Rem
		}
		err, _ := session.Get[string](ctx, settingsErrKey, true)
		context.HtmlOk(ctx, "user/settings", typ.Resp[typ.UpdUser]{Data: user, Msg: err})
	}
}

// 用户设置
func settings(ctx *gin.Context) {
	// 错误重定向到设置页
	errRedirect := func(user typ.UpdUser, err any) {
		user.OrigPasswd = ""
		user.Passwd = ""
		user.NewPasswd = ""
		user.ReNewPasswd = ""
		session.Set(ctx, settingsUserKey, user)
		session.Set(ctx, settingsErrKey, util_string.String(err))
		context.Redirect(ctx, "/user/settings")
	}

	user := typ.UpdUser{}

	// 绑定
	err := context.ShouldBind(ctx, &user)
	if err != nil {
		errRedirect(user, err)
		return
	}

	user.Name = strings.TrimSpace(user.Name)
	user.Nickname = strings.TrimSpace(user.Nickname)
	user.OrigPasswd = strings.TrimSpace(user.OrigPasswd)
	user.NewPasswd = strings.TrimSpace(user.NewPasswd)
	user.Rem = strings.TrimSpace(user.Rem)

	// 校验用户名
	err = validate.UserName(user.Name)
	if err != nil {
		errRedirect(user, err)
		return
	}

	// 校验密码
	err = validate.Passwd(user.NewPasswd)
	if err != nil {
		errRedirect(user, err)
		return
	}

	// 校验原密码
	sessionUser, _ := session.GetUser(ctx)
	if util_crypto_bcrypt.CompareHash(user.OrigPasswd, sessionUser.Passwd) != nil {
		errRedirect(user, i18n.MustGetMessage("i18n.origPasswdIncorrect"))
		return
	}

	// 加密新密码
	newPasswdHash, err := util_crypto_bcrypt.Generate(user.NewPasswd)
	if err != nil {
		errRedirect(user, err)
		return
	}

	// 获取数据库操作实例
	db, err := api.Db(nil)
	if err != nil {
		errRedirect(user, err)
		return
	}

	// 根据用户名查询用户信息
	dbUser := getByName(db, user.Name)
	// 校验数据库用户名
	if dbUser.Id != 0 && dbUser.Id != sessionUser.Id {
		errRedirect(user, i18n.MustGetMessage("i18n.userNameAlreadyExists"))
		return
	}

	// 更新用户信息
	updTime := time.NowUnix()
	err = db.Exec("UPDATE `user` SET `name` = ?, nickname = ?, `passwd` = ?, rem = ?, upd_time = ? WHERE `del` = 0 AND `id` = ?",
		user.Name, user.Nickname, newPasswdHash, user.Rem, updTime, sessionUser.Id).Error
	if err != nil {
		errRedirect(user, err)
		return
	}

	sessionUser.Name = user.Name
	sessionUser.Nickname = user.Nickname
	sessionUser.Passwd = newPasswdHash
	sessionUser.Rem = user.Rem
	sessionUser.UpdTime = updTime

	// 保存用户信息到session
	session.SetUser(ctx, sessionUser)

	// 重定向到首页
	context.Redirect(ctx, "/user/settings")
}
