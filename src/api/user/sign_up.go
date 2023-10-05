// signup
// @author xiangqian
// @date 23:33 2023/07/10
package user

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"log"
	"note/src/api"
	"note/src/context"
	"note/src/session"
	"note/src/typ"
	util_crypto_bcrypt "note/src/util/crypto/bcrypt"
	util_os "note/src/util/os"
	util_string "note/src/util/string"
	"note/src/util/time"
	"note/src/util/validate"
	"strings"
)

const signUpUserKey = "signUpUser"
const signUpErrKey = "signUpErr"

// SignUp 注册
func SignUp(ctx *gin.Context) {
	method, _ := context.PostForm[string](ctx, "_method")
	if method == "PUT" {
		signUp(ctx)
	} else {
		user, _ := session.Get[typ.User](ctx, signUpUserKey, true)
		err, _ := session.Get[any](ctx, signUpErrKey, true)
		context.HtmlOk(ctx, "user/sign_up", typ.Resp[typ.User]{Data: user, Msg: util_string.String(err)})
	}
}

func signUp(ctx *gin.Context) {
	// 错误重定向到注册页
	errRedirect := func(user typ.User, err any) {
		user.OrigPasswd = ""
		user.Passwd = ""
		user.RePasswd = ""
		session.Set(ctx, signUpUserKey, user)
		session.Set(ctx, signUpErrKey, err)
		context.Redirect(ctx, "/user/sign_up")
	}

	user := typ.User{}

	// 是否允许用户注册
	if !typ.GetArg().AllowReg {
		errRedirect(user, i18n.MustGetMessage("i18n.signUpNotOpen"))
		return
	}

	// 绑定
	err := context.ShouldBind(ctx, &user)
	if err != nil {
		errRedirect(user, err)
		return
	}

	// 用户名
	name := strings.TrimSpace(user.Name)
	// 校验用户名
	err = validate.UserName(name)
	if err != nil {
		errRedirect(user, err)
		return
	}

	// 密码
	passwd := strings.TrimSpace(user.Passwd)
	// 校验密码
	err = validate.Passwd(passwd)
	if err != nil {
		errRedirect(user, err)
		return
	}

	// 密码加密
	passwdHash, err := util_crypto_bcrypt.Generate(passwd)
	if err != nil {
		errRedirect(user, err)
		return
	}

	// 获取数据库
	db, err := api.Db(nil)
	if err != nil {
		errRedirect(user, err)
		return
	}

	// 根据用户名查询用户信息
	dbUser := getByName(db, name)
	// 校验数据库用户名
	if dbUser.Id != 0 {
		errRedirect(user, errors.New(i18n.MustGetMessage("i18n.userNameAlreadyExists")))
		return
	}

	nickname := strings.TrimSpace(user.Nickname)
	rem := strings.TrimSpace(user.Rem)

	// ------事务操作------

	// 开始事务
	tx := db.Begin()

	// add
	id, err := tx.Exec("INSERT INTO `user` (`name`, `nickname`, `passwd`, `rem`, `add_time`) VALUES (?, ?, ?, ?, ?)",
		name, user.Nickname, passwdHash, user.Rem, time.NowUnix())
	if err != nil {
		tx.Rollback()
		redirect(user, err)
		return
	}

	// 创建用户数据目录
	dataDir := common.DataDirOnUserId(id)
	if !util_os.IsExist(dataDir) {
		util_os.MkDir(dataDir)
	}
	log.Printf("dataDir: %v\n", dataDir)

	// 复制文件
	dstPath := fmt.Sprintf("%s%s%s", dataDir, util_os.FileSeparator(), "database.db")
	srcPath := fmt.Sprintf("%s%s%s%s%s", common.AppArg.DataDir, util_os.FileSeparator(), "{id}", util_os.FileSeparator(), "database.db")
	_, err = util_os.CopyFile(dstPath, srcPath)
	if err != nil {
		// 回滚事务
		tx.Rollback()
		redirect(user, err)
		return
	}

	// 提交事务
	tx.Commit()

	// 用户注册成功后，重定向到登录页
	api_common_context.Redirect(context, "/user/login", typ.Resp[typ.User]{
		Data: user,
	})
}
