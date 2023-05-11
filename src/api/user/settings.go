// user settings
// @author xiangqian
// @date 22:33 2023/05/11
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

// Settings0 更新用户信息
func Settings0(context *gin.Context) {
	// 更新异常时，重定向到设置页
	redirect := func(user typ.User, err any) {
		resp := typ.Resp[typ.User]{
			Msg:  str.ConvTypeToStr(err),
			Data: user,
		}
		api_common_context.Redirect(context, "/user/settings", resp)
	}

	// bind
	user := typ.User{}
	err := api_common_context.ShouldBind(context, &user)
	if err != nil {
		redirect(user, err)
		return
	}

	// validate name
	err = validate.UserName(user.Name)
	if err != nil {
		redirect(user, err)
		return
	}

	// validate passwd
	err = validate.Passwd(user.Passwd)
	if err != nil {
		redirect(user, err)
		return
	}
	err = validate.Passwd(user.OrigPasswd)
	if err != nil {
		redirect(user, err)
		return
	}

	// 密码加密
	passwd, err := bcrypt.Generate(user.Passwd)
	if err != nil {
		redirect(user, err)
		return
	}

	// name
	sessionUser, err := session.GetUser(context)
	if err != nil {
		redirect(user, err)
		return
	}
	if user.Name != sessionUser.Name {
		// 校验数据库用户名
		err = ValidateUserName(user.Name)
		if err != nil {
			redirect(user, err)
			return
		}
	}

	user.Nickname = strings.TrimSpace(user.Nickname)
	user.Rem = strings.TrimSpace(user.Rem)

	// query
	origPasswdHash, count, err := db.Qry[string](nil, "SELECT `passwd` FROM `user` WHERE `del` = 0 AND `id` = ? LIMIT 1", sessionUser.Id)
	if err != nil || count == 0 {
		redirect(user, err)
		return
	}
	if !bcrypt.CompareHash(origPasswdHash, user.OrigPasswd) {
		redirect(user, i18n.MustGetMessage("i18n.userOrPasswdIncorrect"))
		return
	}

	// update
	updTime := time.NowUnix()
	_, err = db.Upd(nil, "UPDATE `user` SET `name` = ?, nickname = ?, `passwd` = ?, rem = ?, upd_time = ? WHERE id = ?",
		user.Name, user.Nickname, passwd, user.Rem, updTime, sessionUser.Id)
	if err != nil {
		redirect(user, err)
		return
	}

	// 更新session中User信息
	sessionUser.Name = user.Name
	sessionUser.Nickname = user.Nickname
	sessionUser.Rem = user.Rem
	sessionUser.UpdTime = updTime
	session.SetUser(context, sessionUser)

	// redirect
	redirect(user, nil)
}

// Settings 用户设置页
func Settings(context *gin.Context) {
	resp, err := session.Get[typ.Resp[typ.User]](context, api_common_context.RespSessionKey, true)
	if err != nil {
		user, err := session.GetUser(context)
		resp.Msg = str.ConvTypeToStr(err)
		resp.Data = user
	}
	api_common_context.HtmlOk(context, "user/reg.html", resp)
}
