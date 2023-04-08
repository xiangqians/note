// verify
// @author xiangqian
// @date 21:41 2023/04/08
package common

import (
	"errors"
	"regexp"
)

// VerifyName 校验名称
func VerifyName(name string) error {
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
