package main

import (
	"fmt"
	"log"
	"theing/gin-template/utils"
)

func main() {
	fmt.Println("=== 安全密钥生成工具 ===")

	// 生成 JWT 密钥
	jwtKey, err := utils.GenerateJWTKey()
	if err != nil {
		log.Fatalf("生成 JWT 密钥失败: %v", err)
	}
	fmt.Printf("JWT 密钥: %s\n", jwtKey)

	// 生成数据库密码
	dbPassword, err := utils.GenerateSecureKey(16)
	if err != nil {
		log.Fatalf("生成数据库密码失败: %v", err)
	}
	fmt.Printf("数据库密码: %s\n", dbPassword)

	fmt.Println("\n=== 使用说明 ===")
	fmt.Println("1. 将生成的 JWT 密钥设置到环境变量 TPL_JWT_SECRET")
	fmt.Println("2. 将生成的数据库密码设置到相应的环境变量")
	fmt.Println("3. 确保 .env 文件不要提交到版本控制系统")
}
