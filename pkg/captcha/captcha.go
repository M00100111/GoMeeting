package captcha

import (
	"math/rand"
	"strings"
	"time"
)

// 定义字符集（包含数字、大写字母、小写字母）
const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const captchaLength = 6

// 初始化随机数种子
func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateCaptcha() string {
	var builder strings.Builder
	builder.Grow(captchaLength)

	// 生成随机字符
	for i := 0; i < captchaLength; i++ {
		builder.WriteByte(charset[rand.Intn(len(charset))])
	}

	return builder.String()
}
