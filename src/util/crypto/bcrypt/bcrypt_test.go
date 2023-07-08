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

	passwd := "123456"
	hash := encrypt(passwd)
	fmt.Println(hash)
	fmt.Println(CompareHash(hash, passwd))
	fmt.Println(CompareHash(hash, "passwd"))
	fmt.Println()

	hash = encrypt(passwd)
	fmt.Println(hash)
	fmt.Println(CompareHash(hash, passwd))
	fmt.Println()

}
