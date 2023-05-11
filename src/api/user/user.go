// user
// @author xiangqian
// @date 13:20 2023/02/04
package user

import (
	"errors"
	"github.com/gin-contrib/i18n"
	"note/src/api/common/db"
)

const TmpUserSessionKey = "tmp_user"

//// Upd 更新用户信息
//func Upd(context *gin.Context) {
//	// 更新异常时，重定向到设置页
//	redirect := func(user typ.User, err any) {
//		resp := typ.Resp[typ.User]{
//			Msg:  util_str.ConvTypeToStr(err),
//			Data: user,
//		}
//		context.Redirect(context, "/user/settings", resp)
//	}
//
//	// bind
//	user := typ.User{}
//	err := context.ShouldBind(context, &user)
//	if err != nil {
//		redirect(user, err)
//		return
//	}
//
//	// name
//	err = util_validate.UserName(user.Name)
//	if err != nil {
//		redirect(user, err)
//		return
//	}
//
//	// passwd
//	err = util_validate.Passwd(user.Passwd)
//	if err != nil {
//		redirect(user, err)
//		return
//	}
//
//	// name
//	sessionUser, err := session.GetSessionUser(context)
//	if err != nil {
//		redirect(user, err)
//		return
//	}
//	if user.Name != sessionUser.Name {
//		// 校验数据库用户名
//		err = VerifyDbName(user.Name)
//		if err != nil {
//			redirect(user, err)
//			return
//		}
//	}
//
//	user.Nickname = strings.TrimSpace(user.Nickname)
//	user.Rem = strings.TrimSpace(user.Rem)
//
//	// update
//	updTime := util_time.NowUnix()
//	_, err = db.DbUpd(nil, "UPDATE `user` SET `name` = ?, nickname = ?, `passwd` = ?, rem = ?, upd_time = ? WHERE id = ?",
//		user.Name, user.Nickname, EncryptPasswd(user.Passwd), user.Rem, updTime, sessionUser.Id)
//	if err != nil {
//		redirect(user, err)
//		return
//	}
//
//	// 更新session中User信息
//	sessionUser.Name = user.Name
//	sessionUser.Nickname = user.Nickname
//	sessionUser.Rem = user.Rem
//	sessionUser.UpdTime = updTime
//	session.SetSessionUser(context, sessionUser)
//
//	redirect(user, nil)
//}
//
//// Settings 用户设置页
//func Settings(context *gin.Context) {
//	resp, err := session.GetSessionV[typ.Resp[typ.User]](context, session.RespSessionKey, true)
//	if err != nil {
//		user, err := session.GetSessionUser(context)
//		resp.Msg = util_str.ConvTypeToStr(err)
//		resp.Data = user
//	}
//
//	context.HtmlOk(context, "user/settings.html", resp)
//}
//

// ValidateUserName 校验数据库用户名是否存在
func ValidateUserName(name string) error {
	_, count, err := db.Qry[int64](nil, "SELECT `id` FROM `user` WHERE `del` = 0 AND `name` = ? LIMIT 1", name)
	if err != nil {
		return err
	}

	if count != 0 {
		return errors.New(i18n.MustGetMessage("i18n.userNameAlreadyExists"))
	}

	return nil
}
