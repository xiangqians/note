// md5
// @author xiangqian
// @date 20:13 2023/05/07
package md5

import (
	"crypto/md5"
	"encoding/hex"
)

// Cipher.ENCRYPT_MODE
// Cipher.DECRYPT_MODE

// Encrypt 加密
// data: 要加密的数据
// salt: 盐
func Encrypt(data []byte, salt []byte) string {
	// 不加盐
	if salt == nil {
		digest := md5.Sum(data)
		return hex.EncodeToString(digest[:])
	}

	// 加盐
	//digest := md5.New()
	//str := ""
	//for i := 0; i < len(data); i++ {
	//	digest.Write(data[:i])
	//	str += fmt.Sprintf("%c", data[i])
	//	if i%2 == 0 {
	//		str += salt
	//	}
	//}
	//
	//_, err := io.WriteString(digest, str)
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//return fmt.Sprintf("%x", digest.Sum(nil))

	panic("暂不支持加盐")
}
