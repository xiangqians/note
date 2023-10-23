// 用户注册
// @author xiangqian
// @date 23:33 2023/07/10
package user

import (
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"note/src/api"
	"note/src/context"
	"note/src/model"
	"note/src/session"
	util_crypto_bcrypt "note/src/util/crypto/bcrypt"
	util_os "note/src/util/os"
	"note/src/util/time"
	"note/src/util/validate"
	"os"
	"strings"
)

const signUpUserKey = "signUpUser"

// SignUp 注册
func SignUp(ctx *gin.Context) {
	method, _ := context.PostForm[string](ctx, "_method")
	// 注册
	if method == "POST" {
		signUp(ctx)
	} else
	// 注册页
	{
		user, _ := session.Get[model.AddUser](ctx, signUpUserKey, true)

		msg := ""
		arg := model.GetArg()
		if !arg.AllowSignUp {
			msg = i18n.MustGetMessage("i18n.signUpNotOpen")
		}
		context.HtmlOk(ctx, "user/signup", model.Resp[model.AddUser]{Data: user, Msg: msg})
	}
}

// 注册
func signUp(ctx *gin.Context) {
	// 错误重定向到注册页
	errRedirect := func(tx *gorm.DB, user model.AddUser, err any) {
		if tx != nil {
			// 回滚事务
			tx.Rollback()
		}

		user.Passwd = ""
		user.RePasswd = ""
		session.Set(ctx, signUpUserKey, user)
		context.Redirect(ctx, "/user/signup", nil, err)
	}

	user := model.AddUser{}

	// 是否允许用户注册
	arg := model.GetArg()
	if !arg.AllowSignUp {
		errRedirect(nil, user, i18n.MustGetMessage("i18n.signUpNotOpen"))
		return
	}

	// 绑定
	err := context.ShouldBind(ctx, &user)
	if err != nil {
		errRedirect(nil, user, err)
		return
	}

	user.Name = strings.TrimSpace(user.Name)
	user.Nickname = strings.TrimSpace(user.Nickname)
	user.Passwd = strings.TrimSpace(user.Passwd)
	user.Rem = strings.TrimSpace(user.Rem)

	// 校验用户名
	err = validate.UserName(user.Name)
	if err != nil {
		errRedirect(nil, user, err)
		return
	}

	// 校验密码
	err = validate.Passwd(user.Passwd)
	if err != nil {
		errRedirect(nil, user, err)
		return
	}

	// 加密密码
	passwdHash, err := util_crypto_bcrypt.Generate(user.Passwd)
	if err != nil {
		errRedirect(nil, user, err)
		return
	}

	// 获取数据库操作实例
	db, err := api.Db(nil)
	if err != nil {
		errRedirect(nil, user, err)
		return
	}

	// 根据用户名查询用户信息
	dbUser := getByName(db, user.Name)
	// 校验数据库用户名
	if dbUser.Id != 0 {
		errRedirect(nil, user, i18n.MustGetMessage("i18n.userNameAlreadyExists"))
		return
	}

	// ------事务操作------

	// 开始事务
	tx := db.Begin()

	// 新增用户
	err = tx.Exec("INSERT INTO `user` (`name`, `nickname`, `passwd`, `rem`, `add_time`) VALUES (?, ?, ?, ?, ?)",
		user.Name, user.Nickname, passwdHash, user.Rem, time.NowUnix()).Error
	if err != nil {
		errRedirect(tx, user, err)
		return
	}

	// 获取用户id
	dbUser = getByName(tx, user.Name)
	id := dbUser.Id

	// 创建用户数据目录
	dataDir := arg.DataDir
	userDataDir := util_os.Path(dataDir, fmt.Sprintf("%d", id))
	if !util_os.Stat(userDataDir).IsExist() {
		err = util_os.MkDir(userDataDir, os.ModePerm)
		if err != nil {
			errRedirect(tx, user, err)
			return
		}
	}
	log.Printf("User DataDir %v\n", userDataDir)

	// 复制文件
	srcPath := util_os.Path(dataDir, "id", "database.db")
	dstPath := util_os.Path(userDataDir, "database.db")
	_, err = util_os.CopyFile(srcPath, dstPath)
	if err != nil {
		errRedirect(tx, user, err)
		return
	}

	// 提交事务
	tx.Commit()

	// 用户注册成功后，重定向到登录页
	session.Set(ctx, signInNameKey, user.Name)
	context.Redirect(ctx, "/user/signin", nil, nil)
}
