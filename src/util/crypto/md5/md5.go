// md5
// @author xiangqian
// @date 20:13 2023/05/07
package md5

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
)

// Cipher.ENCRYPT_MODE
// Cipher.DECRYPT_MODE

// Encrypt 加密
// data: 要加密的数据
// salt: 盐
func Encrypt(data []byte, salt string) string {
	if salt == "" {
		return fmt.Sprintf("%x", md5.Sum(data))
	}

	d := md5.New()
	str := ""
	for i := 0; i < len(data); i++ {
		str += fmt.Sprintf("%c", data[i])
		if i%2 == 0 {
			str += salt
		}
	}

	_, err := io.WriteString(d, str)
	if err != nil {
		log.Println(err)
	}

	return hex.EncodeToString(d.Sum(nil))
}
