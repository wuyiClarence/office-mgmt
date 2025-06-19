package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT 配置
var (
	secretKey       = []byte("051f568467f5525c05765fca1671dee7") // 用于生成 JWT token 的密钥
	refreshSecret   = []byte("e4a4c5e266925640a83fb18e1bbd7c5d") // 用于生成 Refresh token 的密钥
	accessTokenTTL  = time.Minute * 30                           // Access token 有效期
	refreshTokenTTL = time.Hour * 24                             // Refresh token 有效期
)

// Claims 定义 JWT 的声明
type Claims struct {
	UserID               int64  `json:"user_id"`   // 用户 ID
	UserName             string `json:"user_name"` // 用户名
	jwt.RegisteredClaims        // 内嵌标准声明
}

// GenerateTokens 生成 Token
func GenerateTokens(userID int64, userName string) (string, string, error) {
	// 创建 Access Token 的 Claims
	accessTokenClaims := Claims{
		UserID:   userID,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims).SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	// 创建 Refresh Token 的 Claims
	refreshTokenClaims := Claims{
		UserID:   userID,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenTTL)),
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims).SignedString(refreshSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// VerifyToken 验证 Token
func VerifyToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func RefreshTokens(refreshTokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(refreshTokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return GenerateAccessToken(claims.UserID, claims.UserName)
	}
	return "", errors.New("invalid refresh token")
}

func GenerateAccessToken(userID int64, userName string) (string, error) {
	claims := Claims{
		UserID:   userID,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey)
}
