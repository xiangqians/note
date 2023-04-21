// user
// @author xiangqian
// @date 13:20 2023/02/04
package user

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"note/src/api/common"
	"note/src/db"
	typ_api "note/src/typ/api"
	typ_resp "note/src/typ/resp"
	util_os "note/src/util/os"
	util_str "note/src/util/str"
	util_time "note/src/util/time"
	"os"
	"regexp"
	"strings"
)

// Upd 更新用户信息
func Upd(context *gin.Context) {
	// 更新异常时，重定向到设置页
	redirect := func(user typ_api.User, err any) {
		resp := typ_resp.Resp[typ_api.User]{
			Msg:  util_str.TypeToStr(err),
			Data: user,
		}
		common.Redirect(context, "/user/settings", resp)
	}

	// bind
	user := typ_api.User{}
	err := common.ShouldBind(context, &user)
	if err != nil {
		redirect(user, err)
		return
	}

	// name
	err = VerifyName(user.Name)
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

	// name
	sessionUser, err := common.GetSessionUser(context)
	if err != nil {
		redirect(user, err)
		return
	}
	if user.Name != sessionUser.Name {
		// 校验数据库用户名
		err = VerifyDbName(user.Name)
		if err != nil {
			redirect(user, err)
			return
		}
	}

	user.Nickname = strings.TrimSpace(user.Nickname)
	user.Rem = strings.TrimSpace(user.Rem)

	// update
	updTime := util_time.NowUnix()
	_, err = common.DbUpd(nil, "UPDATE `user` SET `name` = ?, nickname = ?, `passwd` = ?, rem = ?, upd_time = ? WHERE id = ?",
		user.Name, user.Nickname, EncryptPasswd(user.Passwd), user.Rem, updTime, sessionUser.Id)
	if err != nil {
		redirect(user, err)
		return
	}

	// 更新session中User信息
	sessionUser.Name = user.Name
	sessionUser.Nickname = user.Nickname
	sessionUser.Rem = user.Rem
	sessionUser.UpdTime = updTime
	common.SetSessionUser(context, sessionUser)

	redirect(user, nil)
}

// Settings 用户设置页
func Settings(context *gin.Context) {
	resp, err := common.GetSessionV[typ_resp.Resp[typ_api.User]](context, common.RespSessionKey, true)
	if err != nil {
		user, err := common.GetSessionUser(context)
		resp.Msg = util_str.TypeToStr(err)
		resp.Data = user
	}

	common.HtmlOk(context, "user/settings.html", resp)
}

// Logout 用户登出
func Logout(context *gin.Context) {
	// 清除session
	common.ClearSession(context)

	// 重定向
	common.Redirect(context, "/user/login", typ_resp.Resp[any]{})
}

// Login0 用户登录
func Login0(context *gin.Context) {
	// name
	name := strings.TrimSpace(context.PostForm("name"))

	// redirect
	redirect := func(err any) {
		resp := typ_resp.Resp[typ_api.User]{
			Msg:  util_str.TypeToStr(err),
			Data: typ_api.User{Name: name},
		}
		common.Redirect(context, "/user/login", resp)
	}

	// verify name
	err := VerifyName(name)
	if err != nil {
		redirect(err)
		return
	}

	// passwd
	passwd := strings.TrimSpace(context.PostForm("passwd"))
	err = VerifyPasswd(passwd)
	if err != nil {
		redirect(err)
		return
	}

	// query
	user, count, err := common.DbQry[typ_api.User](nil,
		"SELECT `id`, `name`, `nickname`, `rem`, `add_time`, `upd_time` FROM `user` WHERE `del` = 0 AND `name` = ? AND `passwd` = ? LIMIT 1",
		name, EncryptPasswd(passwd))
	if err != nil {
		redirect(err)
		return
	}

	if count == 0 {
		redirect(i18n.MustGetMessage("i18n.userOrPasswdIncorrect"))
		return
	}

	// 保存用户信息到session
	common.SetSessionUser(context, user)

	// 重定向
	common.Redirect(context, "/", typ_resp.Resp[any]{})
}

// Login 用户登录页
func Login(context *gin.Context) {
	resp, _ := common.GetSessionV[typ_resp.Resp[typ_api.User]](context, common.RespSessionKey, true)
	common.HtmlOk(context, "user/login.html", resp)
}

// Add 添加用户（用户注册）
func Add(context *gin.Context) {
	// 注册异常时，重定向到注册页
	redirect := func(user typ_api.User, err any) {
		resp := typ_resp.Resp[typ_api.User]{
			Msg:  util_str.TypeToStr(err),
			Data: user,
		}
		common.Redirect(context, "/user/reg", resp)
	}

	// 是否允许用户注册
	if common.AppArg.AllowReg != 1 {
		redirect(typ_api.User{}, i18n.MustGetMessage("i18n.regNotOpen"))
		return
	}

	// bind
	user := typ_api.User{}
	err := common.ShouldBind(context, &user)
	if err != nil {
		redirect(user, err)
		return
	}

	// name
	err = VerifyName(user.Name)
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
	err = VerifyDbName(user.Name)
	if err != nil {
		redirect(user, err)
		return
	}

	// db
	// get
	_db := db.Get(common.Dsn(nil))
	defer _db.Close()
	// open
	err = _db.Open()
	if err != nil {
		redirect(user, err)
		return
	}
	// begin
	err = _db.Begin()
	if err != nil {
		redirect(user, err)
		return
	}

	// add
	id, err := _db.Add("INSERT INTO `user` (`name`, `nickname`, `passwd`, `rem`, `add_time`) VALUES (?, ?, ?, ?, ?)",
		user.Name, strings.TrimSpace(user.Nickname), EncryptPasswd(user.Passwd), strings.TrimSpace(user.Rem), util_time.NowUnix())
	if err != nil {
		_db.Rollback()
		redirect(user, err)
		return
	}

	// 创建用户数据目录
	dataDir := fmt.Sprintf("%s%s%d", common.AppArg.DataDir, util_os.FileSeparator(), id)
	if !util_os.IsExist(dataDir) {
		util_os.MkDir(dataDir)
	}
	log.Printf("dataDir: %v\n", dataDir)

	// 复制文件
	// src
	src, err := os.Open(fmt.Sprintf("%s%s%s%s%s", common.AppArg.DataDir, util_os.FileSeparator(), "{id}", util_os.FileSeparator(), "database.db"))
	if err != nil {
		_db.Rollback()
		redirect(user, err)
		return
	}
	defer src.Close()
	// dst
	dst, err := os.Create(fmt.Sprintf("%s%s%s", dataDir, util_os.FileSeparator(), "database.db"))
	if err != nil {
		_db.Rollback()
		redirect(user, err)
		return
	}
	defer dst.Close()
	// copy
	err = util_os.CopyIo(src, dst, 0)
	if err != nil {
		_db.Rollback()
		redirect(user, err)
		return
	}

	// db commit
	_db.Commit()

	// 用户注册成功后，重定向到登录页
	resp := typ_resp.Resp[typ_api.User]{
		Msg:  i18n.MustGetMessage("i18n.accountRegSuccess"),
		Data: user,
	}
	common.Redirect(context, "/user/login", resp)
}

// Reg 用户注册页
func Reg(context *gin.Context) {
	resp, _ := common.GetSessionV[typ_resp.Resp[typ_api.User]](context, common.RespSessionKey, true)
	common.HtmlOk(context, "user/reg.html", resp)
}

func VerifyDbName(name string) error {
	_, count, err := common.DbQry[int64](nil, "SELECT `id` FROM `user` WHERE `del` = 0 AND `name` = ? LIMIT 1", name)
	if err != nil {
		return err
	}

	if count != 0 {
		return errors.New(i18n.MustGetMessage("i18n.userNameAlreadyExists"))
	}

	return nil
}

// EncryptPasswd 加密密码
func EncryptPasswd(passwd string) string {
	d := md5.New()
	salt := "123456"
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

// VerifyName 校验用户名
// 1-16位长度（字母，数字，下划线，减号）
func VerifyName(username string) error {
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
