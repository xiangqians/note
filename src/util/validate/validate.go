// @author xiangqian
// @date 14:30 2023/05/07
package validate

import (
	"errors"
	"github.com/gin-contrib/i18n"
	"regexp"
)

// UserName 校验用户名
func UserName(userName string) error {
	// 1-60位长度（字母，数字，下划线，减号）
	matched, err := regexp.MatchString("^[a-zA-Z0-9_-]{1,60}$", userName)
	if err != nil {
		return err
	}

	if !matched {
		return errors.New(i18n.MustGetMessage("i18n.validateUserName"))
	}

	return nil
}

// Passwd 校验密码
func Passwd(passwd string) error {
	// 1-20位长度（字母，数字，特殊字符）
	matched, err := regexp.MatchString("^[a-zA-Z0-9!@#$%^&*()-_=+]{1,20}$", passwd)
	if err != nil {
		return err
	}

	if !matched {
		return errors.New(i18n.MustGetMessage("i18n.validatePasswd"))
	}

	return nil
}

// FileName 校验文件名
func FileName(fileName string) error {
	// 名称不能包含字符：\ / : * ? " < > |

	// ^[^\\/:*?"<>|]*$
	matched, err := regexp.MatchString("^[^\\\\/:*?\"<>|]*$", fileName)
	if err != nil {
		return err
	}

	if !matched {
		return errors.New(i18n.MustGetMessage("i18n.validateFileName"))
	}

	return nil
}
