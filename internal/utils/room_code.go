package utils

import (
	"math/rand"
	"time"
)

// GenerateRoomCode 生成6位房间代码
func GenerateRoomCode() string {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // 移除易混淆的字符 I, O, 1, 0
	const codeLength = 6

	rand.Seed(time.Now().UnixNano())
	code := make([]byte, codeLength)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}
