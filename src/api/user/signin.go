// signin
// @author xiangqian
// @date 22:40 2023/06/13
package user

import (
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"note/src/api/common"
	"note/src/context"
	"note/src/session"
	"note/src/typ"
	util_crypto_bcrypt "note/src/util/crypto/bcrypt"
	util_time "note/src/util/time"
	util_validate "note/src/util/validate"
	"time"
)

const signInNameKey = "signInNameKey"
const signInErrKey = "signInErrKey"

// SignIn 登录页
func SignIn(ctx *gin.Context) {
	name, _ := session.Get[string](ctx, signInNameKey, true)
	err, _ := session.Get[any](ctx, signInErrKey, true)
	context.HtmlOk(ctx, "user/signin", common.Resp[typ.User](typ.User{Name: name}, err))
}

// SignIn0 登录
func SignIn0(ctx *gin.Context) {
	// 保存用户信息到session
	session.SetUser(ctx, typ.User{
		Abs:  typ.Abs{Id: 1},
		Name: "test"})

	return

	// 用户名
	name, _ := context.PostForm[string](ctx, "name")

	// 错误重定向到登录页
	errRedirect := func(err any) {
		session.Set(ctx, signInNameKey, name)
		session.Set(ctx, signInErrKey, err)
		context.Redirect(ctx, common.Arg.Path+"/user/signin")
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

	// 获取数据库
	db, err := common.Db(nil)
	if err != nil {
		errRedirect(err)
		return
	}

	// 查询用户信息
	var user typ.User
	db.Raw("SELECT `id`, `name`, `nickname`, `passwd`, `rem`, `try`, `add_time`, `upd_time` FROM `user` WHERE `del` = 0 AND `name` = ? LIMIT 1", name).Scan(&user)
	// 用户不存在
	if user.Id == 0 {
		errRedirect(i18n.MustGetMessage("i18n.userNameOrPasswdIncorrect"))
		return
	}

	// 判断账号是否被锁定
	if user.Try >= 3 {
		// 获取账号锁定时间
		lockTime := util_time.ParseUnix(user.UpdTime)
		// 账号锁定持续时间
		duration := time.Now().Sub(lockTime)
		hour := int64(duration.Hours())
		// 如果账号锁定超过24h，则自动解除锁定
		if hour > 24 {
			resetTry(db, user.Id)
		} else
		// 账号已锁定
		{
			errRedirect(i18n.MustGetMessage("i18n.accountHasBeenLocked"))
			return
		}
	}

	// 密码错误
	if util_crypto_bcrypt.CompareHash(user.Passwd, passwd) != nil {
		if user.Try == 1 {
			errRedirect(i18n.MustGetMessage("i18n.accountWillBeLocked"))
		} else {
			errRedirect(i18n.MustGetMessage("i18n.userOrPasswdIncorrect"))
		}
		db.Exec("UPDATE `user` SET `try` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ?", user.Try+1, util_time.NowUnix(), user.Id)
		return
	}

	// 重置 try 状态
	if user.Try != 0 {
		resetTry(db, user.Id)
	}

	// 清理已登录用户
	//app.ClearUser(user.Id)

	// 密码不存于session
	user.Passwd = ""

	// 保存用户信息到session
	session.SetUser(ctx, user)

	// 重定向到首页
	context.Redirect(ctx, common.Arg.Path+"/")
}

// 重置try值
func resetTry(db *gorm.DB, id int64) {
	db.Exec("UPDATE `user` SET `try` = 0, `upd_time` = ? WHERE `del` = 0 AND `id` = ?", util_time.NowUnix(), id)
}
