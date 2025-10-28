package common

import (
	"theing/gin-template/model"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

// 获取 JWT 密钥
func getJWTKey() []byte {
	secret := viper.GetString("jwt.secret")
	if secret == "" {
		// 如果配置中没有设置，使用默认密钥（仅用于开发环境）
		secret = "default-jwt-secret-for-development-only"
	}
	return []byte(secret)
}

// 获取 JWT 过期时间
func getJWTExpirationTime() time.Duration {
	hours := viper.GetInt("jwt.expire_hours")
	if hours <= 0 {
		hours = 168 // 默认 7 天
	}
	return time.Duration(hours) * time.Hour
}

// 获取 JWT 签发者
func getJWTIssuer() string {
	issuer := viper.GetString("jwt.issuer")
	if issuer == "" {
		issuer = "gin-template"
	}
	return issuer
}

// 定义token 的 claims
type Claims struct {
	UserId uint
	jwt.StandardClaims
}

// 调用这个方法发放token
func ReleaseToken(user model.User) (string, error) {
	expirationTime := time.Now().Add(getJWTExpirationTime()) // 从配置获取过期时间
	claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(), // token 发放的时间
			Issuer:    getJWTIssuer(),    // 从配置获取签发者
			Subject:   "user token",      // token 的主题
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // 加密方式应该是hs256
	tokenString, err := token.SignedString(getJWTKey())         // 从配置获取密钥
	if err != nil {
		return "生成token错误", err
	}
	return tokenString, nil
}

// 解析 token 的方法
func ParseToken(tokenstring string) (*jwt.Token, *Claims, error) { // '*' 号表示通过指针修改的内容，返回的也是token 的内存地址的内容。
	claims := &Claims{} // 格式
	token, err := jwt.ParseWithClaims(tokenstring, claims, func(token *jwt.Token) (i interface{}, err error) {
		return getJWTKey(), nil
	})
	return token, claims, err // 解析出claims 然会返回
}
