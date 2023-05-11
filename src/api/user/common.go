// user
// @author xiangqian
// @date 13:20 2023/02/04
package user

import (
	"errors"
	"github.com/gin-contrib/i18n"
	"note/src/api/common/db"
)

// ValidateUserName 校验数据库用户名是否已存在
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
