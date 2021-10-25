package utils

// 各种工具包

import (
	"math/rand"
	"time"
)

// 生成随机字符串的功能，传入一个数字的长度，生成相应长度的随机字符串。
func RandomString(n int) string {
	var letters = []byte("asdfkjj;asdf;lasqpoewitupoizx,cvnb/m")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}
