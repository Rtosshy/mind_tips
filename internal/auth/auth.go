package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_secret_key_here") // 秘密鍵

// JWTのペイロード
type Claims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

// トークンを生成
func GenerateJWT(name string) (string, error) {
	claims := Claims{
		Name: name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // トークンの有効期限（24時間）
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// トークンを署名して返す
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// JWTを検証する関数
func ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}
