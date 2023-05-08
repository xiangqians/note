// test
// @author xiangqian
// @date 21:21 2023/02/16
package test

import (
	"fmt"
	util_crypto_bcrypt "note/src/util/crypto/bcrypt"
	util_validate "note/src/util/validate"
	"testing"
)

func TestCryptoBcrypt(t *testing.T) {
	encrypt := func(passwd string) string {
		hash, err := util_crypto_bcrypt.Generate(passwd)
		if err != nil {
			panic(err)
		}
		return hash
	}

	passwd := "123456"
	hash := encrypt(passwd)
	fmt.Println(hash)
	fmt.Println(util_crypto_bcrypt.CompareHash(hash, passwd))
	fmt.Println(util_crypto_bcrypt.CompareHash(hash, "passwd"))
	fmt.Println()

	hash = encrypt(passwd)
	fmt.Println(hash)
	fmt.Println(util_crypto_bcrypt.CompareHash(hash, passwd))
	fmt.Println()
}

func TestCryptoMd5(t *testing.T) {

}

func TestValidateName(t *testing.T) {

	// \ / : * ? " < > |

	names := []string{
		"test",
		"test\\",
		"test/",
		"test:",
		"test*",
		"test?",
		"test\"",
		"test<",
		"test>",
		"test|",
		"hello",
		"world",
	}

	for i, name := range names {
		err := util_validate.FileName(name)
		fmt.Println(i, name, err)
	}

}
