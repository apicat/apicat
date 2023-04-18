package auth

import (
	"errors"
	"time"

	"github.com/apicat/apicat/models"
	"github.com/dgrijalva/jwt-go"
)

type MyClaims struct {
	User *models.Users `json:"user"`
	jwt.StandardClaims
}

// 定义过期时间
const TokenExpireDuration = time.Hour * 168

// 定义secret
var MySecret = []byte("这是一段用于生成token的密钥")

// 生成JWT
func GenerateToken(user *models.Users) (string, error) {
	c := MyClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "apicat",
		},
	}
	//使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	//使用指定的secret签名并获得完成的编码后的字符串token
	return token.SignedString(MySecret)
}

// 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	//解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}

	// TODO 补充「invalid token」的i18n
	return nil, errors.New("invalid token")
}
