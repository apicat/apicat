package random

import (
	"math/rand"
	"strings"
)

func GenerateRandomString(length int) string {
	// 定义包含所有可能字符的字符串
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 生成随机字符串
	sb := strings.Builder{}
	sb.Grow(length)
	for i := 0; i < length; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}
