package rnum

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateNumber 生成指定位数的数字字符串
func GenerateNumber(length int) string {
	if length <= 0 {
		return ""
	}

	// 确保第一位不为0
	digits := make([]byte, length)
	digits[0] = byte(rand.Intn(9) + 1 + '0') // 1-9

	// 生成剩余位数
	for i := 1; i < length; i++ {
		digits[i] = byte(rand.Intn(10) + '0') // 0-9
	}

	return string(digits)
}
