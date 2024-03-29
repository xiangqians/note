// bcrypt test
// @author xiangqian
// @date 21:21 2023/02/16
package bcrypt

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	encrypt := func(passwd string) string {
		hash, err := Generate(passwd)
		if err != nil {
			panic(err)
		}
		return hash
	}

	passwd := "admin"
	hash := encrypt(passwd)
	fmt.Println(hash)
	fmt.Println(CompareHash(passwd, hash))
	fmt.Println(CompareHash("passwd", hash))
	fmt.Println()

	hash = encrypt(passwd)
	fmt.Println(hash)
	fmt.Println(CompareHash(passwd, hash))
	fmt.Println()
}
