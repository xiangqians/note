// 用户登录
// @author xiangqian
// @date 22:40 2023/06/13
package user

import (
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"note/src/api"
	"note/src/context"
	"note/src/session"
	"note/src/typ"
	util_crypto_bcrypt "note/src/util/crypto/bcrypt"
	util_string "note/src/util/string"
	util_time "note/src/util/time"
	util_validate "note/src/util/validate"
	"time"
)

const signInNameKey = "signInName"
const signInErrKey = "signInErr"

// SignIn 登录
func SignIn(ctx *gin.Context) {
	method, _ := context.PostForm[string](ctx, "_method")
	if method == "POST" {
		signIn(ctx)
	} else {
		name, _ := session.Get[string](ctx, signInNameKey, true)
		err, _ := session.Get[string](ctx, signInErrKey, true)
		context.HtmlOk(ctx, "user/signin", typ.Resp[typ.User]{Data: typ.User{Name: name}, Msg: err})
	}
}

func signIn(ctx *gin.Context) {
	// 用户名
	name, _ := context.PostForm[string](ctx, "name")

	// 错误重定向到登录页
	errRedirect := func(err any) {
		session.Set(ctx, signInNameKey, name)
		session.Set(ctx, signInErrKey, util_string.String(err))
		context.Redirect(ctx, "/user/signin")
	}

	// 校验用户名
	err := util_validate.UserName(name)
	if err != nil {
		errRedirect(err)
		return
	}

	// 密码
	passwd, _ := context.PostForm[string](ctx, "passwd")
	// 校验密码
	err = util_validate.Passwd(passwd)
	if err != nil {
		errRedirect(err)
		return
	}

	// 获取数据库操作实例
	db, err := api.Db(nil)
	if err != nil {
		errRedirect(err)
		return
	}

	// 根据用户名查询用户信息
	user := getByName(db, name)
	// 用户不存在
	if user.Id == 0 {
		errRedirect(i18n.MustGetMessage("i18n.userNameOrPasswdIncorrect"))
		return
	}

	// 判断账号是否被锁定
	try := user.Try
	if try >= 3 {
		// 获取账号锁定时间
		lockTime := util_time.ParseUnix(user.UpdTime)
		// 账号锁定持续时间
		duration := time.Now().Sub(lockTime)
		hour := int64(duration.Hours())
		// 如果账号锁定超过24h，则自动解除锁定
		if hour > 24 {
			try = 0
			updTryById(db, user.Id, try)
		} else
		// 账号已锁定
		{
			errRedirect(i18n.MustGetMessage("i18n.accountHasBeenLocked"))
			return
		}
	}

	// 密码错误
	if util_crypto_bcrypt.CompareHash(user.Passwd, passwd) != nil {
		updTryById(db, user.Id, try+1)
		if try == 1 {
			errRedirect(i18n.MustGetMessage("i18n.accountWillBeLocked"))
		} else {
			errRedirect(i18n.MustGetMessage("i18n.userNameOrPasswdIncorrect"))
		}
		return
	}

	// 重置try值
	if try != 0 {
		updTryById(db, user.Id, 0)
	}

	// 密码不存于session
	user.Passwd = ""

	// 保存用户信息到session
	session.SetUser(ctx, user)

	// 重定向到首页
	context.Redirect(ctx, "/")
}

// 根据用户id更新try值
func updTryById(db *gorm.DB, id int64, try byte) {
	db.Exec("UPDATE `user` SET `try` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ?", try, util_time.NowUnix(), id)
}

// 根据用户名查询用户信息
func getByName(db *gorm.DB, name string) typ.User {
	var user typ.User
	db.Raw("SELECT `id`, `name`, `nickname`, `passwd`, `rem`, `try`, `add_time`, `upd_time` FROM `user` WHERE `del` = 0 AND `name` = ? LIMIT 1", name).Scan(&user)
	return user
}
