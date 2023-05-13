// user registration
// @author xiangqian
// @date 19:40 2023/05/11
package user

import (
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"log"
	"note/src/api/common"
	api_common_context "note/src/api/common/context"
	api_common_db "note/src/api/common/db"
	"note/src/api/common/session"
	"note/src/typ"
	"note/src/util/crypto/bcrypt"
	util_os "note/src/util/os"
	"note/src/util/str"
	"note/src/util/time"
	"note/src/util/validate"
	"strings"
)

// Reg0 用户注册
func Reg0(context *gin.Context) {
	// 注册异常时，重定向到注册页
	redirect := func(user typ.User, err any) {
		resp := typ.Resp[typ.User]{
			Msg:  str.ConvTypeToStr(err),
			Data: user,
		}
		api_common_context.Redirect(context, "/user/reg", resp)
	}

	// user
	user := typ.User{}

	// 是否允许用户注册
	if common.AppArg.AllowReg == 0 {
		redirect(user, i18n.MustGetMessage("i18n.regNotOpen"))
		return
	}

	// bind
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

	// 密码加密
	passwdHash, err := bcrypt.Generate(user.Passwd)
	if err != nil {
		redirect(user, err)
		return
	}

	// 校验数据库用户名
	err = ValidateUserName(user.Name)
	if err != nil {
		redirect(user, err)
		return
	}

	user.Nickname = strings.TrimSpace(user.Nickname)
	user.Rem = strings.TrimSpace(user.Rem)

	// db
	db, err := api_common_db.Db(nil)
	if err != nil {
		redirect(user, err)
		return
	}
	defer db.Close()

	// begin
	err = db.Begin()
	if err != nil {
		redirect(user, err)
		return
	}

	// add
	id, err := db.Add("INSERT INTO `user` (`name`, `nickname`, `passwd`, `rem`, `add_time`) VALUES (?, ?, ?, ?, ?)",
		user.Name, user.Nickname, passwdHash, user.Rem, time.NowUnix())
	if err != nil {
		db.Rollback()
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
		db.Rollback()
		redirect(user, err)
		return
	}

	// db commit
	db.Commit()

	// 用户注册成功后，重定向到登录页
	api_common_context.Redirect(context, "/user/login", typ.Resp[typ.User]{
		Data: user,
	})
}

// Reg 用户注册页
func Reg(context *gin.Context) {
	resp, _ := session.Get[typ.Resp[typ.User]](context, api_common_context.RespSessionKey, true)
	api_common_context.HtmlOk(context, "user/reg.html", resp)
}
