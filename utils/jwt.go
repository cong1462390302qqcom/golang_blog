package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtSecret = []byte("your_jwt_secret_key_here")

type Claims struct {
	UserID               uint   `json:"user_id"`
	UserName             string `json:"username"`
	jwt.RegisteredClaims        //引用标准配置
}

func GenerateToken(userID uint, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		UserName: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now())},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if len(jwtSecret) == 0 {
		return "", fmt.Errorf("JWT_SECRET 环境变量未设置")
	}
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	//判断claims的格式是否指定类型
	claims, ok := token.Claims.(*Claims)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
