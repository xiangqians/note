// signin
// @author xiangqian
// @date 22:40 2023/06/13
package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"note/src/api"
	"note/src/app"
	"note/src/context"
	"note/src/session"
	"note/src/typ"
	util_validate "note/src/util/validate"
)

const signInNameKey = "signInNameKey"
const signInErrKey = "signInErrKey"

// SignIn 登录页
func SignIn(ctx *gin.Context) {
	name, _ := session.Get[string](ctx, signInNameKey, true)
	err, _ := session.Get[any](ctx, signInErrKey, true)
	context.HtmlOk(ctx, "user/signin", api.Resp[typ.User](typ.User{Name: name}, err))
}

// SignIn0 登录
func SignIn0(ctx *gin.Context) {
	// 用户名
	name, _ := context.PostForm[string](ctx, "name")

	// 错误重定向到登录页
	errRedirect := func(err any) {
		session.Set(ctx, signInNameKey, name)
		session.Set(ctx, signInErrKey, err)
		context.Redirect(ctx, app.GetArg().Path+"/user/signin")
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
	db, err := api.Db(nil)
	if err != nil {
		errRedirect(err)
		return
	}

	fmt.Println(db)

	//var typ.User user
	//db.Raw("SELECT `id`, `name`, `nickname`, `passwd`, `rem`, `try`, `add_time`, `upd_time` FROM `user` WHERE `del` = 0 AND `name` = ? LIMIT 1")
	//
	//// query
	//user, count, err := db.Qry[](nil,
	//	, name)
	//if err != nil {
	//	redirect(err)
	//	return
	//}
	//
	//// 校验用户信息是否存在
	//if count == 0 {
	//	redirect(i18n.MustGetMessage("i18n.userOrPasswdIncorrect"))
	//	return
	//}
	//
	//// 重置try
	//resetTry := func() {
	//	db.Upd(nil, "UPDATE `user` SET `try` = 0, `upd_time` = ? WHERE `del` = 0 AND `id` = ?", time.NowUnix(), user.Id)
	//	user.Try = 0
	//}
	//
	//// lock ?
	//if user.Try >= 3 {
	//	// lock time
	//	lockTime := time.ParseUnix(user.UpdTime)
	//	// Duration
	//	duration := time.Now().Sub(lockTime)
	//	// hour
	//	hour := int64(duration.Hours())
	//	// 如果超过24h则自动解除锁定
	//	if hour > 24 {
	//		resetTry()
	//	} else {
	//		redirect(i18n.MustGetMessage("i18n.accountHasBeenLocked"))
	//		return
	//	}
	//}
	//
	//// 密码是否正确
	//if !bcrypt.CompareHash(user.Passwd, passwd) {
	//	if user.Try == 1 {
	//		redirect(i18n.MustGetMessage("i18n.accountWillBeLocked"))
	//	} else {
	//		redirect(i18n.MustGetMessage("i18n.userOrPasswdIncorrect"))
	//	}
	//	db.Upd(nil, "UPDATE `user` SET `try` = ?, `upd_time` = ? WHERE `del` = 0 AND `id` = ?", user.Try+1, time.NowUnix(), user.Id)
	//	return
	//}
	//
	//// 重置 try 状态
	//if user.Try != 0 {
	//	resetTry()
	//}
	//
	//// 密码不存于session
	//user.Passwd = ""
	//
	//// 保存用户信息到session
	//session.SetUser(context, user)
	//
	//// 重定向到首页
	//api_common_context.Redirect(context, "/", typ.Resp[any]{})

	////////////////

	app.ClearUser(1)

	// 保存用户信息到session
	session.SetUser(ctx, typ.User{Abs: typ.Abs{Id: 1}, Name: "test", Nickname: "测试"})

	// 重定向到首页
	context.Redirect(ctx, app.GetArg().Path+"/")
}
