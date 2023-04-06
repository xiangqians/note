// user test
// @author xiangqian
// @date 21:17 2023/04/06
package test

import (
	"fmt"
	"note/src/api/user"
	"testing"
)

func TestPasswd(t *testing.T) {
	passwd := user.EncryptPasswd("1") // eeafb716f93fa090d7716749a6eefa72
	fmt.Println(passwd)
}
