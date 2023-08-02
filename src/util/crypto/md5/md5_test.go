// md5 test
// @author xiangqian
// @date 23:01 2023/08/02
package md5

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	data := []byte("123456")
	md5 := Encrypt(data, nil)
	fmt.Println(md5)
}
