package utils

// 各种工具包

import (
	"math/rand"
	"time"
)

func RandomString(n int) string {
	var letters = []byte("asdfkjj;asdf;lasqpoewitupoizx,cvnb/m")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}
