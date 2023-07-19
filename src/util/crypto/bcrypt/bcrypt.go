// bcrypt
// @author xiangqian
// @date 20:00 2023/05/08
package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

// CompareHash 对比密码
// hash: 密码密文
// passwd: 密码原文
func CompareHash(hash, passwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd))
}

// Generate 加密密码，每次生成的hash密文都不同
func Generate(passwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
