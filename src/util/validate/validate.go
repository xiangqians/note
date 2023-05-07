// validate util
// @author xiangqian
// @date 14:30 2023/05/07
package validate

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/i18n"
	"regexp"
)

// FileName 校验文件名
func FileName(name string) error {
	// 名称不能包含字符：
	// \ / : * ? " < > |

	// ^[^\\/:*?"<>|]*$
	matched, err := regexp.MatchString("^[^\\\\/:*?\"<>|]*$", name)
	if err != nil {
		return err
	}

	if !matched {
		return errors.New("文件名不能包含字符：\\ / : * ? \" < > |")
	}

	return nil
}

// Passwd 校验密码
func Passwd(passwd string) error {
	// 1-16位长度（字母，数字，特殊字符）
	matched, err := regexp.MatchString("^[a-zA-Z0-9!@#$%^&*()-_=+]{1,16}$", passwd)
	if err == nil && matched {
		return nil
	}

	return errors.New(fmt.Sprintf(i18n.MustGetMessage("i18n.xMastNBitsLong"), i18n.MustGetMessage("i18n.passwd")))
}

// UserName 校验用户名
func UserName(username string) error {
	if username == "" {
		return errors.New(fmt.Sprintf(i18n.MustGetMessage("i18n.xCannotEmpty"), i18n.MustGetMessage("i18n.userName")))
	}

	// 1-16位长度（字母，数字，下划线，减号）
	matched, err := regexp.MatchString("^[a-zA-Z0-9_-]{1,16}$", username)
	if err == nil && matched {
		return nil
	}

	return errors.New(fmt.Sprintf(i18n.MustGetMessage("i18n.xMastNBitsLong"), i18n.MustGetMessage("i18n.userName")))
}
