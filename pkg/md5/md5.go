package md5

import (
	"crypto/md5"
	"fmt"
)

// Encrypt 对字符串进行MD5加密
func Encrypt(data string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

// Verify 验证明文和MD5密文是否匹配
func Verify(plainText, cipherText string) bool {
	// 对明文进行MD5加密
	encrypted := Encrypt(plainText)
	// 比较加密后的结果是否与密文一致
	return encrypted == cipherText
}
