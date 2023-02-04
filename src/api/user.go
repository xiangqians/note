// user
// @author xiangqian
// @date 13:20 2023/02/04
package api

import (
	"github.com/gin-gonic/gin"
	"note/src/typ"
)

const SessionKeyUser = "_user_"

func SessionUser(pContext *gin.Context) (typ.User, error) {
	user, err := SessionV[typ.User](pContext, SessionKeyUser, false)

	// 如果返回指针值，有可能会发生逃逸
	//return &user

	return user, err
}

// UserRegPage 用户注册页
func UserRegPage(pContext *gin.Context) {
	user, err := SessionV[typ.User](pContext, "user", true)
	if err != nil {
		user = typ.User{}
	}
	Html(pContext, "user/reg.html", gin.H{"user": user}, nil)
}

// 用户注册
//func UserAdd(pContext *gin.Context) {
//	// 用户注册成功后，重定向到登录页
//	redirectToRegPage := func(message any) {
//		Redirect(pContext, "/user/regpage", message, nil)
//	}
//
//	setting, count, err := Setting()
//	if err != nil || count == 0 {
//		redirectToRegPage(err)
//		return
//	}
//
//	if setting.AllowReg != 1 {
//		redirectToRegPage(i18n.MustGetMessage("i18n.regNotOpen"))
//		return
//	}
//
//	user := typ.User{}
//	err = ShouldBind(pContext, &user)
//	if err == nil {
//		err = VerifyUserNameAndPasswd(user.Name, user.Passwd)
//	}
//
//	if err == nil {
//		err = VerifyDbUserName(user.Name)
//	}
//
//	session := sessions.Default(pContext)
//	if err != nil {
//		session.Set("user", user)
//		session.Set("message", err.Error())
//		session.Save()
//		pContext.Redirect(http.StatusMovedPermanently, "/user/regpage")
//		return
//	}
//
//	db.Add("INSERT INTO `user` (`name`, `nickname`, `passwd`, `rem`, `add_time`) VALUES (?, ?, ?, ?, ?)",
//		user.Name, strings.TrimSpace(user.Nickname), PasswdEncrypt(user.Passwd), strings.TrimSpace(user.Rem), time.Now().Unix())
//
//	session.Set("username", user.Name)
//	session.Set("message", i18n.MustGetMessage("i18n.accountRegSuccess"))
//	session.Save()
//
//	// 用户注册成功后，重定向到登录页
//	Redirect(pContext, "/user/loginpage", nil, nil)
//}
//

// UserLoginPage 用户登录html
func UserLoginPage(pContext *gin.Context) {
	username, _ := SessionV[string](pContext, "username", true)
	Html(pContext, "user/login.html", gin.H{"username": username}, nil)
}

//
//// 用户登录
//func UserLogin(pContext *gin.Context) {
//	name := strings.TrimSpace(pContext.PostForm("name"))
//	passwd := strings.TrimSpace(pContext.PostForm("passwd"))
//
//	err := VerifyUserNameAndPasswd(name, passwd)
//
//	var user typ.User
//	if err == nil {
//		user, _, err = db.Qry[typ.User]("SELECT u.id, u.`name`, u.nickname, u.rem, u.add_time, u.upd_time FROM `user` u WHERE u.del = 0 AND u.`name` = ? AND u.passwd = ? LIMIT 1", name, PasswdEncrypt(passwd))
//	}
//
//	if err == nil && user.Id == 0 {
//		err = errors.New(i18n.MustGetMessage("i18n.userOrPasswdIncorrect"))
//	}
//
//	// 初始化session对象
//	session := sessions.Default(pContext)
//
//	if err != nil {
//		session.Set("username", name)
//		session.Set("message", err.Error())
//		session.Save()
//		pContext.Redirect(http.StatusMovedPermanently, "/user/loginpage")
//		return
//	}
//
//	// session
//	// 设置session数据
//	session.Set(SessionKeyUser, user)
//	// 保存session数据
//	session.Save()
//
//	// 重定向
//	pContext.Redirect(http.StatusMovedPermanently, "/")
//}
//
//// 用户登出
//func UserLogout(pContext *gin.Context) {
//	// 清除session
//	SessionClear(pContext)
//
//	// 重定向
//	pContext.Redirect(http.StatusMovedPermanently, "/user/loginpage")
//}
//
//func UserSettingPage(pContext *gin.Context) {
//	session := sessions.Default(pContext)
//	user := session.Get("user")
//	message := session.Get("message")
//	session.Delete("user")
//	session.Delete("message")
//	session.Save()
//
//	if user == nil {
//		user = SessionUser(pContext)
//	}
//	pContext.HTML(http.StatusOK, "user/setting.html", gin.H{
//		"user":    user,
//		"message": message,
//	})
//}
//
//func UserUpd(pContext *gin.Context) {
//	user := typ.User{}
//	err := ShouldBind(pContext, &user)
//
//	if err == nil {
//		err = VerifyUserNameAndPasswd(user.Name, user.Passwd)
//	}
//
//	sessionUser := SessionUser(pContext)
//	if err == nil && user.Name != sessionUser.Name {
//		err = VerifyDbUserName(user.Name)
//	}
//
//	user.Nickname = strings.TrimSpace(user.Nickname)
//	user.Rem = strings.TrimSpace(user.Rem)
//
//	session := sessions.Default(pContext)
//	if err != nil {
//		session.Set("user", user)
//		session.Set("message", err.Error())
//		session.Save()
//		pContext.Redirect(http.StatusMovedPermanently, "/user/settingpage")
//		return
//	}
//
//	db.Add("UPDATE `user` SET `name` = ?, nickname = ?, `passwd` = ?, rem = ?, upd_time = ? WHERE id = ?",
//		user.Name, user.Nickname, PasswdEncrypt(user.Passwd), user.Rem, time.Now().Unix(), sessionUser.Id)
//
//	// 更新session中User信息
//	sessionUser.Name = user.Name
//	sessionUser.Nickname = user.Nickname
//	sessionUser.Rem = user.Rem
//	session.Set(SessionKeyUser, sessionUser)
//	session.Save()
//
//	pContext.Redirect(http.StatusMovedPermanently, "/user/settingpage")
//}
//
//func VerifyUserNameAndPasswd(name, passwd string) error {
//	err := util.VerifyUserName(name)
//	if err == nil {
//		err = util.VerifyPasswd(passwd)
//	}
//	return err
//}
//
//func VerifyDbUserName(name string) error {
//	id, count, err := db.Qry[int64]("SELECT u.id FROM `user` u WHERE u.del = 0 AND u.`name` = ? LIMIT 1", name)
//	if err != nil || count == 0 {
//		return err
//	}
//
//	if id != 0 {
//		return errors.New(i18n.MustGetMessage("i18n.userAlreadyExists"))
//	}
//
//	return nil
//}
//
//func PasswdEncrypt(passwd string) string {
//	d := md5.New()
//	salt := "test"
//	str := ""
//	for i := 0; i < len(passwd); i++ {
//		str += fmt.Sprintf("%c", passwd[i])
//		if i%2 == 0 {
//			str += salt
//		}
//	}
//
//	_, err := io.WriteString(d, str)
//	if err != nil {
//		log.Println(err)
//		return passwd
//	}
//
//	return hex.EncodeToString(d.Sum(nil))
//}
