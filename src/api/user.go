// user
// @author xiangqian
// @date 13:20 2023/02/04
package api

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"note/src/arg"
	"note/src/typ"
	"note/src/util"
	"regexp"
	"strings"
	"time"
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

// UserAdd 添加用户
func UserAdd(pContext *gin.Context) {
	// 注册异常时，重定向到注册页
	redirect := func(user typ.User, message any) {
		Redirect(pContext, "/user/regpage", message, map[string]any{"user": user})
	}

	if arg.AllowReg != 1 {
		redirect(typ.User{}, i18n.MustGetMessage("i18n.regNotOpen"))
		return
	}

	// bind
	user := typ.User{}
	err := ShouldBind(pContext, &user)
	if err != nil {
		redirect(user, err)
		return
	}

	// name
	err = VerifyUserName(user.Name)
	if err != nil {
		redirect(user, err)
		return
	}

	// passwd
	err = VerifyPasswd(user.Passwd)
	if err != nil {
		redirect(user, err)
		return
	}

	// 校验数据库用户名
	err = VerifyDbUserName(user.Name)
	if err != nil {
		redirect(user, err)
		return
	}

	id, err := DbAdd(nil, "INSERT INTO `user` (`name`, `nickname`, `passwd`, `rem`, `add_time`) VALUES (?, ?, ?, ?, ?)",
		user.Name, strings.TrimSpace(user.Nickname), PasswdEncrypt(user.Passwd), strings.TrimSpace(user.Rem), time.Now().Unix())
	if err != nil {
		redirect(user, err)
		return
	}

	// 创建用户数据目录
	userDataDir := fmt.Sprintf("%s%s%d", arg.DataDir, util.FileSeparator, id)
	if !util.IsExistOfPath(userDataDir) {
		util.Mkdir(userDataDir)
	}
	// 初始化用户数据目录
	pCmd := util.CopyDir(fmt.Sprintf("%s%sid", arg.DataDir, util.FileSeparator), userDataDir)
	buf, err := pCmd.CombinedOutput()
	if err != nil {
		log.Println(err)
	}
	log.Println(util.DecodeBuf(buf))

	// 用户注册成功后，重定向到登录页
	Redirect(pContext, "/user/loginpage", i18n.MustGetMessage("i18n.accountRegSuccess"), map[string]any{"userName": user.Name})
}

// UserLoginPage 用户登录页
func UserLoginPage(pContext *gin.Context) {
	userName, _ := SessionV[string](pContext, "userName", true)
	Html(pContext, "user/login.html", gin.H{"userName": userName}, nil)
}

// UserLogin 用户登录
func UserLogin(pContext *gin.Context) {
	// name
	name := strings.TrimSpace(pContext.PostForm("name"))

	redirect := func(message any) {
		Redirect(pContext, "/user/loginpage", message, map[string]any{"userName": name})
	}

	err := VerifyUserName(name)
	if err != nil {
		redirect(err)
		return
	}

	// passwd
	passwd := strings.TrimSpace(pContext.PostForm("passwd"))
	err = VerifyPasswd(passwd)
	if err != nil {
		redirect(err)
		return
	}

	user, count, err := DbQry[typ.User](nil,
		"SELECT u.id, u.`name`, u.nickname, u.rem, u.add_time, u.upd_time FROM `user` u WHERE u.del = 0 AND u.`name` = ? AND u.passwd = ? LIMIT 1",
		name, PasswdEncrypt(passwd))
	if err != nil {
		redirect(err)
		return
	}

	if count == 0 {
		err = errors.New(i18n.MustGetMessage("i18n.userOrPasswdIncorrect"))
		redirect(err)
		return
	}

	// 保存用户信息到session
	SessionKv(pContext, SessionKeyUser, user)

	// 重定向
	Redirect(pContext, "/dir/list", nil, nil)
}

// // 用户登出
//
//	func UserLogout(pContext *gin.Context) {
//		// 清除session
//		SessionClear(pContext)
//
//		// 重定向
//		pContext.Redirect(http.StatusMovedPermanently, "/user/loginpage")
//	}
//
//	func UserSettingPage(pContext *gin.Context) {
//		session := sessions.Default(pContext)
//		user := session.Get("user")
//		message := session.Get("message")
//		session.Delete("user")
//		session.Delete("message")
//		session.Save()
//
//		if user == nil {
//			user = SessionUser(pContext)
//		}
//		pContext.HTML(http.StatusOK, "user/setting.html", gin.H{
//			"user":    user,
//			"message": message,
//		})
//	}
//
//	func UserUpd(pContext *gin.Context) {
//		user := typ.User{}
//		err := ShouldBind(pContext, &user)
//
//		if err == nil {
//			err = VerifyUserNameAndPasswd(user.Name, user.Passwd)
//		}
//
//		sessionUser := SessionUser(pContext)
//		if err == nil && user.Name != sessionUser.Name {
//			err = VerifyDbUserName(user.Name)
//		}
//
//		user.Nickname = strings.TrimSpace(user.Nickname)
//		user.Rem = strings.TrimSpace(user.Rem)
//
//		session := sessions.Default(pContext)
//		if err != nil {
//			session.Set("user", user)
//			session.Set("message", err.Error())
//			session.Save()
//			pContext.Redirect(http.StatusMovedPermanently, "/user/settingpage")
//			return
//		}
//
//		db.Add("UPDATE `user` SET `name` = ?, nickname = ?, `passwd` = ?, rem = ?, upd_time = ? WHERE id = ?",
//			user.Name, user.Nickname, PasswdEncrypt(user.Passwd), user.Rem, time.Now().Unix(), sessionUser.Id)
//
//		// 更新session中User信息
//		sessionUser.Name = user.Name
//		sessionUser.Nickname = user.Nickname
//		sessionUser.Rem = user.Rem
//		session.Set(SessionKeyUser, sessionUser)
//		session.Save()
//
//		pContext.Redirect(http.StatusMovedPermanently, "/user/settingpage")
//	}
//
//	func VerifyUserNameAndPasswd(name, passwd string) error {
//		err := util.VerifyUserName(name)
//		if err == nil {
//			err = util.VerifyPasswd(passwd)
//		}
//		return err
//	}
//
//	func VerifyDbUserName(name string) error {
//		id, count, err := db.Qry[int64]("SELECT u.id FROM `user` u WHERE u.del = 0 AND u.`name` = ? LIMIT 1", name)
//		if err != nil || count == 0 {
//			return err
//		}
//
//		if id != 0 {
//			return errors.New(i18n.MustGetMessage("i18n.userAlreadyExists"))
//		}
//
//		return nil
//	}

func VerifyDbUserName(name string) error {
	_, count, err := DbQry[int64](nil, "SELECT u.id FROM `user` u WHERE u.del = 0 AND u.`name` = ? LIMIT 1", name)
	if err != nil {
		return err
	}

	if count != 0 {
		return errors.New(i18n.MustGetMessage("i18n.userNameAlreadyExists"))
	}

	return nil
}

func PasswdEncrypt(passwd string) string {
	d := md5.New()
	salt := "test"
	str := ""
	for i := 0; i < len(passwd); i++ {
		str += fmt.Sprintf("%c", passwd[i])
		if i%2 == 0 {
			str += salt
		}
	}

	_, err := io.WriteString(d, str)
	if err != nil {
		log.Println(err)
		return passwd
	}

	return hex.EncodeToString(d.Sum(nil))
}

// VerifyUserName 校验用户名
// 1-16位长度（字母，数字，下划线，减号）
func VerifyUserName(username string) error {
	if username == "" {
		return errors.New(fmt.Sprintf(i18n.MustGetMessage("i18n.xCannotEmpty"), i18n.MustGetMessage("i18n.userName")))
	}

	matched, err := regexp.MatchString("^[a-zA-Z0-9_-]{1,16}$", username)
	if err == nil && matched {
		return nil
	}

	return errors.New(fmt.Sprintf(i18n.MustGetMessage("i18n.xMastNBitsLong"), i18n.MustGetMessage("i18n.userName")))
}

// VerifyPasswd 校验密码
// 1-16位长度（字母，数字，特殊字符）
func VerifyPasswd(passwd string) error {
	matched, err := regexp.MatchString("^[a-zA-Z0-9!@#$%^&*()-_=+]{1,16}$", passwd)
	if err == nil && matched {
		return nil
	}

	return errors.New(fmt.Sprintf(i18n.MustGetMessage("i18n.xMastNBitsLong"), i18n.MustGetMessage("i18n.passwd")))
}
