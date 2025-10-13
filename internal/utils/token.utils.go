package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zivwu/reminder-note-api/internal/config"
	"github.com/zivwu/reminder-note-api/internal/consts"
)

type TokenClaims struct {
	UserID string
	Name   string
	Exp    int64
	Iat    int64
	jwt.RegisteredClaims
}

func GenToken(claims TokenClaims) (string, error) {
	// 產生token
	claims.Exp = time.Now().Add(consts.TokenExpireTime).Unix() // 過期時間
	claims.Iat = time.Now().Unix()                             // 發行時間

	// 建立 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 簽名
	signedToken, err := token.SignedString([]byte(config.Env.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ParseToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (any, error) {
		// 驗證簽名方法（防止被改成別的算法）
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Env.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// 驗證 payload 是否正確
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
