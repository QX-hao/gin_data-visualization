package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// getSecret 动态获取 JWT 密钥
func getSecret() []byte {
	secret := viper.GetString("jwt.secret")
	if secret == "" {
		// 如果未配置，使用默认值（不推荐用于生产环境）
		secret = "your-super-secret-jwt-key-change-in-production"
	}
	return []byte(secret)
}

// 生成访问令牌（JWT）
func GenerateAccessToken(userID uint) (string, error) {
	// 读取过期时间配置，默认15分钟
	expire := viper.GetString("jwt.access_token_expire")
	if expire == "" {
		expire = "15m"
	}
	dur, err := time.ParseDuration(expire)
	if err != nil {
		dur = 15 * time.Minute
	}
	expirationTime := time.Now().Add(dur)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getSecret())
	return tokenString, err
}

// 解析并验证 token
func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return getSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}
